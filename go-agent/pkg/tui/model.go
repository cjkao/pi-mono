package tui

import (
	"context"
	"fmt"

	"github.com/badlogic/pi-mono/go-agent/pkg/agent"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	session    *agent.Session
	viewport   viewport.Model
	textarea   textarea.Model
	err        error

	history    string
	streamChan chan string
	thinking   bool
}

type streamStartMsg struct {
	ch chan string
}

type streamChunkMsg string
type streamDoneMsg struct{}
type errMsg error

func InitialModel(session *agent.Session) Model {
	ta := textarea.New()
	ta.Placeholder = "Ask me anything..."
	ta.Focus()
	ta.CharLimit = 0
	ta.SetHeight(3)
	ta.ShowLineNumbers = false

	vp := viewport.New(80, 20)
	vp.SetContent("Welcome to Pi Coding Agent!\n\n")

	return Model{
		session:  session,
		textarea: ta,
		viewport: vp,
		history:  "Welcome to Pi Coding Agent!\n\n",
		err:      nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.thinking {
				v := m.textarea.Value()
				if v == "" {
					return m, nil
				}

				// Clear textarea but don't let it handle the Enter key
				m.textarea.Reset()

				m.history += "> " + v + "\n"
				m.viewport.SetContent(m.history)
				m.viewport.GotoBottom()

				m.thinking = true
				return m, m.startStream(v)
			}
			// If thinking, ignore Enter or maybe cancel?
			return m, nil
		}

	case streamStartMsg:
		m.streamChan = msg.ch
		return m, waitForChunk(m.streamChan)

	case streamChunkMsg:
		content := string(msg)
		m.history += content
		m.viewport.SetContent(m.history)
		m.viewport.GotoBottom()
		return m, waitForChunk(m.streamChan)

	case streamDoneMsg:
		m.thinking = false
		m.streamChan = nil
		m.history += "\n"
		m.viewport.SetContent(m.history)
		m.viewport.GotoBottom()
		return m, nil // tea.Batch(tiCmd, vpCmd) done below

	case errMsg:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		headerHeight := 0
		footerHeight := m.textarea.Height() + 2
		verticalMarginHeight := headerHeight + footerHeight

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
		m.textarea.SetWidth(msg.Width)
	}

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		m.viewport.View(),
		m.textarea.View(),
	)
}

func waitForChunk(sub chan string) tea.Cmd {
	return func() tea.Msg {
		if sub == nil {
			return streamDoneMsg{}
		}
		chunk, ok := <-sub
		if !ok {
			return streamDoneMsg{}
		}
		return streamChunkMsg(chunk)
	}
}

func (m Model) startStream(input string) tea.Cmd {
	return func() tea.Msg {
		ch := make(chan string)
		go func() {
			defer close(ch)
			_, err := m.session.StreamPrompt(context.Background(), input, func(chunk string) {
				ch <- chunk
			})
			if err != nil {
				ch <- fmt.Sprintf("Error: %v\n", err)
			}
		}()
		return streamStartMsg{ch: ch}
	}
}

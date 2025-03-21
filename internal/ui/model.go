package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	ID        string
	Title     string
	Completed bool
}

type CompleteTaskMsg struct {
	ID string
}

// KeyMap defines keybindings for the application
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Complete key.Binding
	Toggle   key.Binding
	Quit     key.Binding
}

// DefaultKeyMap returns a set of default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),
		Complete: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "mark as complete"),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle completion"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (i Item) FilterValue() string { return i.Title }

type Model struct {
	list     list.Model
	choice   string
	quitting bool
	loading  bool
	spinner  spinner.Model
	keyMap   KeyMap
	showHelp bool
}

func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = spinnerStyle

	return Model{
		loading:  true,
		spinner:  s,
		keyMap:   DefaultKeyMap(),
		showHelp: false,
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case CompleteTaskMsg:
		// This message will be handled by the App
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		if !m.loading {
			m.list.SetWidth(msg.Width)
		}
		return m, nil
	case list.Model:
		m.list = msg
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		switch msg.String() {
		case "?":
			// Toggle help view
			m.showHelp = !m.showHelp
			return m, nil
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if i, ok := m.list.SelectedItem().(Item); ok {
				m.choice = i.Title
			}
			return m, tea.Quit
		case "c":
			if i, ok := m.list.SelectedItem().(Item); ok {
				items := m.list.Items()
				i.Completed = true
				items[m.list.Index()] = i
				m.list.SetItems(items)
				return m, func() tea.Msg { return CompleteTaskMsg{ID: i.ID} }
			}
		case " ":
			if i, ok := m.list.SelectedItem().(Item); ok {
				items := m.list.Items()
				i.Completed = !i.Completed
				items[m.list.Index()] = i
				m.list.SetItems(items)
			}
		}
	}

	if !m.loading {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

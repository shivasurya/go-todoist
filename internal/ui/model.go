package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	Title     string
	Completed bool
}

func (i Item) FilterValue() string { return i.Title }

type Model struct {
	list     list.Model
	choice   string
	quitting bool
	loading  bool
	spinner  spinner.Model
}

func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = spinnerStyle

	return Model{
		loading: true,
		spinner: s,
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if i, ok := m.list.SelectedItem().(Item); ok {
				m.choice = i.Title
			}
			return m, tea.Quit
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

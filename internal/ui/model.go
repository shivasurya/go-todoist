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

type TaskCreatedMsg struct {
	Task Item
}

type CreateTaskMsg struct {
	Content     string
	Description string
	DueDate     string
	Priority    int
}

type RefreshTasksMsg struct{}

type AppPage int

const (
	ListPage AppPage = iota
	CreateTaskPage
)

// KeyMap defines keybindings for the application
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Complete key.Binding
	Toggle   key.Binding
	Quit     key.Binding
	New      key.Binding
	Back     key.Binding
	Refresh  key.Binding
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
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new task"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh tasks"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (i Item) FilterValue() string { return i.Title }

type Model struct {
	list           list.Model
	choice         string
	quitting       bool
	loading        bool
	spinner        spinner.Model
	keyMap         KeyMap
	showHelp       bool
	currentPage    AppPage
	taskContent    string
	taskDescription string
	taskDueDate     string
	taskPriority    int
	focusedField    int
}

func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = spinnerStyle

	return Model{
		loading:     true,
		spinner:     s,
		keyMap:      DefaultKeyMap(),
		showHelp:    false,
		currentPage: ListPage,
		taskPriority: 1, // Default to normal priority (P4 in Todoist)
		focusedField: 0,
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
	case TaskCreatedMsg:
		// Add the new task to the list
		items := m.list.Items()
		items = append([]list.Item{msg.Task}, items...)
		m.list.SetItems(items)
		m.currentPage = ListPage
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
		// First handle keys that work on all pages
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		if m.loading {
			return m, nil
		}

		// Handle keys based on current page
		switch m.currentPage {
		case ListPage:
			switch msg.String() {
			case "?":
				// Toggle help view
				m.showHelp = !m.showHelp
				return m, nil
			case "r":
				// Refresh tasks
				m.loading = true
				return m, func() tea.Msg { return RefreshTasksMsg{} }
			case "n":
				// Switch to create task page
				m.currentPage = CreateTaskPage
				m.taskContent = ""
				m.taskDescription = ""
				m.taskDueDate = ""
				m.focusedField = 0 // Focus on the first field
				return m, nil
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

			if !m.loading {
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}
		case CreateTaskPage:

			switch msg.String() {
			case "esc":
				// Go back to list page
				m.currentPage = ListPage
				return m, nil
			case "enter":
				// Submit the form if enter is pressed
				if m.taskContent != "" {
					return m, func() tea.Msg {
						return CreateTaskMsg{
							Content:     m.taskContent,
							Description: m.taskDescription,
							DueDate:     m.taskDueDate,
							Priority:    m.taskPriority,
						}
					}
				}
				return m, nil
			case "tab":
				// Move to next field
				m.focusedField = (m.focusedField + 1) % 4
				return m, nil
			case "shift+tab":
				// Move to previous field
				m.focusedField = (m.focusedField - 1 + 4) % 4
				return m, nil
			case "backspace":
				// Handle backspace for the focused field
				switch m.focusedField {
				case 0:
					if len(m.taskContent) > 0 {
						m.taskContent = m.taskContent[:len(m.taskContent)-1]
					}
				case 1:
					if len(m.taskDescription) > 0 {
						m.taskDescription = m.taskDescription[:len(m.taskDescription)-1]
					}
				case 2:
					if len(m.taskDueDate) > 0 {
						m.taskDueDate = m.taskDueDate[:len(m.taskDueDate)-1]
					}
				case 3:
					// Cycle through priority levels when backspace is pressed
					m.taskPriority = (m.taskPriority % 4) + 1
				}
				return m, nil
			default:
				// Handle typing in the focused field
				if msg.Type == tea.KeyRunes || msg.String() == " " {
					switch m.focusedField {
					case 0:
						m.taskContent += msg.String()
					case 1:
						m.taskDescription += msg.String()
					case 2:
						m.taskDueDate += msg.String()
					case 3:
						// For priority field, we cycle through values with any key press
						m.taskPriority = (m.taskPriority % 4) + 1
					}
				}
				return m, nil
			}
		}
	}

	return m, nil
}

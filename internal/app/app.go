package app

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/shivasurya/go-todoist/internal/todoist"
	"github.com/shivasurya/go-todoist/internal/ui"
	"github.com/shivasurya/go-todoist/pkg/config"
)

type App struct {
	config *config.Config
	client *todoist.Client
	model  ui.Model
}

func New(cfg *config.Config) (*App, error) {
	client := todoist.NewClient(cfg)
	model := ui.NewModel()

	return &App{
		config: cfg,
		client: client,
		model:  model,
	}, nil
}

func completeTask(client *todoist.Client, id string) tea.Cmd {
	return func() tea.Msg {
		err := client.CompleteTask(id)
		if err != nil {
			// We could create an error message here, but for now just return nil
			return nil
		}
		return nil
	}
}

// Define a custom model that can handle CompleteTaskMsg
type todoistModel struct {
	baseModel ui.Model
	client    *todoist.Client
}

// Implement the tea.Model interface
func (m todoistModel) Init() tea.Cmd {
	return m.baseModel.Init()
}

func (m todoistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.CompleteTaskMsg:
		return m, completeTask(m.client, msg.ID)
	default:
		updatedModel, cmd := m.baseModel.Update(msg)
		m.baseModel = updatedModel.(ui.Model)
		return m, cmd
	}
}

func (m todoistModel) View() string {
	return m.baseModel.View()
}

func (a *App) Run() error {
	// Create our custom model that adds task completion functionality
	customModel := todoistModel{
		baseModel: a.model,
		client:    a.client,
	}

	p := tea.NewProgram(customModel, tea.WithAltScreen())

	go func() {
		tasks, err := a.client.GetTasks()
		if err != nil {
			return
		}

		var items []list.Item
		for _, task := range tasks {
			if task.Due.Date == time.Now().Format("2006-01-02") {
				items = append(items, ui.Item{
					ID:        task.Id,
					Title:     task.Content,
					Completed: task.Completed,
				})
			}
		}

		l := list.New(items, ui.ItemDelegate{}, a.config.DefaultWidth, a.config.ListHeight)
		l.Title = "Todoist Tasks"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = ui.TitleStyle
		l.Styles.PaginationStyle = ui.PaginationStyle
		l.Styles.HelpStyle = ui.HelpStyle

		// Set up the list view with task navigation and completion options
		l.Title = "Todoist Tasks"
		l.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{
				key.NewBinding(
					key.WithKeys("c"),
					key.WithHelp("c", "mark task as complete"),
				),
			}
		}
		l.SetStatusBarItemName("task", "tasks")
		l.SetFilteringEnabled(false)
		l.SetShowHelp(true)
		l.SetShowStatusBar(true)

		p.Send(l)
	}()

	_, err := p.Run()
	return err
}

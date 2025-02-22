package app

import (
	"time"

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

func (a *App) Run() error {
	p := tea.NewProgram(a.model)

	go func() {
		tasks, err := a.client.GetTasks()
		if err != nil {
			return
		}

		var items []list.Item
		for _, task := range tasks {
			if task.Due.Date == time.Now().Format("2006-01-02") {
				items = append(items, ui.Item{
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

		p.Send(l)
	}()

	_, err := p.Run()
	return err
}

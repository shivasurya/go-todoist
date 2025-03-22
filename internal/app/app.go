package app

import (
	"strings"
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

func createTask(client *todoist.Client, msg ui.CreateTaskMsg) tea.Cmd {
	return func() tea.Msg {
		// Todoist uses priorities 1-4 where:
		// P1 = normal (1)
		// P2 = medium (2)
		// P3 = high (3)
		// P4 = urgent (4)
		// But we need to reverse the mapping since the UI selection is inverted
		todoistPriority := 5 - msg.Priority

		// Format due string using natural language input
		// Combine date and time into a single due string
		dueString := ""

		// If we have a date, start with that
		if msg.DueDate != "" {
			dueString = msg.DueDate

			// Add time if specified
			if msg.DueTime != "" {
				dueString += " at " + msg.DueTime
			}
		} else if msg.DueTime != "" {
			// If only time specified, assume it's today
			dueString = "today at " + msg.DueTime
		}

		taskReq := todoist.CreateTaskRequest{
			Content:     msg.Content,
			Description: msg.Description,
			DueString:   dueString,
			Priority:    todoistPriority,
		}

		newTask, err := client.CreateTask(taskReq)
		if err != nil {
			// Handle error (in a real app, we'd return an error message)
			return nil
		}

		// Determine if this is a recurring task based on the response
		isRecurring := false
		if newTask.Due.IsRecurring {
			isRecurring = true
		}

		// Extract time if present in the datetime string
		dueTime := ""
		if newTask.Due.Datetime != "" {
			// If datetime is available, try to extract just the time portion
			// This is simplified - in a real app you would parse the datetime properly
			parts := strings.Split(newTask.Due.Datetime, "T")
			if len(parts) > 1 {
				// Extract time and remove seconds/timezone
				timeOnly := strings.Split(parts[1], "+")
				dueTime = strings.Split(timeOnly[0], ":")[0] + ":" + strings.Split(timeOnly[0], ":")[1]
			}
		}

		// Format the due date for display
		dueDate := newTask.Due.Date
		if newTask.Due.String != "" {
			// Extract just the date part without time for consistent display
			dateTimeParts := strings.Split(newTask.Due.String, " at ")
			if len(dateTimeParts) > 0 {
				dueDate = dateTimeParts[0]
			}
		}

		// Return the created task to add it to the UI
		return ui.TaskCreatedMsg{
			Task: ui.Item{
				ID:          newTask.Id,
				Title:       newTask.Content,
				Completed:   newTask.Completed,
				Priority:    todoistPriority,
				DueDate:     dueDate,
				DueTime:     dueTime,
				IsRecurring: isRecurring,
			},
		}
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
	case ui.CreateTaskMsg:
		return m, createTask(m.client, msg)
	case ui.RefreshTasksMsg:
		return m, fetchTasksForDate(m.client, m.baseModel.CurrentDate, m.baseModel.ShowOverdue)
	case ui.ChangeTaskDateMsg:
		return m, fetchTasksForDate(m.client, m.baseModel.CurrentDate, m.baseModel.ShowOverdue)
	case ui.GoToTodayMsg:
		return m, fetchTasksForDate(m.client, time.Now(), m.baseModel.ShowOverdue)
	case ui.ToggleOverdueMsg:
		return m, fetchTasksForDate(m.client, m.baseModel.CurrentDate, m.baseModel.ShowOverdue)
	default:
		updatedModel, cmd := m.baseModel.Update(msg)
		m.baseModel = updatedModel.(ui.Model)
		return m, cmd
	}
}

func (m todoistModel) View() string {
	return m.baseModel.View()
}

func fetchTasksForDate(client *todoist.Client, date time.Time, showOverdue bool) tea.Cmd {
	return func() tea.Msg {
		// Format the date for the API query
		targetDate := date.Format("2006-01-02")
		var allTasks []todoist.Task

		// Get either date tasks or overdue tasks based on filter
		if showOverdue {
			// Only show overdue tasks when filtered
			overdueTasks, err := client.GetOverdueTasks()
			if err != nil {
				return nil
			}
			allTasks = overdueTasks
		} else {
			// Show regular date tasks
			dateTasks, err := client.GetTasksByDate(targetDate)
			if err != nil {
				return nil
			}
			allTasks = dateTasks
		}

		// Convert tasks to UI items
		var items []list.Item
		for _, task := range allTasks {
			// Extract time if present in the datetime string
			dueTime := ""
			if task.Due.Datetime != "" {
				// If datetime is available, extract just the time portion
				parts := strings.Split(task.Due.Datetime, "T")
				if len(parts) > 1 {
					// Extract time and remove seconds/timezone
					timeOnly := strings.Split(parts[1], "+")
					dueTime = strings.Split(timeOnly[0], ":")[0] + ":" + strings.Split(timeOnly[0], ":")[1]
				}
			}

			// If we have a due string from the API, use it for display
			dueDate := task.Due.Date
			if task.Due.String != "" {
				// Extract just the date part without time for consistent display
				dateTimeParts := strings.Split(task.Due.String, " at ")
				if len(dateTimeParts) > 0 {
					dueDate = dateTimeParts[0]
				}
			}

			items = append(items, ui.Item{
				ID:          task.Id,
				Title:       task.Content,
				Completed:   task.Completed,
				Priority:    task.Priority,
				DueDate:     dueDate,
				DueTime:     dueTime,
				IsRecurring: task.Due.IsRecurring,
			})
		}

		l := list.New(items, ui.ItemDelegate{}, 20, 14) // Using hardcoded width and height for refresh
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = ui.TitleStyle
		l.Styles.PaginationStyle = ui.PaginationStyle
		l.Styles.HelpStyle = ui.HelpStyle
		l.Title = "" // Set empty title to avoid duplicate headings

		// Set up the list view with task navigation and completion options
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

		return l
	}
}

func (a *App) Run() error {
	// Create our custom model that adds task completion functionality
	customModel := todoistModel{
		baseModel: a.model,
		client:    a.client,
	}

	p := tea.NewProgram(customModel, tea.WithAltScreen())

	go func() {
		// Use fetchTasksForDate to load today's tasks initially
		cmd := fetchTasksForDate(a.client, time.Now(), false)
		l := cmd()
		p.Send(l)
	}()

	_, err := p.Run()
	return err
}

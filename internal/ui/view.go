package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ItemDelegate struct{}

func (d ItemDelegate) Height() int  { return 1 }
func (d ItemDelegate) Spacing() int { return 0 }
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// You could add more list-specific key handlers here
		// This would intercept keys before they reach the main model
		}
	}
	return nil
}
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, item.Title)

	if item.Completed {
		str = strikedStyle.Render(str)
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (m Model) View() string {
	if m.loading {
		return fmt.Sprintf("\n %s Loading your Todoist tasks...\n", m.spinner.View())
	}
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("Task selected: %s", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Thanks for using Todoist CLI!")
	}

	if m.showHelp {
		helpView := "\n ✨ Todoist TUI Keyboard Controls \n\n"
		helpView += " • j/↓: Move cursor down\n"
		helpView += " • k/↑: Move cursor up\n"
		helpView += " • c: Mark task as complete\n"
		helpView += " • Space: Toggle task completion status\n"
		helpView += " • n: Create new task\n"
		helpView += " • r: Refresh tasks\n"
		helpView += " • Enter: Select task\n"
		helpView += " • q/Ctrl+C: Quit\n"
		helpView += " • ?: Toggle this help menu\n\n"
		helpView += " Press any key to return to tasks"
		return helpStyle.Render(helpView)
	}

	switch m.currentPage {
	case ListPage:
		return "\n" + m.list.View() + "\n\n" + subtleStyle.Render(" Press ? for help, n for new task, r to refresh ")
	case CreateTaskPage:
		s := "\n" + titleStyle.Render("Create New Task") + "\n\n"

		// Task content field
		if m.focusedField == 0 {
			s += focusedInputStyle.Render("Task: " + m.taskContent + "_") + "\n"
		} else {
			s += unfocusedInputStyle.Render("Task: " + m.taskContent) + "\n"
		}

		// Description field
		if m.focusedField == 1 {
			s += focusedInputStyle.Render("Description: " + m.taskDescription + "_") + "\n"
		} else {
			s += unfocusedInputStyle.Render("Description: " + m.taskDescription) + "\n"
		}

		// Due date field
		if m.focusedField == 2 {
			s += focusedInputStyle.Render("Due date: " + m.taskDueDate + "_") + "\n"
		} else {
			s += unfocusedInputStyle.Render("Due date: " + m.taskDueDate) + "\n"
		}

		// Priority field - show priority levels from P1 (normal) to P4 (urgent)
		priorityLabels := map[int]string{1: "P1 (normal)", 2: "P2 (medium)", 3: "P3 (high)", 4: "P4 (urgent)"}
		if m.focusedField == 3 {
			s += focusedInputStyle.Render("Priority: " + priorityLabels[m.taskPriority] + " [press any key to cycle]") + "\n"
		} else {
			s += unfocusedInputStyle.Render("Priority: " + priorityLabels[m.taskPriority]) + "\n"
		}

		s += "\n" + subtleStyle.Render(" Tab: Next field • Enter: Submit • Esc: Cancel ") + "\n"
		return s
	default:
		return "Unknown page"
	}
}

package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
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
	return "\n" + m.list.View()
}

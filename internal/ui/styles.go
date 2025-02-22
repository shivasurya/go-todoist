package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.AdaptiveColor{Light: "#F72585", Dark: "#7209B7"})
	PaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	spinnerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	strikedStyle      = lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color("240"))
)

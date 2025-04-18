package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle           = lipgloss.NewStyle().MarginLeft(2)
	titleStyle           = lipgloss.NewStyle().MarginLeft(2).Bold(true)
	itemStyle            = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle    = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.AdaptiveColor{Light: "#F72585", Dark: "#7209B7"})
	PaginationStyle      = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle            = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle        = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	spinnerStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	strikedStyle         = lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color("240"))
	StatusMessageStyle   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04D483"}).Bold(true)
	helpStyle            = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63"))
	subtleStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Italic(true)
	focusedInputStyle    = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#F72585", Dark: "#7209B7"}).Bold(true)
	unfocusedInputStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	dueStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("110")).Italic(true)
	overdueStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Background(lipgloss.Color("196")).Bold(true).Padding(0, 1).BorderStyle(lipgloss.RoundedBorder())
	overdueSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Background(lipgloss.Color("196")).Bold(true).Padding(0, 1).BorderStyle(lipgloss.RoundedBorder())
)

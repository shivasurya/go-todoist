# 📝 Golang ToDoist Client App with Bubble Tea UI 🫧

<img height="32" width="32" src="https://cdn.simpleicons.org/todoist" /> <img height="32" width="32" src="https://cdn.simpleicons.org/go" />

</br>

🚧 **This project is a work in progress.** 🚧

A delightful terminal-based Todoist client built with Go and the charming Bubble Tea framework.

## ✨ Features

- 📋 List tasks
- 🎨 Beautiful terminal UI
- ⌨️ Keyboard-driven interface

## 🚀 Installation

go install github.com/shivasurya/go-todoist/cmd/todoist@latest

## 🎮 Usage

1. Launch the application:

`todoist` should be available in your $GOPATH/bin directory.

2. Navigate using these keyboard shortcuts:
- `j/k` or `↑/↓`: Navigate tasks
- `space`: Toggle task completion status
- `c`: Mark task as complete
- `?`: Show help menu
- `q`: Quit application

## 🛠️ Development

Requirements:
- Go 1.19 or higher
- Todoist API key

Build from source:

- git clone https://github.com/shivasurya/go-todoist
- cd todoist-tui
- export TODOIST_TOKEN=your_api_key
- go build -o todo cmd/todoist/main.go

Ensure you have a valid Todoist API key in your environment variable as TODOIST_TOKEN.

## 📄 License

MIT License - feel free to use and modify!

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the amazing TUI framework
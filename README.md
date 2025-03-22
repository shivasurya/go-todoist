# 📝 Golang ToDoist Client App with Bubble Tea UI 🫧

<img height="32" width="32" src="https://cdn.simpleicons.org/todoist" /> <img height="32" width="32" src="https://cdn.simpleicons.org/go" />

</br>

🚧 **This project is a work in progress.** 🚧

A delightful terminal-based Todoist client built with Go and the charming Bubble Tea framework.

## ✨ Features

- 📋 List tasks and manage them from your terminal
- 📅 Browse tasks by date with day-to-day navigation
- ➕ Create new tasks with descriptions, due dates, times, and priority levels
- 📓️ Quick toggle between daily view and overdue tasks with visual indicators
- 🔁 Support for recurring tasks (daily, weekly, monthly, etc.)
- ⏰ Add specific times to your tasks
- ✓ Mark tasks as complete or toggle completion status
- 🔄 Refresh task list to sync with Todoist
- 📆 Assign priority levels (P1-P4) to your tasks
- 🎨 Beautiful terminal UI with intuitive navigation
- ⌨️ Fully keyboard-driven interface

## 🚀 Installation

go install github.com/shivasurya/go-todoist/cmd/todoist@latest

## 🎮 Usage

1. Launch the application:

`todoist` should be available in your $GOPATH/bin directory.

2. Navigate using these keyboard shortcuts:

### Task List Navigation
- `j/k` or `↑/↓`: Navigate up and down through tasks
- `h/l` or `←/→`: Navigate between previous/next day's tasks
- `t`: Jump to today's tasks
- `o`: Toggle between viewing current day's tasks and all overdue tasks
- `space`: Toggle task completion status
- `c`: Mark task as complete
- `n`: Create new task
- `r`: Refresh tasks list from Todoist
- `Enter`: Select currently highlighted task
- `?`: Toggle help menu
- `q` or `Ctrl+C`: Quit application

### Task Creation
- `Tab/Shift+Tab`: Navigate between task fields (Title, Description, Due Date, Due Time, Priority)
- Any key: Cycle through priority levels when the priority field is selected:
  - P1: Normal priority (default)
  - P2: Medium priority
  - P3: High priority
  - P4: Urgent priority
- Due date field: Enter dates using natural language (e.g., `tomorrow`, `next Monday`, `Dec 25`)
- Due time field: Enter time using natural language (e.g., `9am`, `2pm`, `morning`, `evening`)
- `Enter`: Submit and create the task
- `Esc`: Cancel and return to task list

## 🛠️ Development

Requirements:
- Go 1.19 or higher
- Todoist API key

Build from source:

```bash
# Clone the repository
git clone https://github.com/shivasurya/go-todoist

# Navigate to the project directory
cd go-todoist

# Set your Todoist API token
export TODOIST_TOKEN=your_api_key

# Build the application
go build -o todoist cmd/todoist/main.go

# Run it
./todoist
```

Ensure you have a valid Todoist API key in your environment variable as TODOIST_TOKEN.

## 📄 License

MIT License - feel free to use and modify!

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the amazing TUI framework
package todoist

import "time"

type Project struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	CommentCount int    `json:"comment_count"`
	Order        int    `json:"order"`
	Color        string `json:"color"`
	IsShared     bool   `json:"is_shared"`
	IsFavorite   bool   `json:"is_favorite"`
	ParentId     int    `json:"parent_id"`
	IsInbox      bool   `json:"is_inbox_project"`
	IsTeamInbox  bool   `json:"is_team_inbox"`
	ViewStyle    string `json:"view_style"`
	URL          string `json:"url"`
}

type Task struct {
	Id           string    `json:"id"`
	ProjectId    string    `json:"project_id"`
	Content      string    `json:"content"`
	Description  string    `json:"description"`
	Completed    bool      `json:"completed"`
	Labels       []string  `json:"labels"`
	LabelIds     []string  `json:"label_ids"`
	Priority     int       `json:"priority"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	Due          struct {
		Date        string `json:"date"`
		Datetime    string `json:"datetime"`
		String      string `json:"string"`
		IsRecurring bool   `json:"is_recurring"`
		Timezone    string `json:"timezone,omitempty"`
	} `json:"due"`
}

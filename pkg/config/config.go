package config

import "os"

type Config struct {
	Token          string
	ProjectURLPath string
	TaskURLPath    string
	BaseURL        string
	DefaultWidth   int
	ListHeight     int
}

func New() *Config {
	return &Config{
		Token:          os.Getenv("TODOIST_TOKEN"),
		ProjectURLPath: "projects",
		TaskURLPath:    "tasks",
		BaseURL:        "https://api.todoist.com/rest/v2",
		DefaultWidth:   20,
		ListHeight:     14,
	}
}

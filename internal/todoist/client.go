package todoist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shivasurya/go-todoist/pkg/config"
)

type Client struct {
	config *config.Config
	http   *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		http: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) GetProjects() ([]Project, error) {
	body, err := c.makeGetRequest("projects")
	if err != nil {
		return nil, err
	}

	var projects []Project
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (c *Client) GetTasks() ([]Task, error) {
	body, err := c.makeGetRequest("tasks")
	if err != nil {
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (c *Client) makeGetRequest(feature string) ([]byte, error) {
	url := c.resolveFeatureURL(feature)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.config.Token)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) CompleteTask(taskID string) error {
	url := c.resolveFeatureURL("tasks/" + taskID + "/close")
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.config.Token)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) resolveFeatureURL(feature string) string {
	switch feature {
	case "projects":
		return c.config.BaseURL + "/" + c.config.ProjectURLPath
	case "tasks":
		return c.config.BaseURL + "/" + c.config.TaskURLPath
	default:
		return c.config.BaseURL + "/" + feature
	}
}

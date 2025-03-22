package todoist

import (
	"bytes"
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

type CreateTaskRequest struct {
	Content     string   `json:"content"`
	ProjectID   string   `json:"project_id,omitempty"`
	Description string   `json:"description,omitempty"`
	DueDate     string   `json:"due_string,omitempty"`
	LabelIDs    []string `json:"label_ids,omitempty"`
	Priority    int      `json:"priority,omitempty"`
}

func (c *Client) CreateTask(task CreateTaskRequest) (*Task, error) {
	url := c.resolveFeatureURL("tasks")

	data, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.config.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// Read the response body to get more details on the error
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned status code %d: %s", resp.StatusCode, string(respBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var newTask Task
	if err := json.Unmarshal(body, &newTask); err != nil {
		return nil, err
	}

	return &newTask, nil
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

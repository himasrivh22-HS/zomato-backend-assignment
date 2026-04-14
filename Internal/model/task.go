package model

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	ProjectID   string `json:"project_id"`
	AssigneeID  string `json:"assignee_id"`
	CreatorID   string `json:"creator_id"`   
	DueDate     string `json:"due_date"`     
}
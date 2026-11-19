package domain

type Status string

const (
	StatusInProgress Status = "in_progress"
	StatusReady      Status = "ready"
)

type Task struct {
	Code     string
	Language string
	Status   Status
	ID       string
	Result   string
}

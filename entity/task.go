package entity

// Task ...
type Task struct {
	ID       uint
	Text     string
	Done     bool
	AssignTo User
}

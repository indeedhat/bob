package transaction

import (
	"os"
)

type FileMakerFunc func(*TaskLog[string])

type FileMaker struct {
	Tasks *TaskLog[string]
	fn    FileMakerFunc
}

// NewFileMaker creates a new FileMaker instance and runs the maker func against it
func NewFileMaker(f FileMakerFunc) *FileMaker {
	return &FileMaker{
		fn:    f,
		Tasks: new(TaskLog[string]),
	}
}

func (m *FileMaker) Run() {
	m.fn(m.Tasks)
}

// Rollback reverts changes made by the maker
func (m *FileMaker) Rollback() {
	for _, task := range *m.Tasks {
		os.Remove(task.Task)
	}
}

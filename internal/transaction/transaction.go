package transaction

type TaskEntry[T any] struct {
	Task  T
	Error error
}

type TaskLog[T any] []TaskEntry[T]

func (t *TaskLog[T]) Do(task T, fn func() error) {
	t.Add(task)
	t.Failed(fn())
}

// Add adds a task entry to the task log
// this task can then be interpreted by the runner to revert any changes made if needed
func (t *TaskLog[T]) Add(task T) {
	*t = append(*t, TaskEntry[T]{task, nil})
}

// Failed tracks the error for the current task
func (t *TaskLog[T]) Failed(err error) {
	(*t)[len(*t)-1].Error = err
}

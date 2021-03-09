package task

// Task should be ignored for now
type Task interface {
	Process() (string, error)
	Send() error // NOTE dummy to satisfy machinery.Task requirement
}

package task

import "context"

type Doer interface {
	Do(ctx context.Context) error
}

type Undoer interface {
	Undo(ctx context.Context) error
}

type Interface interface {
	Doer
	Undoer
}

// Do the given tasks in order
func Do(ctx context.Context, tasks ...Interface) error {
	for _, task := range tasks {
		err := task.Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Undo the given tasks in order
func Undo(ctx context.Context, tasks ...Interface) error {
	for _, task := range tasks {
		err := task.Undo(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Go runs the given tasks concurrently
func Go(ctx context.Context, tasks ...Interface) error {
	return nil
}

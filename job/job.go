package job

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

var ErrEmptyQueue = errors.New("queue is empty")

// Message is a wrapper around the job payload with common metadata
type Message struct {
	ID        string
	Payload   []byte
	Timestamp time.Time
}

type Error struct {
	error
	*Message
}

// Handler handles incoming messages
type Handler interface {
	Handle(ctx context.Context, msg *Message) error
	HandleError(ctx context.Context, err error) error
}

type Queue interface {
	// Push messages to the queue
	Push(ctx context.Context, messages ...*Message) error
	// Pull jobs from the queue and handle it
	Pull(ctx context.Context, handle func(ctx context.Context, msg *Message) error) error
	// Worker creates a worker for the queue
	Worker(handler Handler) Worker
}

// type Handle func(ctx context.Context, msg *Message) error

// var _ Handler = (Handle)(nil)

// func (fn Handle) Handle(ctx context.Context, msg *Message) error {
// 	return fn(ctx, msg)
// }

// func (fn Handle) HandleError(ctx context.Context, err error) error {
// 	return err
// }

// Worker continuously pulls jobs from the queue and handles them
type Worker interface {
	Work(ctx context.Context) error
}

func Looper(queue Queue, handler Handler) Worker {
	return &looper{queue, handler}
}

type looper struct {
	queue   Queue
	handler Handler
}

func (l *looper) Work(ctx context.Context) error {
	handle := errorWrap(recoverWrap(l.handler.Handle))
	for {
		if err := l.queue.Pull(ctx, handle); err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			if err := l.handler.HandleError(ctx, err); err != nil {
				return err
			}
		}
	}
}

// Loop over the queue and handle incoming messages
func Loop(ctx context.Context, queue Queue, handler Handler) error {
	handle := errorWrap(recoverWrap(handler.Handle))
	for {
		if err := queue.Pull(ctx, handle); err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			if err := handler.HandleError(ctx, err); err != nil {
				return err
			}
		}
	}
}

func Work(ctx context.Context, workers ...Worker) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, worker := range workers {
		worker := worker
		eg.Go(func() error {
			return worker.Work(ctx)
		})
	}
	return eg.Wait()
}

func recoverWrap(fn func(ctx context.Context, msg *Message) error) func(ctx context.Context, msg *Message) error {
	return func(ctx context.Context, msg *Message) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("job: recovered from panic: %v", p)
			}
		}()
		return fn(ctx, msg)
	}
}

func errorWrap(fn func(ctx context.Context, msg *Message) error) func(ctx context.Context, msg *Message) error {
	return func(ctx context.Context, msg *Message) error {
		if err := fn(ctx, msg); err != nil {
			return Error{err, msg}
		}
		return nil
	}
}

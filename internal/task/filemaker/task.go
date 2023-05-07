package filemaker

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/matthewmueller/bud/internal/task"
	"github.com/matthewmueller/bud/log"
)

func New(log log.Log, path string, data []byte, force bool) *Task {
	return &Task{log, path, data, force}
}

type Task struct {
	log   log.Log
	path  string
	data  []byte
	force bool
}

var _ task.Interface = (*Task)(nil)

func (m *Task) Do(ctx context.Context) error {
	existing, err := os.ReadFile(m.path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return m.createFile(ctx)
	}
	if bytes.Equal(existing, m.data) {
		m.log.Infof("identical %s", m.path)
		return nil
	}
	if m.force {
		m.log.Infof("overwriting %s", m.path)
		return m.createFile(ctx)
	}
	return fmt.Errorf("%s already exists", m.path)
}

func (m *Task) createFile(ctx context.Context) error {
	m.log.Infof("creating %s", m.path)
	if err := os.MkdirAll(filepath.Dir(m.path), 0755); err != nil {
		return err
	}
	return os.WriteFile(m.path, m.data, 0644)
}

func (m *Task) Undo(ctx context.Context) error {
	return fmt.Errorf("undo %s is not implemented", m.path)
}

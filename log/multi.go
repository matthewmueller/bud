package log

import "errors"

// Multi logs to multiple places
func Multi(handlers ...Handler) Handler {
	return &multi{handlers}
}

type multi struct {
	handlers []Handler
}

func (m *multi) Log(entry *Entry) (err error) {
	for _, handler := range m.handlers {
		err = errors.Join(err, handler.Log(entry))
	}
	return err
}

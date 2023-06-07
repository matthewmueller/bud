package slot

import (
	"bytes"
	"io"
	"sync"
)

func New() *Set {
	return &Set{
		main:  newPipe(),
		pipes: map[string]*Pipe{},
		prev:  nil,
	}
}

type Set struct {
	main  *Pipe
	mu    sync.RWMutex
	pipes map[string]*Pipe
	prev  *Set
}

var _ io.ReadWriteCloser = (*Set)(nil)

func (s *Set) Write(b []byte) (n int, err error) {
	return write(s.main, b)
}

func (s *Set) pipe(name string) (*Pipe, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	pipe, ok := s.pipes[name]
	return pipe, ok
}

func (s *Set) WriteTo(name string, b []byte) (n int, err error) {
	pipe, ok := s.pipe(name)
	if !ok {
		s.mu.Lock()
		s.pipes[name] = newPipe()
		pipe = s.pipes[name]
		s.mu.Unlock()
	}
	return write(pipe, b)
}

func write(pipe *Pipe, b []byte) (n int, err error) {
	return pipe.Write(b)
}

func (s *Set) Read(b []byte) (n int, err error) {
	if s.prev == nil {
		return 0, io.EOF
	}
	return read(s.prev.main, b)
}

func (s *Set) Reader() (io.Reader, bool) {
	if s.prev == nil {
		return nil, false
	}
	return s.prev.main, true
}

func (s *Set) ReaderFrom(name string) (io.Reader, bool) {
	if s.prev == nil {
		return nil, false
	}
	pipe, ok := s.prev.pipes[name]
	return pipe, ok
}

func (s *Set) Readers() []io.Reader {
	if s.prev == nil {
		return nil
	}
	return append([]io.Reader{s.prev.main}, s.prev.Readers()...)
}

func (s *Set) ReadersFrom(name string) (readers []io.Reader) {
	if s.prev == nil {
		return readers
	}
	pipe, ok := s.prev.pipes[name]
	if ok {
		readers = append(readers, pipe)
	}
	return append(readers, s.prev.Readers()...)
}

func read(pipe *Pipe, b []byte) (n int, err error) {
	return pipe.Read(b)
}

func (s *Set) Close() error {
	s.main.Close()
	for _, pipe := range s.pipes {
		pipe.Close()
	}
	return nil
}

func (s *Set) New() *Set {
	slot := New()
	slot.prev = s
	return slot
}

func (s *Set) Slot(name string) io.ReadWriter {
	return s.New()
}

func newPipe() *Pipe {
	return &Pipe{
		done: make(chan struct{}),
	}
}

type Pipe struct {
	mu   sync.Mutex
	b    bytes.Buffer  // Written data
	done chan struct{} // Writes are done
}

var _ io.ReadWriteCloser = (*Pipe)(nil)

func (p *Pipe) Write(b []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.b.Write(b)
}

func (p *Pipe) Read(b []byte) (int, error) {
	<-p.done
	return p.b.Read(b)
}

func (p *Pipe) Close() error {
	close(p.done)
	return nil
}

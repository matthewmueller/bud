package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/matthewmueller/bud/router"
	"golang.org/x/sync/errgroup"
)

type Input struct {
	Props   any
	slots   Slots
	Context map[string]any
}

type Slots map[string][]io.Reader

func (i *Input) Slot() (template.HTML, error) {
	return i.slot(i.slots["main"]...)
	// if len(names) == 0 {
	// 	return i.slot(i.slots["main"]...)
	// }
	// return i.slot(i.slots[names[0]]...)
}

func (i *Input) slot(readers ...io.Reader) (template.HTML, error) {
	var buf bytes.Buffer
	for _, r := range readers {
		_, err := io.Copy(&buf, r)
		if err != nil {
			return "", err
		}
	}
	return template.HTML(buf.String()), nil
}

func (i *Input) Title(title string) {
	// i.slots["title"] = title
}

func (i *Input) Meta(key, value string) {
	// i.Context["title"] = title
}

func main() {

	indexTpl := template.Must(template.New("index.gohtml").Parse(`{{ with $.Props }}Hello {{ . }}!{{ end }}`))
	frameTpl := template.Must(template.New("frame.gohtml").Parse(`{{ with $.Props }}<main>{{ $.Slot }}</main>{{ end }}`))
	layoutTpl := template.Must(template.New("layout.gohtml").Parse(`{{ with $.Props }}
		<html>
			<body>{{ $.Slot }}</body>
		</html>
	{{ end }}`))

	index := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("rendering index")
		time.Sleep(100 * time.Millisecond)
		indexTpl.Execute(w, &Input{
			Props:   "Earth",
			slots:   Slots{},
			Context: map[string]any{},
		})
		fmt.Println("rendered index")
	})
	frame := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("rendering frame")
		time.Sleep(100 * time.Millisecond)
		frameTpl.Execute(w, &Input{
			Props: struct{}{},
			slots: Slots{
				"main": []io.Reader{r.Body},
			},
			Context: map[string]any{},
		})
		fmt.Println("rendered frame")
	})
	layout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("rendering layout")
		layoutTpl.Execute(w, &Input{
			Props: struct{}{},
			slots: Slots{
				"main": []io.Reader{r.Body},
			},
			Context: map[string]any{},
		})
		fmt.Println("rendered layout")
	})

	r := router.New()
	r.Get("/index", index)
	r.Get("/frame", frame)
	r.Get("/layout", layout)
	r.Get("/", wrap(layout, frame, index))
	fmt.Println("listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

type Pipes map[string]*Pipe

func (pipes Pipes) Close() {
	for _, pipe := range pipes {
		pipe.Close()
	}
}

func wrap(layout, frame, index http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		eg := new(errgroup.Group)
		indexPipes := Pipes{
			"main": &Pipe{
				done: make(chan struct{}),
			},
		}
		framePipes := Pipes{
			"main": &Pipe{
				done: make(chan struct{}),
			},
		}
		eg.Go(func() error {
			index.ServeHTTP(&ResponseWriter{indexPipes["main"]}, r)
			indexPipes.Close()
			return nil
		})
		eg.Go(func() error {
			r.Body = indexPipes["main"]
			frame.ServeHTTP(&ResponseWriter{framePipes["main"]}, r)
			framePipes.Close()
			return nil
		})
		eg.Go(func() error {
			r.Body = framePipes["main"]
			layout.ServeHTTP(w, r)
			return nil
		})
		eg.Wait()
	})
}

type ResponseWriter struct {
	w io.Writer
}

var _ http.ResponseWriter = (*ResponseWriter)(nil)

func (w *ResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func (w *ResponseWriter) WriteHeader(code int) {
	// w.w.WriteHeader(code)
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

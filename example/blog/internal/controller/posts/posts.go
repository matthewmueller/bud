package posts

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/livebud/buddy/controller"
	"github.com/livebud/buddy/request"
	"github.com/livebud/buddy/router"
	"github.com/livebud/buddy/view"
	"github.com/livebud/buddy/web"
)

func New(db *sql.DB, view view.Interface) *Controller {
	return &Controller{db, view}
}

type Controller struct {
	db   *sql.DB
	view view.Interface
}

func (c *Controller) Mount(r router.Interface) {
	// router := web.NewRouter()
	// router.Get("/posts").Action(c.Index)

	r.Get("/posts", controller.Action(c.view, c.Index5))
	r.Get("/posts", web.Action2(c.view, c.Index2))
	r.Get("/posts", web.Action(c.Index))
	r.Get("/posts/:id", web.Action(c.Show))
	r.Get("/posts/new", http.HandlerFunc(c.New))
	// r.Get("/posts/new", web.Handler(c.Another))
}

type IndexIn struct {
	Title string `json:"title"`
}

type IndexOut struct {
}

func (c *Controller) Index5(ctx *controller.Context[IndexIn, IndexOut]) error {
	fmt.Println(ctx.Map.Title)
	return ctx.Render("posts/index", IndexOut{})
}

func (c *Controller) Index2(ctx context.Context, in *IndexIn) (out string, err error) {
	return "", nil
}

type CreatePost struct {
}

type Post struct {
}

func (c *Controller) Create(ctx context.Context, in *CreatePost) (out *Post, err error) {

	return nil, err
}

func (c *Controller) Index(req *web.Request[IndexIn], res web.Response[IndexOut]) error {
	return c.view.Render(res, "posts/index", &IndexOut{})
	// return c.viewer.Render(req, "posts/index", view.Props{
	// 	"title": req.Params.Title,
	// })
	// return res.Render(IndexOut{})
}

type ShowIn struct {
	ID string `json:"id"`
}

type ShowOut struct {
}

func (c *Controller) Show(req *web.Request[ShowIn], res web.Response[ShowOut]) error {
	return nil
}

func (c *Controller) New(w http.ResponseWriter, r *http.Request) {
	params, err := request.Query[struct {
		Title string `json:"title"`
	}](r)
	if err != nil {
		return
	}
	c.view.Render(w, "posts/new", params)
}

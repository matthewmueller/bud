package web

// type FS = fs.FS

// type Viewer interface {
// 	Mount(fsys fs.FS) error
// }

// type View interface {
// 	Render(ctx context.Context, w io.Writer, key string, props any) error
// }

// func provideViewer(in di.Injector) (Viewer, error) {
// 	return newViewer(), nil
// }

// type webViewer struct {
// }

// func provideView(in di.Injector) (View, error) {
// 	viewer, err := di.Load[Viewer](in)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return newView(), nil
// }

// type webView struct {
// }

// func Viewer(fsys fs.FS) (View, error) {
// 	pages, err := view.Find(fsys)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &viewer{fsys, pages}, nil
// }

// type viewer struct {
// 	fsys  fs.FS
// 	pages map[string]*view.Page
// }

// func (v *viewer) Render(ctx context.Context, w io.Writer, key string, props any) error {
// 	// page, ok := v.pages[key]
// 	// if !ok {
// 	// 	// return view.ErrNotFound
// 	// }
// 	return nil
// 	// return page.Render(ctx, w, props)
// }

// func (v *viewer) Render(ctx context.Context, w io.Writer, key string, props any) error {
// 	page, ok := v.pages[key]
// 	if !ok {
// 		// return view.ErrNotFound
// 	}
// 	// return page.Render(ctx, w, props)
// }

// func New()

// type View interface {
// 	Render()
// }

// type Viewer interface {
// }

// func View(fsys fs.FS, path string) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// http.ServeFile(w, r, path)
// 	})
// }

// func provideView(in di.Injector) (View, error) {
// 	return newView(), nil
// }

// type View interface {
// 	Mount(fsys fs.FS) error
// }

// func newView() *view {
// 	return &view{}
// }

// type view struct {
// }

// func (v *view) Mount(fsys fs.FS) error {
// 	return nil
// }

// func (v *view)

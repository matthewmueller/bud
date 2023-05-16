package di

// type contextKey string

// var key contextKey = "di"

// // New returns a middleware that injects a new injector into the request context.
// // TODO: inherit from the parent injector
// func Middleware(parent Injector) middleware.Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			in := New()
// 			r = r.WithContext(context.WithValue(r.Context(), key, in))
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // From returns the injector from a context.
// func From(ctx context.Context) (Injector, bool) {
// 	in, ok := ctx.Value(key).(Injector)
// 	return in, ok
// }

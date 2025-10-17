# Malcolm - tiny wrapper over net/http to handle middleware
No fancy extensions, no api changes, just net/http with middleware

My goal was to keep it dependency free and under 100sloc, it doesn't even come close to that

## What does middleware look like?
```go
// middleware has the following type def
type Middleware func(http.HandlerFunc) http.HandlerFunc

// so a very basic middleware for athentication may look something like this
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !verifyAuthHeader(r.Header.Get("Authorization")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
```

## The base router
```go

// Create with http.DefaultServerMuk (middleware is optional)
r := malcolm.NewDefaultRouter(middleware1, midleware2 ...)

// Alternatively you can pass in a custom server mux (middleware is optional)
mux := http.NewServerMux()
r := malcolm.NewRouter(mux, middleware1, midleware2 ...)

// middleware assigned at the router level will be assigned to all handlers
r.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
	...
})
```

## Route Groups
You can create groups to apply a base path and middleware chain to any sub routes
```go
r := malcolm.NewDefaultRouter()

// when creating a group it will copy the middleware chain from the parent route and add any provided
// handlers on to the end of the chain
private := r.Group("/private", AuthMiddleware)

// defined routes will use the group prefix as well as the middleware chain assigned to it
// in this case the final path will be "/private/me" and it will have the AuthMiddleware assigned
private.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
	...
})

// this also works with method syntax
// the final path will be "POST /private/me/change-email"
private.HandleFunc("POST /me/change-email", func(w http.ResponseWriter, r *http.Request) {
	...
})
```

## Assign middleware at handler level
```go
// Middleware can also be assigned at the handler level, any parent middleware from router/groups will still apply first
// middleware added at this level will be added onto the end of the chain
func homeHandler(w http.ResponseWriter, r *http.Request) {}

private.HandleFunc("/home", handler, AuthMiddleware)

// The same applies to the Handle method
type UserHandler struct {}
func (h UserHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {}

private.HandleFunc("/user", UserHandler{}, AuthMiddleware)
```

package routing

import (
	"net/http"
	"time"

	"forum/handlers"
	"forum/middlewares"
)

func RegisterRoutes() {
	http.HandleFunc(
		"/posts/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreatePost, true),
			2*time.Second,
		),
	)

	http.HandleFunc(
		"/api/posts/{id}/{endpoint}",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.PostResolver, true),
			2*time.Second,
		),
	)

	http.HandleFunc(
		"/comments/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreateComment, true),
			2*time.Second,
		),
	)

	http.HandleFunc(
		"/api/comments/{id}/{endpoint}",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CommentResolver, true),
			2*time.Second,
		),
	)

	http.HandleFunc(
		"/login",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Login, false),
			2*time.Second,
		),
	)

	http.HandleFunc(
		"/register",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Register, false),
			2*time.Second,
		),
	)

	http.HandleFunc(
		"/logout",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Logout, true),
			2*time.Second,
		),
	)
}

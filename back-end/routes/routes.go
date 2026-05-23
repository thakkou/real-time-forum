package routes

import (
	"net/http"
	"time"

	"forum/handlers"
	"forum/middlewares"
)

func RegisterRoutes() {
	// http.HandleFunc(
	// 	"/api/posts/create",
	// 	middlewares.RateLimit(
	// 		middlewares.CheckSessionCookie(handlers.CreatePost, true),
	// 		3*time.Second,
	// 	),
	// )

	// http.HandleFunc(
	// 	"/api/comments/create",
	// 	middlewares.RateLimit(
	// 		middlewares.CheckSessionCookie(handlers.CreateComment, true),
	// 		3*time.Second,
	// 	),
	// )

	// brute force targets
	http.HandleFunc(
		"/login",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Login, false),
			2*time.Second,
		),
	)
	http.HandleFunc(
		"/logout",
		middlewares.CheckSessionCookie(handlers.Logout, true),
	)

	http.HandleFunc(
		"/register",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Register, false),
			2*time.Second,
		),
	)

	// http.HandleFunc(
	// 	"/api/posts/{id}/{endpoint}",
	// 	middlewares.CheckSessionCookie(handlers.PostResolver, true),
	// )

	// http.HandleFunc(
	// 	"/api/comments/{id}/{endpoint}",
	// 	middlewares.CheckSessionCookie(handlers.CommentResolver, true),
	// )

	// http.HandleFunc(
	// 	"/auth/{provider}",
	// 	middlewares.CheckSessionCookie(handlers.OAuthLogin, false),
	// )

	// http.HandleFunc(
	// 	"/auth/{provider}/callback",
	// 	middlewares.CheckSessionCookie(handlers.OAuthCallback, false),
	// )
}

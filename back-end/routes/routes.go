package routes

import (
	"net/http"
	"time"

	"forum/handlers"
	"forum/middlewares"
)

func RegisterRoutes() {
	// authentification

	http.HandleFunc(
		"/api/login",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Login, false),
			2*time.Second,
		),
	)

	http.HandleFunc(
		"/api/logout",
		middlewares.CheckSessionCookie(handlers.Logout, true),
	)

	http.HandleFunc(
		"/api/register",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Register, false),
			2*time.Second,
		),
	)

	// auth providers

	// http.HandleFunc(
	// 	"/api/auth/{provider}",
	// 	middlewares.CheckSessionCookie(handlers.OAuthLogin, false),
	// )

	// http.HandleFunc(
	// 	"/api/auth/{provider}/callback",
	// 	middlewares.CheckSessionCookie(handlers.OAuthCallback, false),
	// )

	// postes

	http.HandleFunc(
		"/api/posts/getPosts",
		middlewares.RateLimit(
			handlers.GetPosts,
			3*time.Second,
		),
	)

	http.HandleFunc(
		"/api/posts/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreatePost, true),
			3*time.Second,
		),
	)
	
	//this resolver for liking,disliking,delete
	http.HandleFunc(
		"/api/posts/{id}/{endpoint}",
		middlewares.CheckSessionCookie(handlers.PostResolver, true),
	)

	// comments

	http.HandleFunc(
		"/api/comments/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreateComment, true),
			3*time.Second,
		),
	)

	//this resolver for liking,disliking,delete
	http.HandleFunc(
		"/api/comments/{id}/{endpoint}",
		middlewares.CheckSessionCookie(handlers.CommentResolver, true),
	)
}

package routing

import (
	"net/http"
	"time"

	"forum/handlers"
	"forum/middlewares"
)

func RegisterRoutes() {
    http.HandleFunc(
        "/",
        middlewares.RateLimit(
            handlers.Forum,
            200*time.Millisecond, // fast enough for normal browsing
        ),
    )

    http.HandleFunc(
        "/posts/create",
        middlewares.RateLimit(
            middlewares.CheckSessionCookie(handlers.CreatePost, true),
            2*time.Second, // writing a post takes time
        ),
    )

    http.HandleFunc(
        "/api/posts/{id}/{endpoint}",
        middlewares.RateLimit(
            middlewares.CheckSessionCookie(handlers.PostResolver, true),
            500*time.Millisecond, // like/dislike — 1s is fair
        ),
    )

    http.HandleFunc(
        "/comments/create",
        middlewares.RateLimit(
            middlewares.CheckSessionCookie(handlers.CreateComment, true),
            2*time.Second, // same as post
        ),
    )

    http.HandleFunc(
        "/api/comments/{id}/{endpoint}",
        middlewares.RateLimit(
            middlewares.CheckSessionCookie(handlers.CommentResolver, true),
            1*time.Second, // same as post reactions
        ),
    )

    http.HandleFunc(
        "/login",
        middlewares.RateLimit(
            middlewares.CheckSessionCookie(handlers.Login, false),
            2*time.Second, // brute force protection
        ),
    )

    http.HandleFunc(
        "/register",
        middlewares.RateLimit(
            middlewares.CheckSessionCookie(handlers.Register, false),
            2*time.Second, // spam protection
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
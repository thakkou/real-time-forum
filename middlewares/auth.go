package middlewares

import (
	"database/sql"
	"net/http"
	"time"

	"forum/database"
	"forum/handlers"
	"forum/models"
)

// CheckSessionCookie validates session cookie and redirects depending on requiresAuth
// true: verifies the user's session cookie before allowing access to protected routes.
// false: prevents already-logged-in users from accessing login/register pages.
func CheckSessionCookie(handler http.HandlerFunc, requiresAuth bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		// 1. cookie existence in request
		if err == nil && cookie.Value != "" {
			var expiryTime time.Time
			err = database.Database.QueryRow(
				"SELECT expires_at FROM sessions WHERE id = ?", cookie.Value,
			).Scan(&expiryTime)

			switch err {
			case nil:
				// session found
				if expiryTime.After(time.Now()) {
					if requiresAuth {
						handler(w, r)
					} else {
						http.Redirect(w, r, "/", http.StatusSeeOther)
					}
					return
				}

				// expired
				models.DeleteSession(cookie.Value)
				clearSessionCookie(w)

			case sql.ErrNoRows:
				// session not found
				clearSessionCookie(w)

			default:
				handlers.HandleError(w, http.StatusInternalServerError, "Database error")
				return
			}
		}
		if requiresAuth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			handler(w, r)
		}
	}
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

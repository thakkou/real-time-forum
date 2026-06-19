package middlewares

import (
	"database/sql"
	"net/http"
	"time"

	"forum/database"
	"forum/utilities"
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
				utilities.DeleteSession(cookie.Value)
				utilities.ClearSessionCookie(w)

			case sql.ErrNoRows:
				// session not found
				utilities.ClearSessionCookie(w)

			default:
				utilities.WriteJSON(w, http.StatusInternalServerError, "database error", nil)
				return
			}
		}
		if requiresAuth {
			utilities.WriteJSON(w, http.StatusSeeOther, "user should be login ", nil)
			return
		} else {
			handler(w, r)
		}
	}
}

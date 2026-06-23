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
		// No session cookie
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			if requiresAuth {
				utilities.WriteJSON(w, http.StatusUnauthorized, "login required", nil)
				return
			}

			handler(w, r)
			return
		}

		// Check session in database
		var expiryTime time.Time
		err = database.Database.QueryRow(
			"SELECT expires_at FROM sessions WHERE id = ?",
			cookie.Value,
		).Scan(&expiryTime)

		switch err {
		case nil:
			// continue below

		case sql.ErrNoRows:
			utilities.ClearSessionCookie(w)

			if requiresAuth {
				utilities.WriteJSON(w, http.StatusUnauthorized, "login required", nil)
				return
			}

			handler(w, r)
			return

		default:
			utilities.WriteJSON(w, http.StatusInternalServerError, "database error", nil)
			return
		}

		// Session expired
		if expiryTime.Before(time.Now()) {
			utilities.DeleteSession(cookie.Value)
			utilities.ClearSessionCookie(w)

			if requiresAuth {
				utilities.WriteJSON(w, http.StatusUnauthorized, "session expired", nil)
				return
			}

			handler(w, r)
			return
		}

		// Valid session
		if requiresAuth {
			handler(w, r)
			return
		}

		// Guest-only route (/login, /register, ...)
		utilities.WriteJSON(
			w,
			http.StatusConflict,
			"Unauthorized",
			nil,
		)
		return
	}
}

package utilities

import (
	"fmt"
	"time"
)

func TimeAgo(t time.Time) string {
	d := time.Since(t)

	if d < time.Minute {
		return fmt.Sprintf("%d seconds ago", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	}
	if d < 30*24*time.Hour {
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	}
	return fmt.Sprintf("%d months ago", int(d.Hours()/(24*30)))
}

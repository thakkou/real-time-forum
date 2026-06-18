package models

import (
	"time"
)

type Post struct {
	Id                      int
	UserId                  int
	Nickname                string
	Created_at              time.Time
	TimeAgo                 string
	Title                   string
	Text                    string
	LikeCount, DislikeCount int
	IsLiked                 int // 1:liked, 0:none, -1:disliked
	Comments                []Comment
	Categories              []string
	Image                   string
}

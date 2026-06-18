package models

import (
	"time"
)

type Comment struct {
	Id                      int
	UserId                  int
	Nickname                string
	Created_at              time.Time
	TimeAgo                 string
	Text                    string
	LikeCount, DislikeCount int
	IsLiked                 int // 1:liked, 0:none, -1:disliked
}

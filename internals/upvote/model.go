package upvote

import "time"

type Upvote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BugID     uint      `json:"bug_id"`
	IPHash    string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}
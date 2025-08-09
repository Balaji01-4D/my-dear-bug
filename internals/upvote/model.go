package upvote

import "time"

type Upvote struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ConfessionID uint      `json:"confessionId"`
	IPHash       string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

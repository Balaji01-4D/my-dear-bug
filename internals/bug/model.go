package bug

import (
	"time"
)

type Bug struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Language    string    `json:"language"`
	Snippet     string    `gorm:"type:text" json:"snippet"`
	CreatedAt   time.Time `json:"createdAt"`
	Upvotes     int       `json:"upvotes"`
}

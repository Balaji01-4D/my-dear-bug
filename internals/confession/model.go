package confession

import (
	"time"

	"github.com/Balaji01-4D/shit-happens/internals/tag"
)

type Confession struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Language    string    `gorm:"size:50" json:"language"`
	Snippet     string    `gorm:"type:text" json:"snippet"`
	Tags        []tag.Tag `gorm:"many2many:confession_tags;" json:"tags"`
	Sentiment   string    `gorm:"size:20" json:"sentiment"` // e.g., "positive", "negative", "neutral"
	IsFlagged   bool      `gorm:"default:false" json:"isFlagged"`
	CreatedAt   time.Time `json:"createdAt"`
	Upvotes     int       `gorm:"default:0" json:"upvotes"`
}

package tag

type Tag struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"uniqueIndex;size:20" json:"name"`
}


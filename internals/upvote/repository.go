package upvote

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// Checks whether the user is already upvoted the bug
func (r *Repository) HasUpvoted(bugID uint, ipHash string) bool {
	var upvote Upvote
	err := r.DB.Where("bug_id = ?", bugID).Where("ip_hash = ?", ipHash).First(&upvote).Error
	return err == nil
}



func (r *Repository) Save(upvote *Upvote) error {
	return r.DB.Create(upvote).Error
}

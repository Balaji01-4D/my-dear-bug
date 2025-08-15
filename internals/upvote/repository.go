package upvote

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// Checks whether the user is already upvoted the confessions by IP or client hash
func (r *Repository) HasUpvoted(confessionID uint, ipHash, clientHash string) bool {
	var upvote Upvote
	q := r.DB.Where("confession_id = ?", confessionID)
	if ipHash != "" && clientHash != "" {
		q = q.Where(r.DB.Where("ip_hash = ?", ipHash).Or("client_hash = ?", clientHash))
	} else if ipHash != "" {
		q = q.Where("ip_hash = ?", ipHash)
	} else if clientHash != "" {
		q = q.Where("client_hash = ?", clientHash)
	}
	err := q.First(&upvote).Error
	return err == nil
}

func (r *Repository) Save(upvote *Upvote) error {
	// Rely on unique indexes; if a duplicate insert happens, surface the error to caller
	return r.DB.Create(upvote).Error
}

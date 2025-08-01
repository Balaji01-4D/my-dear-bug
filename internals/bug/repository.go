package bug

import "gorm.io/gorm"

type Repo struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo  {
	return &Repo{DB: db}
}

func (r *Repo) Create(bug *Bug) error {
	return  r.DB.Create(bug).Error
}

func (r *Repo) ListAll(offset, limit int) ([]Bug, error) {
	var out []Bug

	err := r.DB.
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&out).Error

	return out, err
}


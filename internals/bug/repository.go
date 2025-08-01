package bug

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetToBugs(offest int, limit int) ([]Bug, error) {
	var bugs []Bug

	err := r.DB.
	Offset(offest).
	Limit(limit).
	Order("upvoteS DESC").
	Find(&bugs).Error

	return bugs, err

}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(bug *Bug) error {
	return r.DB.Create(bug).Error
}

func (r *Repository) ListAll(offset, limit int) ([]Bug, error) {
	var out []Bug

	err := r.DB.
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&out).Error

	return out, err
}

func (r *Repository) Get(id uint) (Bug, error) {
	var bug Bug

	err := r.DB.Find(&bug, id).Error

	return bug, err
}

func (r *Repository) Delete(id uint) error {
	err := r.DB.Delete(&Bug{}, id).Error

	return err
}

func (r *Repository) GetByLanguage(language string, offset int, limit int) ([]Bug, error) {

	var bugs []Bug
	err := r.DB.
		Offset(offset).
		Limit(limit).
		Where("language ILIKE ?", language).
		Order("created_at DESC").
		Find(&bugs).Error

	return bugs, err
}

package confession

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetTopConfessions(offest int, limit int) ([]Confession, error) {
	var confessions []Confession

	err := r.DB.
		Offset(offest).
		Limit(limit).
		Order("upvoteS DESC").
		Find(&confessions).Error

	return confessions, err

}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(confession *Confession) error {
	return r.DB.Create(confession).Error
}

func (r *Repository) List(offset, limit int) ([]Confession, error) {
	var out []Confession

	err := r.DB.
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&out).Error

	return out, err
}

func (r *Repository) Get(id uint) (Confession, error) {
	var confession Confession

	err := r.DB.Find(&confession, id).Error

	return confession, err
}

func (r *Repository) Delete(id uint) error {
	err := r.DB.Delete(&Confession{}, id).Error

	return err
}

func (r *Repository) GetByLanguage(language string, offset int, limit int) ([]Confession, error) {

	var confessions []Confession
	err := r.DB.
		Offset(offset).
		Limit(limit).
		Where("language ILIKE ?", language).
		Order("created_at DESC").
		Find(&confessions).Error

	return confessions, err
}

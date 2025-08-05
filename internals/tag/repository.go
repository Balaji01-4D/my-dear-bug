package tag

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Save(tag *Tag) error {
	return r.DB.Create(tag).Error
}

func (r *Repository) GetTags() ([]Tag, error) {
	var tags []Tag
	err := r.DB.Find(&tags).Error
	return tags, err
}


func (r *Repository) SuggestTags(query string) ([]Tag, error) {

	var tags []Tag
	err := r.DB.
		Where("name ILIKE ?", query+"%").
		Order("name ASC").
		Limit(6).
		Find(&tags).Error;

	return  tags, err

}


func (r *Repository) DeleteTags(id int) error {
	return r.DB.Delete(&Tag{}, id).Error
}

func (r *Repository) GetTagByName(name string) (Tag, error) {
	var tag Tag
	err := r.DB.Where("name = ?", name).First(&tag).Error
	return tag, err
}
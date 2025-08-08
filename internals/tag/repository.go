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
	tx := r.DB.Begin()

	if err := tx.Error; err != nil {
		return err
	}

		// Removing join rows first
	if err := tx.Exec("DELETE FROM confession_tags WHERE tag_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Now delete the tag
	if err := tx.Delete(&Tag{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *Repository) GetTagByName(name string) (Tag, error) {
	var tag Tag
	err := r.DB.Where("name = ?", name).First(&tag).Error
	return tag, err
}
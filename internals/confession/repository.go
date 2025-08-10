package confession

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetTopConfessions(offset int, limit int) ([]Confession, error) { // renamed param
	var confessions []Confession

	err := r.DB.
		Preload("Tags").
		Offset(offset).
		Limit(limit).
		Order("upvotes DESC").
		Find(&confessions).Error

	return confessions, err
}

// GetTopConfessionsSince returns top confessions since a given time (weekly/monthly trending)
func (r *Repository) GetTopConfessionsSince(since time.Time, offset int, limit int) ([]Confession, error) {
	var confessions []Confession
	err := r.DB.
		Preload("Tags").
		Where("created_at >= ?", since).
		Offset(offset).
		Limit(limit).
		Order("upvotes DESC").
		Find(&confessions).Error
	return confessions, err
}

// HallOfFame returns all‑time top confessions (larger limit by caller) - could add thresholds later
func (r *Repository) HallOfFame(offset, limit int) ([]Confession, error) {
	var confessions []Confession
	err := r.DB.
		Preload("Tags").
		Offset(offset).
		Limit(limit).
		Order("upvotes DESC").
		Find(&confessions).Error
	return confessions, err
}

// RandomConfession returns a single random confession (PostgreSQL RANDOM())
func (r *Repository) RandomConfession() (Confession, error) {
	var c Confession
	err := r.DB.Preload("Tags").Order("RANDOM()").Limit(1).First(&c).Error
	return c, err
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
		Preload("Tags").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&out).Error

	return out, err
}

func (r *Repository) Get(id uint) (Confession, error) {
	var confession Confession

	err := r.DB.Preload("Tags").First(&confession, id).Error

	return confession, err
}

func (r *Repository) Delete(id uint) error {
	tx := r.DB.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	// Removing join rows first (ignore if none)
	if err := tx.Exec("DELETE FROM confession_tags WHERE confession_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	res := tx.Delete(&Confession{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	if res.RowsAffected == 0 { // not found
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	return tx.Commit().Error
}

func (r *Repository) GetByLanguage(language string, offset int, limit int) ([]Confession, error) {
	var confessions []Confession
	err := r.DB.
		Preload("Tags").
		Offset(offset).
		Limit(limit).
		Where("language ILIKE ?", language).
		Order("created_at DESC").
		Find(&confessions).Error

	return confessions, err
}

func (r *Repository) Search(q, language, tag string, offset, limit int) ([]Confession, error) {
	var confessions []Confession
	db := r.DB.Model(&Confession{}).Preload("Tags")

	if tag != "" {

		db = db.Joins("JOIN confession_tags ct ON ct.confession_id = confessions.id").
			Joins("JOIN tags ON tags.id = ct.tag_id").
			Where("tags.name ILIKE ?", tag)
	}

	if language != "" {
		db = db.Where("language ILIKE ?", language)
	}

	if q != "" {
		like := "%" + q + "%"
		db = db.Where("(title ILIKE ? OR description ILIKE ? OR snippet ILIKE ?)", like, like, like)
	}

	err := db.Distinct().
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&confessions).Error
	return confessions, err
}


func isNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }

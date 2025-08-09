package confession

import (
	"errors"
	"strings"
	"time"

	"github.com/Balaji01-4D/my-dear-bug/internals/tag"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// used to create the confessions from the dto and save to database
func (s *Service) Create(dto ConfessionRequest) (Confession, error) {

	confession := Confession{
		Title:       dto.Title,
		Description: dto.Description,
		Language:    dto.Language,
		Snippet:     dto.Snippet,
		Sentiment:   "happy", // hardcoded for just now
		IsFlagged:   dto.IsFlagged,
		CreatedAt:   now(),
		Upvotes:     0,
	}

	var tags []tag.Tag

	// Deduplicate incoming tags
	seen := make(map[string]struct{}, len(dto.Tags))
	for _, tagName := range dto.Tags {
		tagName = strings.TrimSpace(strings.ToLower(tagName))
		if tagName == "" {
			continue
		}
		if _, ok := seen[tagName]; ok {
			continue
		}
		seen[tagName] = struct{}{}
	}

	for tagName := range seen {
		var t tag.Tag
		// Try to find existing
		if err := s.repo.DB.Where("name = ?", tagName).First(&t).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create if absent, ignoring conflict (atomic)
				if err := s.repo.DB.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "name"}},
					DoNothing: true,
				}).Create(&tag.Tag{Name: tagName}).Error; err != nil {
					return Confession{}, err
				}
				// Fetch the row (handles both created and conflicted cases)
				if err := s.repo.DB.Where("name = ?", tagName).First(&t).Error; err != nil {
					return Confession{}, err
				}
			} else {
				return Confession{}, err
			}
		}
		tags = append(tags, t)
	}
	confession.Tags = tags

	err := s.repo.Create(&confession)
	return confession, err
}

// list confessions based on the offset and limit
func (s *Service) List(offset, limit int) ([]Confession, error) {
	return s.repo.List(offset, limit)
}

// get confession by its id (primary key)
func (s *Service) Get(id uint) (Confession, error) {
	return s.repo.Get(id)
}

// Delete the confession by its id
func (s *Service) Delete(id uint) error {
	return s.repo.Delete(id)
}

// Return the confessions based on the language
func (s *Service) GetByLanguage(language string, offset int, limit int) ([]Confession, error) {
	return s.repo.GetByLanguage(language, offset, limit)
}

func (s *Service) GetTopConfessions(offset int, limit int) ([]Confession, error) {
	return s.repo.GetTopConfessions(offset, limit)
}

func (s *Service) TrendingWeekly(offset, limit int) ([]Confession, error) {
	weekAgo := now().AddDate(0, 0, -7)
	return s.repo.GetTopConfessionsSince(weekAgo, offset, limit)
}

func (s *Service) TrendingMonthly(offset, limit int) ([]Confession, error) {
	monthAgo := now().AddDate(0, -1, 0)
	return s.repo.GetTopConfessionsSince(monthAgo, offset, limit)
}

func (s *Service) HallOfFame(offset, limit int) ([]Confession, error) {
	return s.repo.HallOfFame(offset, limit)
}

func (s *Service) Random() (Confession, error) {
	return s.repo.RandomConfession()
}

// Search confessions by free text / language / tag
func (s *Service) Search(q, language, tag string, offset, limit int) ([]Confession, error) {
	return s.repo.Search(q, language, tag, offset, limit)
}

func now() time.Time {
	return time.Now()
}

package confession

import (
	"strings"
	"time"

	"github.com/Balaji01-4D/my-dear-bug/internals/tag"
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

	for _, tagName := range dto.Tags {
		tagName = strings.TrimSpace(strings.ToLower(tagName))
		if tagName == "" {
			continue
		}
		
		tagService := tag.NewService((*tag.Repository)(s.repo)) // Create an instance of tag.Service
		resultTag, err := tagService.GetTagByName(tagName)
		if err != nil {
			if err := s.repo.DB.Create(&resultTag).Error; err != nil {
				return Confession{}, err
			}
		}
		tags = append(tags, resultTag)
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

func (s *Service) GetTopConfessions(offest int, limit int) ([]Confession, error) {
	return s.repo.GetTopConfessions(offest, limit)
}

func now() time.Time {
	return time.Now()
}

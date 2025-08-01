package bug

import "time"

type Service struct {
	repo *Repository
}


func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(dto CreateBugDTO) (Bug, error) {

	bug := Bug{
		Title:       dto.Title,
		Description: dto.Description,
		Language:    dto.Language,
		Upvotes:     0,
		Snippet:     dto.Snippet,
		CreatedAt:   now(),
	}
	err := s.repo.Create(&bug)
	return bug, err
}

func (s *Service) List(offset, limit int) ([]Bug, error) {
	return s.repo.ListAll(offset, limit)
}

func (s *Service) Get(id uint) (Bug, error) {
	return s.repo.Get(id)
}

func (s *Service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func now() time.Time {
	return time.Now()
}

func (s *Service) GetByLanguage(language string, offset int, limit int) ([]Bug, error) {
	return s.repo.GetByLanguage(language, offset, limit)
}

func (s *Service) GetTopBugs(offest int, limit int) ([]Bug, error) {
	return s.repo.GetToBugs(offest, limit)
}
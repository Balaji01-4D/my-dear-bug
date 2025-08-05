package tag

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateTag(name string) error {

	tag := Tag{
		Name: name,
	}

	return  s.repo.Save(&tag)
}


func (s *Service) SuggestTags(query string) ([]Tag, error) {
	return s.repo.SuggestTags(query)
}

func (s *Service) DeleteTags(id int) error {
	return s.repo.DeleteTags(id)
}


func (s *Service) GetTags() ([]Tag, error) {
	return s.repo.GetTags()
}

func (s *Service) GetTagByName(name string) (Tag, error) {
	return s.repo.GetTagByName(name)
}
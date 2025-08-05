package upvote

import (
	"time"

	confession "github.com/Balaji01-4D/my-dear-bug/internals/confession"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Upvote(confessionID uint, ipHash string) error {

	err := s.repo.Save(&Upvote{
		ConfessionID: confessionID,
		IPHash:       ipHash,
		CreatedAt:    time.Now(),
	})

	if err != nil {
		return err
	}

	if err = s.repo.DB.Model(&confession.Confession{}).
		Where("id = ?", confessionID).
		UpdateColumn("upvotes", gorm.Expr("upvotes + 1")).Error; err != nil {
		return err
	}

	return nil

}

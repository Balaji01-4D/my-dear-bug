package upvote

import (
	"time"

	"github.com/Balaji01-4D/my-dear-bug/internals/bug"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Upvote(bugId uint, ipHash string) error {

	err := s.repo.Save(&Upvote{
		BugID:     bugId,
		IPHash:    ipHash,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	if err = s.repo.DB.Model(&bug.Bug{}).
		Where("id = ?", bugId).
		UpdateColumn("upvotes", gorm.Expr("upvotes + 1")).Error; err != nil {
		return err
	}

	return nil

}

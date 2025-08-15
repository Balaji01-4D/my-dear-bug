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

func (s *Service) Upvote(confessionID uint, ipHash, clientHash string) error {
	now := time.Now()
	up := &Upvote{ConfessionID: confessionID, IPHash: ipHash, ClientHash: clientHash, CreatedAt: now}
	if err := s.repo.Save(up); err != nil {
		// If insert fails (likely due to unique constraint), do not bump counter
		return err
	}
	// Only bump the confession's upvote count after successful insert
	if err := s.repo.DB.Model(&confession.Confession{}).
		Where("id = ?", confessionID).
		UpdateColumn("upvotes", gorm.Expr("upvotes + 1")).Error; err != nil {
		return err
	}
	return nil
}

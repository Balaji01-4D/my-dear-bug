package upvote

type UpvoteDTO struct {
	ConfessionID uint `json:"confessionId" binding:"required"`
}

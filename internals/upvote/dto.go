package upvote

type UpvoteDTO struct {
	BugID uint `json:"bug_id" binding:"required"`
}

package bug

type CreateBugDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Language    string `json:"language" binding:"required"`
}

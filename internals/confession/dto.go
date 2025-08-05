package confession

type ConfessionRequest struct {
	Title       string   `json:"title" binding:"required,min=5,max=100"`
	Description string   `json:"description" binding:"required,min=10"`
	Snippet     string   `json:"snippet" binding:"omitempty"`
	Language    string   `json:"language" binding:"required"`
	Tags        []string `json:"tags" binding:"omitempty,dive,min=1"`
	IsFlagged   bool     `json:"isFlagged" binding:"required"`
}

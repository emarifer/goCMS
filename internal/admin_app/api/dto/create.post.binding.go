package dto

type AddPostRequest struct {
	Title   string `json:"title" validate:"required"`
	Excerpt string `json:"excerpt" validate:"required"`
	Content string `json:"content" validate:"required"`
}

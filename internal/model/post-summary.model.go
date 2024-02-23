package model

type PostSummary struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
}

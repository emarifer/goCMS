package dto

type PostBinding struct {
	Id string `uri:"id" binding:"required"`
}

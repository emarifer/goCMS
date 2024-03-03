package dto

type ImageBinding struct {
	UUID string `uri:"uuid" binding:"required"`
}

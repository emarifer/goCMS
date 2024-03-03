package entity

import "github.com/google/uuid"

type Image struct {
	UUID uuid.UUID `db:"uuid"`
	Name string    `db:"name"`
	Alt  string    `db:"alt"`
}

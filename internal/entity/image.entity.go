package entity

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	UUID      uuid.UUID `db:"uuid"`
	Name      string    `db:"name"`
	Alt       string    `db:"alt"`
	CreatedAt time.Time `db:"created_at"`
}

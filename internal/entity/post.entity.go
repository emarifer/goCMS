package entity

import "time"

type Post struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Excerpt   string    `db:"excerpt"`
	CreatedAt time.Time `db:"created_at"`
}

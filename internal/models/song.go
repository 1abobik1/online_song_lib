package models

import "time"

type Song struct {
	ID          int       `db:"id"`
	GroupName   string    `db:"group_name"`
	SongName    string    `db:"song_name"`
	ReleaseDate time.Time `db:"release_date"`
	Text        string    `db:"text"`
	Link        string    `db:"link"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

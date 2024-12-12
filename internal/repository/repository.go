package repository

import (
	"context"

	"github.com/1abobik1/online_song_lib/internal/models"
	"github.com/1abobik1/online_song_lib/internal/transport/http/dto"
)

type SongFilter struct {
	GroupName   string
	SongName    string
	ReleaseDate string
	Limit       int
	Offset      int
}

type SongUpdate struct {
	GroupName   *string
	SongName    *string
	ReleaseDate *string
	Text        *string
	Link        *string
}

type SongRepository interface {
	GetSongs(ctx context.Context, filter SongFilter) ([]models.Song, error)
	GetSongByID(ctx context.Context, id int) (*models.Song, error)
	DeleteSongByID(ctx context.Context, id int) error
	UpdateSong(ctx context.Context, id int, update SongUpdate) error
	CreateSong(ctx context.Context, song *dto.SongResponse) (int, error)
}

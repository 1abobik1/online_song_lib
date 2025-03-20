package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/1abobik1/online_song_lib/internal/models"
	"github.com/1abobik1/online_song_lib/internal/repository"
	"github.com/1abobik1/online_song_lib/internal/transport/http/dto"

	_ "github.com/lib/pq"
)

type Storage struct {
	db     *sql.DB
	logger *slog.Logger
}

// New создает новое подключение к базе данных PostgreSQL.
func New(storagePath string, logger *slog.Logger) (*Storage, error) {
	const op = "storage.postresql.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		logger.Warn("failed to open database connection", "error", err)
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	logger.Debug("database connected", "db_url", storagePath)
	return &Storage{db: db, logger: logger}, nil
}

// Stop закрывает подключение к базе данных.
func (s *Storage) Stop() error {
	s.logger.Debug("closing database connection")
	return s.db.Close()
}

func (s *Storage) GetSongs(ctx context.Context, filter repository.SongFilter) ([]models.Song, error) {
	var args []interface{}
	var where []string

	baseQuery := `SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at FROM songs`

	if filter.GroupName != "" {
		args = append(args, filter.GroupName)
		where = append(where, fmt.Sprintf("group_name ILIKE $%d", len(args)))
	}

	if filter.SongName != "" {
		args = append(args, filter.SongName)
		where = append(where, fmt.Sprintf("song_name ILIKE $%d", len(args)))
	}

	if filter.ReleaseDate != "" {
		args = append(args, filter.ReleaseDate)
		where = append(where, fmt.Sprintf("release_date = $%d", len(args)))
	}

	if len(where) > 0 {
		baseQuery += " WHERE " + strings.Join(where, " AND ")
	}

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		baseQuery += fmt.Sprintf(" LIMIT $%d", len(args))
	}

	if filter.Offset > 0 {
		args = append(args, filter.Offset)
		baseQuery += fmt.Sprintf(" OFFSET $%d", len(args))
	}

	rows, err := s.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		s.logger.Warn("GetSongs query failed", "error", err)
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link, &song.CreatedAt, &song.UpdatedAt); err != nil {
			s.logger.Warn("failed to scan row", "error", err)
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

// update postgres
func (s *Storage) GetSongByID(ctx context.Context, id int) (*models.Song, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at FROM songs WHERE id = $1`, id)

	var song models.Song
	if err := row.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link, &song.CreatedAt, &song.UpdatedAt); err != nil {
		s.logger.Debug("GetSongByID not found or error", "id", id, "error", err)
		return nil, err
	}

	return &song, nil
}

func (s *Storage) DeleteSongByID(ctx context.Context, id int) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		s.logger.Warn("failed to delete song", "id", id, "error", err)
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no song found with id %d", id)
	}
	return nil
}

func (s *Storage) UpdateSong(ctx context.Context, id int, update repository.SongUpdate) error {
	var sets []string
	var args []interface{}
	argPos := 1

	if update.GroupName != nil {
		sets = append(sets, fmt.Sprintf("group_name=$%d", argPos))
		args = append(args, *update.GroupName)
		argPos++
	}
	if update.SongName != nil {
		sets = append(sets, fmt.Sprintf("song_name=$%d", argPos))
		args = append(args, *update.SongName)
		argPos++
	}
	if update.ReleaseDate != nil {
		sets = append(sets, fmt.Sprintf("release_date=$%d", argPos))
		args = append(args, *update.ReleaseDate)
		argPos++
	}
	if update.Text != nil {
		sets = append(sets, fmt.Sprintf("text=$%d", argPos))
		args = append(args, *update.Text)
		argPos++
	}
	if update.Link != nil {
		sets = append(sets, fmt.Sprintf("link=$%d", argPos))
		args = append(args, *update.Link)
		argPos++
	}

	if len(sets) == 0 {
		// Нечего обновлять
		return nil
	}

	query := fmt.Sprintf("UPDATE songs SET %s, updated_at=NOW() WHERE id=$%d", strings.Join(sets, ", "), argPos)
	args = append(args, id)

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Warn("failed to update song", "id", id, "error", err, "query", query, "args", args)
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no song found with id %d", id)
	}
	return nil
}

func (s *Storage) CreateSong(ctx context.Context, song *dto.SongResponse) (int, error) {
	var id int
	query := `
		INSERT INTO songs (group_name, song_name, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	err := s.db.QueryRowContext(ctx, query,
		song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link,
	).Scan(&id)
	if err != nil {
		s.logger.Warn("failed to create song", "error", err)
		return 0, err
	}

	s.logger.Debug("song created successfully", "id", id)
	song.ID = id
	return id, nil
}

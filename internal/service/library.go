package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/1abobik1/online_song_lib/internal/models"
	"github.com/1abobik1/online_song_lib/internal/repository"
	"github.com/1abobik1/online_song_lib/internal/transport/http/dto"
)

type LibraryService struct {
	repo repository.SongRepository
}

type ExternalSongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func NewLibraryService(repo repository.SongRepository) *LibraryService {
	return &LibraryService{repo: repo}
}

func (s *LibraryService) GetSongs(ctx context.Context, filter repository.SongFilter) ([]models.Song, error) {
	return s.repo.GetSongs(ctx, filter)
}

func (s *LibraryService) DeleteSong(ctx context.Context, id int) error {
	return s.repo.DeleteSongByID(ctx, id)
}

func (s *LibraryService) UpdateSong(ctx context.Context, id int, update repository.SongUpdate) error {
	return s.repo.UpdateSong(ctx, id, update)
}

// Пагинация по куплетам. Допустим, куплеты разделены "\n\n"
func (s *LibraryService) GetSongTextByVerse(ctx context.Context, id, verse int) (string, error) {
	song, err := s.repo.GetSongByID(ctx, id)
	if err != nil {
		return "", err
	}

	verses := strings.Split(song.Text, "\n\n")

	// Если verse = 0, вернуть весь текст
	if verse == 0 {
		return song.Text, nil
	}

	if verse < 1 || verse > len(verses) {
		return "", errors.New("verse not found")
	}

	return verses[verse-1], nil
}

func (s *LibraryService) AddSong(ctx context.Context, group, song string, externalAPIURL string) (*dto.SongResponse, error) {
	// Вызов внешнего API
	detail, err := s.fetchExternalSongDetail(ctx, group, song, externalAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch external details: %v", err)
	}

	// Парсим дату
	releaseDate, err := time.Parse("02.01.2006", detail.ReleaseDate)
	if err != nil {
		return nil, fmt.Errorf("invalid release date format: %v", err)
	}

	newSong := &dto.SongResponse{
		GroupName:   group,
		SongName:    song,
		ReleaseDate: releaseDate,
		Text:        detail.Text,
		Link:        detail.Link,
	}

	if _, err := s.repo.CreateSong(ctx, newSong); err != nil {
		return nil, err
	}

	return newSong, nil
}

func (s *LibraryService) fetchExternalSongDetail(ctx context.Context, group, song, externalAPIURL string) (*ExternalSongDetail, error) {
	// Формируем запрос к внешнему API
	req, err := http.NewRequestWithContext(ctx, "GET", externalAPIURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("group", group)
	q.Add("song", song)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	var detail ExternalSongDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, err
	}

	return &detail, nil
}

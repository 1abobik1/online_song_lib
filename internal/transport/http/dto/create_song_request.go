package dto

type CreateSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}
// update dto
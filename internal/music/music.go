package music

import (
	"github.com/ogabrielrodrugues/moodfy/internal/style"
)

type DTO struct {
	ArtistID    int32   `json:"artist_id"`
	Name        string  `json:"name"`
	CoverImage  string  `json:"cover_image"`
	SpotifyLink string  `json:"spotify_link"`
	Styles      []int32 `json:"styles"`
}

type Music struct {
	ID          int32         `json:"id"`
	ArtistID    int32         `json:"artist_id"`
	Name        string        `json:"name"`
	CoverImage  string        `json:"cover_image"`
	SpotifyLink string        `json:"spotify_link"`
	Styles      []style.Style `json:"styles,omitempty"`
}

func New(artist_id int32, name, cover_image, spotify_link string) *Music {
	return &Music{
		ArtistID:    artist_id,
		Name:        name,
		CoverImage:  cover_image,
		SpotifyLink: spotify_link,
	}
}

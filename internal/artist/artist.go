package artist

import "github.com/ogabrielrodrugues/moodfy/internal/music"

type DTO struct {
	Name string `json:"name"`
}

type Artist struct {
	ID     int32          `json:"id"`
	Name   string         `json:"name"`
	Musics []*music.Music `json:"musics,omitempty"`
}

func New(name string) *Artist {
	return &Artist{
		Name: name,
	}
}

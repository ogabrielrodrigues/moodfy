package music

type DTO struct {
	ArtistID    int32   `json:"artist_id"`
	Name        string  `json:"name"`
	SpotifyLink string  `json:"spotify_link"`
	Styles      []int32 `json:"styles"`
}

type Music struct {
	ID          int32    `json:"id"`
	Name        string   `json:"music"`
	Artist      string   `json:"artist,omitempty"`
	CoverImage  string   `json:"cover_image"`
	SpotifyLink string   `json:"spotify_link"`
	Styles      []string `json:"styles,omitempty"`
}

func New(name, cover_image, spotify_link string) *Music {
	return &Music{
		Name:        name,
		CoverImage:  cover_image,
		SpotifyLink: spotify_link,
	}
}

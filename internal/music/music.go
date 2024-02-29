package music

import "github.com/ogabrielrodrugues/moodfy/internal/style"

type Music struct {
	ID    int
	Name  string
	Style *style.Style
}

func New(name string) *Music {
	return &Music{
		Name: name,
	}
}

func (s *Music) AddStyle(style *style.Style) {
	s.Style = style
}

package style

type DTO struct {
	Name string `json:"name"`
}

type Style struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func New(name string) *Style {
	return &Style{
		Name: name,
	}
}

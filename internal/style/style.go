package style

type Style struct {
	ID   int
	Name string
}

func New(name string) *Style {
	return &Style{
		Name: name,
	}
}

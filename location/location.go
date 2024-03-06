package location

type Location struct {
	FilePath string
	Row      int
	Col      int
}

func New(FilePath string, Row, Col int) Location {
	return Location{
		FilePath: FilePath,
		Row:      Row,
		Col:      Col,
	}
}

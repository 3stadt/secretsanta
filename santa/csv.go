package santa

import (
	"github.com/artonge/go-csv-tag"
)

type participant struct {
	Name  string `csv:"name"`
	Email string `csv:"email"`
}

func Import(path string) (*[]participant, error) {
	participants := []participant{}
	err := csvtag.Load(csvtag.Config{
		Path:      path,
		Dest:      &participants,
		Separator: ';',
	})
	return &participants, err
}
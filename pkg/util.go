package pkg

import (
	"github.com/dacjames/persistsql/internal/util"
)

type Placeholders util.Placeholders

func NewPlaceholders() Placeholders {
	return util.NewPlaceholders()
}

func NewPlaceholdersAt(start int) Placeholders {
	return util.NewPlaceholdersAt(start)
}

package pkg

import (
	"database/sql"

	"github.com/dacjames/persistsql/internal/core"
)

type Collection interface {
	Append(f func(db sql.DB) error) error
}

func NewCollection() {
	core.NewCollection()
}

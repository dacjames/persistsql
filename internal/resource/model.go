package resource

import (
	"time"
)

type Meta struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
	Revision  int       `db:"revision_id"`
}

type Resource struct {
	ID      ID   `db:"resource_id"`
	Meta    Meta `db:"meta"`
	Deleted bool
}

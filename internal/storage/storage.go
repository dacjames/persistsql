package storage

import (
	"github.com/dacjames/persistsql/internal/resource"
	"github.com/dacjames/persistsql/internal/tags"
)

type Stater interface {
	ResourceID() resource.ID
	ResourceTags() tags.Tagset
	State()
}

type Filterer interface {
	Filter()
}

type Getter interface {
	GetAny(id resource.ID) (state Stater, err error)
}
type Putter interface {
	PutAny(state Stater) (err error)
}
type Lister interface {
	ListAny(filter Filterer) (states []Stater, err error)
}
type Updater interface {
	UpdateAny(state Stater) (updated Stater, err error)
}
type Deleter interface {
	DeleteAny(id string) (deleted Stater, err error)
}

type Storage interface {
	Getter
	Putter
	Lister
	Updater
	Deleter
}

type Fields []string
type Row []interface{}

type Rower interface {
	Row() (Fields, Row)
}

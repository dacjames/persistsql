package devices

import (
	"database/sql"
	"errors"

	"github.com/dacjames/persistsql/internal/resource"
	"github.com/dacjames/persistsql/internal/storage"
)

type psql struct {
	db         *sql.DB
	collection storage.Collection
}

func NewPsql(db *sql.DB) storage.Storage {
	return &psql{
		db:         db,
		collection: storage.NewCollection(db),
	}
}

func (s *psql) GetAny(id resource.ID) (state storage.Stater, err error) {
	dest := &Device{
		Resource: resource.Resource{ID: id},
	}

	if err := s.collection.Select(dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *psql) PutAny(state storage.Stater) (err error) {
	if st, ok := state.(storage.Revisable); ok {
		return s.collection.Revise(st)
	}
	return errors.New("Type Error")
}

func (s *psql) ListAny(filter storage.Filterer) (states []storage.Stater, err error) {
	return nil, nil
}

func (s *psql) UpdateAny(state storage.Stater) (updated storage.Stater, err error) {
	return nil, nil
}

func (s *psql) DeleteAny(id string) (deleted storage.Stater, err error) {
	return nil, nil
}

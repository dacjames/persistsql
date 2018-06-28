package storage

import (
	"database/sql"
	"errors"

	"github.com/dacjames/persistsql/internal/resource"
	"github.com/dacjames/persistsql/internal/util"
	"github.com/jmoiron/sqlx"
)

type ServiceImpl interface {
	ServiceName() string
	Revise(state Stater, tx *sql.Tx) error
	ScanRows(rows *sqlx.Rows) (Stater, error)
}

type Collection struct {
	ServiceImpl
	DB *sql.DB
}

func (d *Collection) PutAny(state Stater) error {

	err := util.WithTransaction(d.DB, func(tx *sql.Tx) error {
		if _, err := tx.Exec(`
			insert into ledger.resources(resource_id, service_id)
			values ($1, (select service_id from ledger.services where name='`+d.ServiceName()+`'))
			on conflict (resource_id) do nothing
		`, state.ResourceID()); err != nil {
			return err
		}

		if err := state.ResourceTags().Insert(tx, state.ResourceID()); err != nil {
			return err
		}

		if err := d.Revise(state, tx); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (d *Collection) GetAny(id resource.ID) (Stater, error) {
	dbx := sqlx.NewDb(d.DB, "postgres")

	// Unsafe here means that unmatched sql fields will be ignored.
	// This is useful because we ignore aliased duplicate fields
	// while avoid a Impl.QueryFields() or something to supply
	// fields
	result, err := dbx.Unsafe().Queryx(`
		select revision_id as "meta.revision_id",
			   created_at as "meta.created_at",
			   updated_at as "meta.updated_at",
			   d.*
		from latest.`+d.ServiceName()+` d
		where resource_id = $1 limit 1
	`, id)
	if err != nil {
		return nil, err
	}

	var state Stater
	for result.Next() {
		var err error
		state, err = d.ScanRows(result)
		if err != nil {
			return nil, err
		}
	}

	if state != nil {
		return state, nil
	}

	return nil, errors.New("Resource ID %s Not Found")
}

package storage

import (
	"database/sql"
	"database/sql/driver"
	"strings"

	"github.com/dacjames/persistsql/internal/model"
	"github.com/dacjames/persistsql/internal/util"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Valueser interface {
	Values() []driver.NamedValue
}

type ScanRowser interface {
	ScanRows(rows *sqlx.Rows) error
}

type Revisable interface {
	Valueser
	model.Resourcer
}

type Selectable interface {
	ScanRowser
	model.ResourceServicer
	model.ResourceIDer
}

type Collection interface {
	Revise(rev Revisable) error
	Select(dest Selectable) error
}

type collection struct {
	db *sql.DB
}

func NewCollection(db *sql.DB) Collection {
	return &collection{db: db}
}

func (c *collection) Revise(rev Revisable) error {
	if err := util.WithTransaction(c.db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(`
			insert into ledger.resources(resource_id, service_id)
			values ($1, (select service_id from ledger.services where name='`+rev.ResourceService()+`'))
			on conflict (resource_id) do nothing
		`, rev.ResourceID()); err != nil {
			return err
		}

		if err := rev.ResourceTags().Insert(tx, rev.ResourceID()); err != nil {
			return err
		}

		values := rev.Values()
		names := make([]string, len(values))
		vv := make([]interface{}, len(values))
		ph := util.PlaceholderValue(len(values))
		for i, v := range values {
			names[i] = v.Name
			vv[i] = v.Value
		}

		if _, err := tx.Exec(`
			insert into ledger.devices(`+strings.Join(names, ",")+`)
			values `+ph+`
		`, vv...); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *collection) Select(dest Selectable) error {
	dbx := sqlx.NewDb(c.db, "postgres")

	// Unsafe here means that unmatched sql fields will be ignored.
	// This is useful because we ignore aliased duplicate fields
	// while avoid a Impl.QueryFields() or something to supply
	// fields
	result, err := dbx.Unsafe().Queryx(`
		select revision_id as "meta.revision_id",
			   created_at as "meta.created_at",
			   updated_at as "meta.updated_at",
			   d.*
		from latest.`+dest.ResourceService()+` d
		where resource_id = $1 limit 1
	`, dest.ResourceID())
	if err != nil {
		return err
	}

	found := false
	for result.Next() {
		found = true
		err := dest.ScanRows(result)
		if err != nil {
			return err
		}
	}

	if !found {
		return errors.Errorf("Resource ID %s Not Found", dest.ResourceID())
	}

	return nil
}

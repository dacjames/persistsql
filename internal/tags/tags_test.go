package tags_test

import (
	"database/sql"
	"testing"

	"github.com/dacjames/persistsql/internal/resource"
	"github.com/dacjames/persistsql/internal/tags"
	"github.com/dacjames/persistsql/internal/test_util"
	"github.com/dacjames/persistsql/internal/util"

	"github.com/stretchr/testify/require"
)

func TestInsertTags(t *testing.T) {
	test_util.WithMigratedDB(t, func(db *sql.DB) {
		require.Equal(t, true, true)

		rid := resource.NewID(nil)
		require.NotEmpty(t, rid)

		if _, err := db.Query(`
			INSERT INTO ledger.resources(resource_id, service_id)
			VALUES ($1, (SELECT service_id FROM ledger.services WHERE name='devices'))
		`, rid); true {
			require.Nil(t, err)
		}

		util.WithTransaction(db, func(tx *sql.Tx) error {
			err := tags.Must("a=b", "a=c").Insert(tx, rid)
			require.Nil(t, err)
			return nil
		})

		util.WithTransaction(db, func(tx *sql.Tx) error {
			err := tags.Tagset{}.Insert(tx, rid)
			require.Nil(t, err)
			return nil
		})

		util.WithTransaction(db, func(tx *sql.Tx) error {
			err := tags.Tagset{}.Insert(tx, rid)
			require.Nil(t, err)
			return nil
		})

	})
}

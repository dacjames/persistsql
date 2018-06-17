package resource_test

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/stretchr/testify/require"

	"github.com/dacjames/persistsql/internal/resource"
	"github.com/dacjames/persistsql/internal/test_util"
)

func TestResourceID(t *testing.T) {
	require.Equal(t, true, true)

	r := resource.NewID(nil)

	require.NotEmpty(t, r.String())

	_, err := ulid.Parse(r.String())
	require.Nil(t, err)

	_, err = uuid.Parse(r.UUID())
	require.Nil(t, err)
}

// TestResourceIDScanValue roundtrips a resource.ID through sql
// to test its Scanner and Valuer implementations
func TestResourceIDScanValue(t *testing.T) {
	test_util.WithMigratedDB(t, func(db *sql.DB) {
		in := resource.NewID(nil)

		var out resource.ID
		if err := db.QueryRow(`
			INSERT INTO ledger.resources(resource_id, service_id)
			VALUES ($1, (SELECT service_id FROM ledger.services WHERE name='devices'))
			RETURNING resource_id
		`, in).Scan(&out); true {
			require.NoError(t, err)
			require.Equal(t, in, out)
		}
	})
}

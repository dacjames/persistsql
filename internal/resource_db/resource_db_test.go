package resource_db_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dacjames/persistsql/internal/resource"
	"github.com/dacjames/persistsql/internal/test_util"
)

var db *sql.DB

func TestMain(m *testing.M) {
	db, _ = test_util.SetupTestDB()
	mCtx := test_util.Migrate(db)

	code := m.Run()

	mCtx.Teardown()

	os.Exit(code)
}

func TestResourceIDScanValue(t *testing.T) {
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
}

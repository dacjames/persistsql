package persistsql

import (
	_ "github.com/lib/pq"

	"github.com/dacjames/persistsql/pkg"
)

var _ = pkg.NewCollection

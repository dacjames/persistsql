package devices

import (
	"database/sql"
	"database/sql/driver"

	"github.com/dacjames/persistsql/internal/resource"
	"github.com/dacjames/persistsql/internal/tags"
	"github.com/jmoiron/sqlx"
)

type Device struct {
	resource.Resource

	Name string `db:"name"`
	Tags tags.Tagset
}

func (d *Device) State()                    {}
func (d *Device) ResourceID() resource.ID   { return d.ID }
func (d *Device) ResourceTags() tags.Tagset { return d.Tags }
func (d *Device) ResourceService() string   { return "devices" }

type DeviceStorage struct {
	DB *sql.DB
}

func (d *Device) Values() []driver.NamedValue {
	return []driver.NamedValue{
		driver.NamedValue{
			Name:    "resource_id",
			Ordinal: 1,
			Value:   d.ID,
		},
		driver.NamedValue{
			Name:    "name",
			Ordinal: 2,
			Value:   d.Name,
		},
	}
}

func (d *Device) ScanRows(rows *sqlx.Rows) error {
	if err := rows.StructScan(d); err != nil {
		return err
	}
	return nil
}

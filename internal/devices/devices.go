package devices

import (
	"database/sql"
	"errors"

	"github.com/dacjames/persistsql/internal/resource"
	"github.com/dacjames/persistsql/internal/storage"
	"github.com/dacjames/persistsql/internal/tags"
	"github.com/jmoiron/sqlx"
)

type Device struct {
	resource.Resource

	Name string `db:"name"`
	Tags tags.Tagset
}

var _ storage.ServiceImpl = &DeviceStorage{}

func (d *Device) State()                    {}
func (d *Device) ResourceID() resource.ID   { return d.ID }
func (d *Device) ResourceTags() tags.Tagset { return d.Tags }

type DeviceStorage struct {
	DB *sql.DB
}

func (d *DeviceStorage) ServiceName() string {
	return "devices"
}

func (d *DeviceStorage) Revise(state storage.Stater, tx *sql.Tx) error {
	device, ok := state.(*Device)
	if !ok {
		return errors.New("Not a device")
	}

	if _, err := tx.Exec(`
		insert into ledger.devices(resource_id, name)
		values ($1, $2)
	`, device.ID, device.Name); err != nil {
		return err
	}

	return nil
}

func (d *DeviceStorage) ScanRows(rows *sqlx.Rows) (storage.Stater, error) {
	device := &Device{}

	err := rows.StructScan(device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

package util

import "database/sql"

func WithTransaction(db *sql.DB, f func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Setup a "catch" to rollback the transaction
	// if the callback panics
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Call the callback
	err = f(tx)

	// If the callback failed, rollback the transaction
	if err != nil {
		tx.Rollback()
		return err
	}

	// Otherwise, things worked, so commit the transaction
	err = tx.Commit()
	return err
}

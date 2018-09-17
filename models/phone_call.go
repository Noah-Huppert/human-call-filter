package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// PhoneCall holds information about a call received by the server
type PhoneCall struct {
	// ID is a unique identifier
	ID int64

	// PhoneNumberID is the ID of the associated phone number which made
	// the call
	PhoneNumberID int64 `db:"phone_number_id"`

	// DateReceived is the date and time the phone call was received
	DateReceived time.Time `db:"date_received"`
}

// Insert adds a phone call row to the phone calls database. The ID field is
// updated with the inserted row's ID.
func (c *PhoneCall) Insert(db *sqlx.DB) error {
	// Start transaction
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("error starting transaction: %s", err.Error())
	}

	// Insert
	err = tx.QueryRowx("INSERT INTO phone_calls (phone_number_id, "+
		"date_received) VALUES ($1, $2) RETURNING id", c.PhoneNumberID,
		c.DateReceived).StructScan(c)

	if err != nil {
		return fmt.Errorf("error executing insert statement: %s", err.Error())
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %s", err.Error())
	}

	return nil
}

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

	// TwilioCallID is the unique ID twilio assigned to the call
	TwilioCallID string `db:"twilio_call_id"`

	// DateReceived is the date and time the phone call was received
	DateReceived time.Time `db:"date_received"`
}

// QueryAllPhoneCalls retrieves a list of all the phone calls in the
// database.
func QueryAllPhoneCalls(db *sqlx.DB) ([]PhoneCall, error) {
	calls := []PhoneCall{}

	rows, err := db.Queryx("SELECT * FROM phone_calls")
	if err != nil {
		return calls, fmt.Errorf("error executing query: %s",
			err.Error())
	}

	for rows.Next() {
		call := PhoneCall{}

		err = rows.StructScan(&call)
		if err != nil {
			return []PhoneCall{}, fmt.Errorf("error scanning row into "+
				"struct: %s", err.Error())
		}

		calls = append(calls, call)
	}

	return calls, nil
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
		"twilio_call_id, date_received) VALUES ($1, $2, $3) RETURNING id",
		c.PhoneNumberID, c.TwilioCallID, c.DateReceived).StructScan(c)

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

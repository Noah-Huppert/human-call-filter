package models

import (
	"github.com/jmoiron/sqlx"
)

// PhoneNumber holds information about a phone call
type PhoneNumber struct {
	// ID is a unique identifier
	ID int

	// Number is the phone number
	Number string

	// Name is the phone number owner's name, empty if unknown
	Name string

	// State is the state the phone number is registered in
	State string

	// City is the city the phone number is registered in
	City string

	// ZipCode is the zip code the phone number is registered in
	ZipCode string `db:"zip_code"`
}

// QueryByNumber attempts to find a row in the phone numbers table a matching
// Number value.
//
// The sql.ErrNoRows error is returned if no rows are found.
func (n *PhoneNumber) QueryByNumber(db *sqlx.DB) error {
	return db.Get(n, "SELECT id, number, name, state, city, zip_code FROM "+
		"phone_numbers WHERE number = $1", n.Number)
}

package models

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// PhoneNumber holds information about a phone call
type PhoneNumber struct {
	// ID is a unique identifier
	ID int64

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

// QueryPhoneNumbersCount retrieves the number of phone numbers in the database
func QueryPhoneNumbersCount(db *sqlx.DB) (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(1) FROM phone_numbers").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error executing query: %s", err.Error())
	}

	return count, nil
}

// QueryAllPhoneNumbers retrieves a list of all the phone numbers in the
// database.
func QueryAllPhoneNumbers(db *sqlx.DB) ([]PhoneNumber, error) {
	numbers := []PhoneNumber{}

	rows, err := db.Queryx("SELECT * FROM phone_numbers")
	if err != nil {
		return numbers, fmt.Errorf("error executing query: %s",
			err.Error())
	}

	for rows.Next() {
		number := PhoneNumber{}

		err = rows.StructScan(&number)
		if err != nil {
			return []PhoneNumber{}, fmt.Errorf("error scanning row into "+
				"struct: %s", err.Error())
		}

		numbers = append(numbers, number)
	}

	return numbers, nil
}

// QueryByNumber attempts to find a row in the phone numbers table a matching
// Number value.
//
// The sql.ErrNoRows error is returned if no rows are found.
func (n *PhoneNumber) QueryByNumber(db *sqlx.DB) error {
	return db.Get(n, "SELECT id, number, name, state, city, zip_code FROM "+
		"phone_numbers WHERE number = $1", n.Number)
}

// Insert adds a phone number to the database. The ID field is updated with the
// inserted row's id value
func (n *PhoneNumber) Insert(db *sqlx.DB) error {
	// Start transaction
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("error starting transaction: %s", err.Error())
	}

	// Insert
	err = tx.QueryRowx("INSERT INTO phone_numbers (number, name, state, "+
		"city, zip_code) VALUES ($1, $2, $3, $4, $5) RETURNING id", n.Number,
		n.Name, n.State, n.City, n.ZipCode).StructScan(n)

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

// HasAChallengePass queries the database to determine if a phone number with
// a matching ID has every passed a challenge before.
//
// If it has passed a challenge before true is returned, false otherwise.
func (n PhoneNumber) HasAChallengePass(db *sqlx.DB) (bool, error) {
	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM phone_calls, challenges "+
		"WHERE phone_calls.phone_number_id = $1 AND "+
		"challenges.status = $2 LIMIT 1)", n.ID, ChallengeStatusPassed).
		Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("error executing select statement: %s",
			err.Error())
	}

	return exists, nil
}

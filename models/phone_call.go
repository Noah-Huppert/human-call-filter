package models

import (
	"time"
)

// PhoneCall holds information about a call received by the server
type PhoneCall struct {
	// ID is a unique identifier
	ID int64

	// PhoneNumberID is the ID of the associated phone number which made
	// the call
	PhoneNumberID int64 `db:"phone_number_id"`

	// DateReceived is the date and time the phone call was received
	DateReceived time.Time
}

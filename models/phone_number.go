package models

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
	ZipCode string
}

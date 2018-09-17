package models

import (
	"time"
)

const (
	// StatusAnswering indicate that a caller is currently answering
	// a challenge
	StatusAnswering string = "ANSWERING"

	// StatusFailed indicates that a caller failed a challenge
	StatusFailed string = "FAILED"

	// StatusPassed indicates a caller succeeded in answering a challenge
	StatusPassed string = "PASSED"
)

// Challenge holds information about a question asked to a caller. Used to
// verify if the caller is human or not.
type Challenge struct {
	// ID is a unique identifier
	ID int64

	// PhoneCallID is the phone call which the challenge was asked during
	PhoneCallID int64 `db:"phone_call_id"`

	// DateAsked is when the question was asked
	DateAsked time.Time `db:"date_asked"`

	// OperandA is the first operand in the challenge
	OperandA int `db:"operand_a"`

	// OperandB is the second operand in the challenge
	OperandB int `db:"operand_b"`

	// Solution is the solution to the challenge
	Solution int

	// Status indicates if the caller is currently answering the question,
	// failed to answer the question, or succeeded to answer the question
	Status string
}

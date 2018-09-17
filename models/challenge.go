package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	// ChallengeStatusAnswering indicate that a caller is currently answering
	// a challenge
	ChallengeStatusAnswering string = "ANSWERING"

	// ChallengeStatusFailed indicates that a caller failed a challenge
	ChallengeStatusFailed string = "FAILED"

	// ChallengeStatusPassed indicates a caller succeeded in answering a challenge
	ChallengeStatusPassed string = "PASSED"
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

// Insert adds a challenge row to the challenges database. The ID field is
// updated with the inserted row's ID.
func (c *Challenge) Insert(db *sqlx.DB) error {
	// Start transaction
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("error starting transaction: %s", err.Error())
	}

	// Insert
	err = tx.QueryRowx("INSERT INTO challenges (phone_call_id, date_asked, "+
		"operand_a, operand_b, solution, status) VALUES ($1, $2, $3, $4, $5, "+
		"$6)", c.PhoneCallID, c.DateAsked, c.OperandA, c.OperandB, c.Solution,
		c.Status).StructScan(c)

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

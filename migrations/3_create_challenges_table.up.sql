CREATE TABLE challenges (
	id SERIAL PRIMARY KEY,
	phone_call_id INTEGER REFERENCES phone_calls,
	date_asked TIMESTAMP WITH TIME ZONE NOT NULL,
	operand_a INTEGER NOT NULL,
	operand_b INTEGER NOT NULL,
	solution INTEGER NOT NULL,
	status challenge_status NOT NULL
)

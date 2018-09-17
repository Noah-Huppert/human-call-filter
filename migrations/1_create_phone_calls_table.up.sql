CREATE TABLE phone_calls (
	id SERIAL PRIMARY KEY,
	phone_number_id INTEGER REFERENCES phone_numbers,
	date_received TIMESTAMP WITH TIME ZONE NOT NULL
)

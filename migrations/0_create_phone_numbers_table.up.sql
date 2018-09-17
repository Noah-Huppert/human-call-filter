CREATE TABLE phone_numbers (
	id SERIAL PRIMARY KEY,
	number TEXT NOT NULL,
	name TEXT,
	state TEXT NOT NULL,
	city TEXT NOT NULL,
	zip_code TEXT NOT NULL
)

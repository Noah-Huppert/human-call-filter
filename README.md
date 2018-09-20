Project status: Active development

# Human Call Filter
Validates that a caller is human before forwarding calls to your phone.

# Table Of Contents
- [Overview](#overview)
- [Setup](#setup)
	- [Configuration](#configuration)
	- [Database](#database)
	- [Twilio](#twilio)

# Overview
Human call filter operates a separate phone number that verifies callers are 
human before forwarding the call to your actual phone.  

This technique is able to prevent robot spam calls.  

When a call is made to the human call filter operated number a prompt is played 
asking the caller a simple arithmetic question. If the caller enters the 
correct answer to their call is forwarded to your actual phone number. If they 
answer the question incorrectly the call will be ended.

If the same caller calls again they will be immediately forwarded without 
being tested, since they have already been verified.

# Setup
## Configuration
Human call filter uses environment variables for configuration:

- `DESTINATION_NUMBER`: Your personal phone number which verified humans will 
	be forwarded to
- `HTTP_PORT`: Port server responds to requests on, defaults to `8000`
- `DB_HOST`: Postgres connection host
- `DB_PORT`: Postgres connection port, defaults to `5432`
- `DB_NAME`: Postgres database name
- `DB_USERNAME`: Postgres username
- `DB_PASSWORD`: Postgres password, optional

## Database
Human call filter stores information about calls in a Postgres database.  

Set the database related configuration variables (start with `DB`) and run:

```
make migrate
```

This will run migrations on your database to setup the schema.

## Twilio
Human call filter uses a Twilio number.  

Complete the following steps to setup Twilio:

1. Login to Twilio  
2. Navigate to the "Phone Numbers" dashboard  
3. Buy a new number  
4. Configure the number  
  - Under the "Voice" category set the following options:
    - "A call comes in" setting
	  - First setting to `Webhook` 
	  - Second setting to `http://<domain you deploy to>/call` and `HTTP POST`
	- " Primary handler fails" setting
	  - First setting to `TwiML`
	  - Click the plus button
	    - Enter `CallDestination` for the "Friendly Name"
		- Enter the following text in the main box:
			```
			<?xml version="1.0" encoding="UTF-8"?>
            <Response>
                <Dial>Your destination number here</Dial>
            </Response>
            ```
			Make sure to replace `Your destination number here` with your 
			actual phone number.

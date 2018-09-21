Project status: Active development

# Human Call Filter
Validates that a caller is human before forwarding calls to your phone.

# Table Of Contents
- [Overview](#overview)
- [Features](#features)
	- [Phone Number](#phone-number)
	- [Dashboard](#dashboard)
- [Setup](#setup)
	- [Configuration](#configuration)
	- [Database](#database)
	- [Twilio](#twilio)
- [Deploy](#deploy)
	- [Kubernetes](#kubernetes)
	- [Docker](#docker)

# Overview
Human call filter operates a separate phone number that verifies callers are 
human before forwarding the call to your actual phone.  

This technique is able to prevent robot spam calls.  

A dashboard is also provided to view information about calls which have been 
made to the human call filter operated phone number.

# Features
## Phone Number
When a call is made to the human call filter operated number a prompt is played 
asking the caller a simple arithmetic question. If the caller enters the 
correct answer, their call is forwarded to your actual phone number. If they 
answer the question incorrectly the call will be ended.

If the same caller calls again they will be immediately forwarded without 
being tested, since they have already been verified.

## Dashboard
A dashboard is provided which shows information about calls which have been 
made to the human call filter operated phone number.

### Phone Numbers Page
The dashboard shows which numbers have placed calls:

![Dashboard numbers page](/imgs/dashboard-phone-numbers-screenshot.png)

### Phone Calls Page
The dashboard shows all the calls which have been received:

![Dashboard calls page](/imgs/dashboard-phone-calls-screenshot.png)

### Challenges Page
The dashboard shows the status of all the challenges which have been issued to 
users.  

![Dashboard challenges page](/imgs/dashboard-challenges-screenshot.png)

# Setup
## Configuration
Human call filter uses environment variables for configuration:

- `DESTINATION_NUMBER`: Your personal phone number which verified humans will 
	be forwarded to
- `CALLS_HTTP_PORT`: Port server responds to twilio call requests on,
	defaults to `8000`
- `DASHBOARD_HTTP_PORT`: Port server responds to internal dashboard request on,
	defaults to `8001`
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

# Deploy
## Kubernetes
A Kubernetes Helm chart is located in the `deploy/human-call-filter/` 
directory.  

This chart is currently configured to setup my (Noah Huppert) personal instance 
of human call filter.  

However the chart can be modified to deploy your personal instance. It requires 
that you have 
[External DNS](https://github.com/kubernetes-incubator/external-dns) installed 
on your cluster.  

First edit the `deploy/human-call-filter/values.yaml` file to your liking.  

Then set the `DESTINATION_NUMBER`, `DB_HOST`, and `DB_PASSWORD` environment 
variables to their production values.  

Next run the following in the repository root:

```
git submodule init --update
./deploy/k8s-deploy/k8s-deploy deploy
```

## Docker
A Docker image is published which runs the human call filter server.  

The image is named `noahhuppert/human-call-filter`, the container tag is the 
source code git commit. This will change with every release. Check the Docker 
hub for the latest tag.  

Run the image by setting all the required environment variables and exposing 
the correct port for call HTTP traffic.  

Make sure to not expose the dashboard HTTP port, as this shows private 
information without requiring authentication. Instead The dashboard port should 
be accessed via a private network.

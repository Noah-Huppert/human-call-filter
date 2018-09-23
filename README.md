Project status: Complete | Actively maintaining

# Human Call Filter
Captcha for phone calls.  

# Table Of Contents
- [Overview](#overview)
- [Features](#features)
	- [Phone Number](#phone-number)
	- [Dashboard](#dashboard)
- [Setup](#setup)
	- [Configuration](#configuration)
	- [Audio Clips](#audio-clips)
	- [Database](#database)
	- [Twilio](#twilio)
	- [Voice Over IP](#voice-over-ip)
- [Deploy](#deploy)
	- [Kubernetes](#kubernetes)
	- [Docker](#docker)

# Overview
Captcha for phone calls.  

Human call filter operates a separate phone number that verifies callers are 
human before forwarding the call to a VOIP number.

This technique is able to prevent robot spam calls.  

A dashboard is also provided to view information about calls which have been 
made to the human call filter operated phone number.

# Features
## Phone Number
When a call is made to the human call filter operated number the caller is 
asked a simple arithmetic question.  

If the caller enters the correct answer, their call is forwarded to 
a VOIP number. If they answer the question incorrectly the call will be ended.

This allows you to trust that all calls which reach this VOIP number are from 
real humans and not spam bots.

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

- `DESTINATION`: Your VOIP address which verified humans will be forwarded to
	- This value will be in the format:
		```
		username@sip_sub_domain.sip.us1.twilio.com
		```

		Replace `username` with the username you set while configuring the 
		Twilio Elastic SIP Trunking Credentials List.  

		Replace `sip_sub_domain` with the sub domain you set when configuring 
		the Twilio Elastic SIP Trunking domain name
- `CALLS_HTTP_PORT`: Port server responds to twilio call requests on,
	defaults to `8000`
- `DASHBOARD_HTTP_PORT`: Port server responds to internal dashboard request on,
	defaults to `8001`
- `DB_HOST`: Postgres connection host
- `DB_PORT`: Postgres connection port, defaults to `5432`
- `DB_NAME`: Postgres database name
- `DB_USERNAME`: Postgres username
- `DB_PASSWORD`: Postgres password, optional

## Audio Clips
Human call filter plays audio clips during calls to communicate with
the caller.

These clips are located in the `./audio-clips/` directory.

- `./audio-clips/`
	- `intro.mp3`: Played before a caller is asked a challenge to verify they
		are human
	- `success.mp3`: Played after a caller provides the correct answer to a
		challenge and before they are forwarded to the VOIP phone address
	- `fail.mp3`: Played after a caller provides the wrong answer to a 
		challenge and before the call is ended

It is recommended that you record your own versions of each of these clips. To 
ensure that callers who hear these clips know that they have reached 
your number.

## Database
Human call filter stores information about calls in a Postgres database.  

Set the database related configuration variables (start with `DB`) and run:

```
make migrate
```

This will run migrations on your database to setup the schema.

## Twilio
Human call filter uses Twilio to interact with phone calls.  

Login to Twilio and setup the following components:

- [Programmable Voice](#programmable-voice)
- [Elastic SIP Trunking](#elastic-sip-trunking)

### Programmable Voice
Navigate to the "Phone Numbers" dashboard.  

First buy a number:

1. Navigate to the "Buy a Number" page in the dashboard
2. Buy a new number

Then configure the number:

1. Navigate to the "Manage Numbers" page in the dashboard
2. Click on the number you just bought
3. Under the "Voice" category set the following options:
	- "Configure with" setting
		- Set to `Webhooks, TwiML Bins, Functions, Studio, or Proxy`
	- "A call comes in" setting
		- First setting to `Webhook` 
	    - Second setting to `http://<domain you deploy to>/call` and `HTTP POST`
	- " Primary handler fails" setting
        - First setting to `TwiML`
        - Click the plus button, a box will pop up, enter the following:
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

				This will forward calls to your phone if an error occurs.
		- Make sure the box above the plus button you just clicked is set to 
			`CallDestination`

### Elastic SIP Trunking
Navigate to the "Elastic SIP Trunking" dashboard.  

First create a SIP trunk:

1. Navigate to the "Trunks" page in the dashboard
2. Create a new SIP trunk by clicking the red plus button in the upper left of 
	the page

Next configure a login for the SIP trunk:

1. Navigate to the "Authentication" > "Credentials List" page in the dashboard
2. Create a new credentials list by clicking the red plus button in the upper 
	left of the page
	- You will use the username and password you enter to connect to the VOIP 
		server from your phone

Finally configure a domain for the SIP trunk to serve requests from:

1. Navigate to the "Programmable Voice" dashboard
2. Navigate to the "SIP Domains" > "Domains" page in the dashboard
3. Create a voice SIP domain by clicking the red plus in the upper left of
	the page
	- Under the "Properties" section set the "SIP URI" option to a sub-domain 
		of your choosing
	- Under the "Voice Authentication" page set the "Credentials List" option 
		to the credentials list you created above
	- Under the "SIP Registration" section click the "Enabled" option
	- Under the "SIP Registration" section set the "Credential Lists" option 
		to the credentials list your created above

## Voice Over IP
Human call filter will forward calls which it knows are from humans to a voice 
over IP number.  

This lets you be confident that every call that comes into this number is a 
human and not a spam robot caller.  

In order to receive calls with this VOIP number you must install and configure 
a VOIP application on your phone.  

If you are using an Android phone I recommend you use the 
[Zoiper](https://play.google.com/store/apps/details?id=com.zoiper.android.app) 
VOIP application. Unfortunately I do not have experience with VOIP applications 
on iPhones so I can not make a recommendation.

On your phone in the VOIP application of your choice add a new VOIP account:

1. Set the host to:
    ```
	sip_sub_domain.sip.us1.twilio.com
	```

	Replace `sip_sub_domain` with the sub domain you set when configuring 
	the Twilio Elastic SIP Trunking domain name
2. Set the username to the username you set while configuring the Twilio 
	Elastic SIP Trunking Credentials List
3. Set the password to the password you set while configuring the Twilio
	Elastic SIP Trunking Credentials List

With your new VOIP number configured you can now begin receiving VOIP calls 
on your phone.
		
# Deploy
## Kubernetes
A Kubernetes Helm chart is located in the `deploy/human-call-filter/` 
directory.  

### Deploy Overview
The Kubernetes deployment provided is fairly opinionated.  

It makes a few assumptions:

- You are running your Kubernets cluster on the Google Cloud Platform
- [External DNS](https://github.com/kubernetes-incubator/external-dns) is 
	installed on your cluster
- You are hosting a Postgresql database on the Google Cloud Platform using 
	Cloud SQL

If any of these assumptions are not true for your deployment situation you may 
have to edit the Kubernetes deployment provided.

### Deploy Configuration
Edit the `deploy/human-call-filter/values.yaml` file and change the 
following options:

- `gcp.cloudsql.connectionName`: The connection name for your GCP 
	CloudSQL instance

Then set the `DESTINATION`, `DB_HOST`, and `DB_PASSWORD` environment 
variables to their production values.  

### Deploy Run
To deploy run the following in the repository root:

```
git submodule init --update
./deploy/k8s-deploy/k8s-deploy deploy
```

## Docker
A Docker image is published which runs the human call filter server.  

### Docker Overview
The image is named [`noahhuppert/human-call-filter`](https://hub.docker.com/r/noahhuppert/human-call-filter/) and is published to the Docker hub.  

The container tag is the source code git commit the image was built with. 
This will change with every release. Check the Docker hub for the latest tag.  

### Docker Run
Run the image by setting 
[all the required environment variables](#configuration) and exposing the 
correct port for call HTTP traffic.  

Make sure to not expose the dashboard HTTP port, as this shows private 
information without requiring authentication. Instead the dashboard port should 
be accessed via a private network.

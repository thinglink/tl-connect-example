# Thinglink Connect

Greetings, human. This directory contains a simple, sample application to demonstrate how to set up Thinglink Connect.

Exact documentation is available at [our support site](https://support.thinglink.com/articles/developer-resources/thinglink-connect).

## Requirements

* Golang 1.9+

## Setup

First, go to https://www.thinglink.com/developer/home and create a new application keypair. You can call it whatever you like, but
use "http://localhost:8000/callback" as the OAuth Callback URL.

## Configuration

Note the client id for your application, and set it up as an environment variable "CLIENTID". You may also choose a local user id
and set it up as USERID. If this is not defined, "sample_user" is used.

## Running the application

    CLIENTID=<your client id> go run main.go

or, if you wish to also simulate some other user id

    CLIENTID=<your client id> USERID=<userid> go run main.go



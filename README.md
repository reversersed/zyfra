# SSO Service

## Installation

After cloning the repo with project you need to `go mod tidy` to install all dependencies<br>
Before building or starting up the application, please, read the [guide](#guide)

## Build

To build an .exe file you need to use `go build cmd/sso/main.go`<br>
After you can use ./main.exe file as usual in terminal

## Start up without build

You can start the application in console without building it using `go run cmd/sso/main.go` in terminal<br>
With `go run` you can use all flags as usual (`-flag value`)

<a name="guide"></a>
## Guide

This application using a `.json` file to store valid usernames and passwords<br>
When starting up, use flag `-file` to specify **absolute** path to a file<br>
If `-file` flag is not provided, the application will use default user `admin admin` which can be changed using flags `-username` and `-password`<br>
You can see flags and description using `go run cmd/sso/main.go -help` or `./main.exe -help`<br>
There already is `users.json` file in this repo which you can use or edit

## Tests

Application has 50% unit tests coverage<br>
You can find coverage info in [tests](/tests) folder or use `make test` to run tests by yourself<br>

> total:								(statements)	50.0%

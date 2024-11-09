# SSO Service

## Installation

After cloning the repo with project you need to `go mod tidy` to install all dependencies<br>
Before building or starting up the application, please, read the [guide](#guide) and [swagger](#swagger)

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

Application has 68.7% unit tests coverage<br>
You can find coverage info in [tests](/tests) folder or use `make test` to run tests by yourself<br>

> total: (statements) 68.7%

<a name="swagger"></a>

## Swagger

Application has swagger you can get access to with URL `http://localhost/api/swagger/index.html` (with default settings)<br>
![Swagger](https://github.com/user-attachments/assets/d3f9497b-8a17-4962-971a-75c51d286ade)

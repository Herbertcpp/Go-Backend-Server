# Go-Backend-Server

A little Backend-Server I wrote in Go after using it for about 3 Weeks now  

## Features

Features
- User registration
- User validation (checking the password)
- SQL-Lite database implementation
- Password encryption with "bcrypt"

## Usage 

### Clone the repository onto your machine and add dependencies with:
```
git clone https://github.com/Herbertcpp/Go-Backend-Server
cd Go-Backend-Server
go mod tidy
```

## Communication
After running the server with "go run server.go" you have following URLs avaviable

### Add a user 
- http://localhost:8080/register 
  - Expects: 
  application/json
  {"username" : string, "password", string}
  - Returns:
  application/json
  {"success" : bool, "message" : string}

### Authenticate a User:
  - http://localhost:8080/authenticate
    - Expects 
    application/json
    {"username" : string, "password", string}
    - Returns 
    application/json 
    {"success" bool, "message" : string}


### print all the currently registered users (Mostly for testing purposes)
- http://localhost:8080/print 
  - Expects: nothing
  - Returns: into the writer

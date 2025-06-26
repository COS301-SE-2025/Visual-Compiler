_Must be run in the backend directory_  
  
**To run all tests:** 
``` 
    go test ./...  
    go test -v ./...  
    go test ./... -coverprofile=coverage.out
    go tool cover -html=coverage.out
```
**To run specific tests:**  
unit tests:
```
    go test -v ./core/unit-tests/lexer-unit_test.go  
    go test -v ./core/unit-tests/parser-unit_test.go 
    go test -v ./api/unit-tests/users-routers-unit_test.go  
    go test -v ./api/unit-tests/delete-user-unit_test.go  
    go test -v ./api/unit-tests/login-unit_test.go  
    go test -v ./api/unit-tests/register-unit_test.go  
    go test -v ./api/unit-tests/lexing-routers-unit_test.go
    go test -v ./api/unit-tests/lexing-unit_test.go
    go test -v ./api/unit-tests/parsing-router-unit_test.go
    go test -v ./api/unit-tests/parsing-unit_test.go
```
integration tests:  
_(requires .env file to run)_
```
    go test -v ./integration-tests/api-mongoDB-integration_test.go  
    go test -v ./integration-tests/api-core-lexer-integration_test.go  
    go test -v ./integration-tests/api-core-parser-integration_test.go  
```
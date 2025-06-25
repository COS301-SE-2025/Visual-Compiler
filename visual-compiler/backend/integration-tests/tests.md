_Must be run in the backend directory_  
  
*To run all tests:* 
``` 
    go test ./...  
    go test -v ./...  
```
**To run specific tests:**  
unit tests:
```
    go test -v ./core/unit-tests/core-lexing-unit_test.go  
    go test -v ./api/unit-tests/api-users-unit_test.go  
    go test -v ./api/unit-tests/api-lexing-unit_test.go
```
integration tests:  
_(requires .env file to run)_
```
    go test -v ./integration-tests/api-mongoDB-integration_test.go  
    go test -v ./integration-tests/api-core-integration_test.go  
```
_Must be run in the backend directory_  
  
*To run all tests:*  
    go test ./...  
    go test -v ./...  
  
*To run specific tests:*  
_unit tests:_  
    go test -v ./tests/unit-tests/core-lexing-unit_test.go  
    go test -v ./tests/unit-tests/api-users-unit_test.go
    go test -v ./tests/unit-tests/api-phases-unit_test.go

_integration tests:_
    go test -v ./tests/integration-tests/api-mongoDB-integration_test.go
    go test -v ./tests/integration-tests/api-core-integration_test.go
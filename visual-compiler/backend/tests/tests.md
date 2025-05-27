_Must be run in the backend directory_  
  
*To run all tests:*  
    go test ./...  
    go test -v ./...  
  
*To run specific tests:*  
    go test -v ./tests/unit-tests/lexing-unit_test.go  
    go test -v ./tests/unit-tests/api-users-unit_test.go
    go test -v ./tests/unit-tests/api-phases-unit_test.go
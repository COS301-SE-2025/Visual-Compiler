FROM golang:1.24.3

WORKDIR /visual-compiler-backend

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080

CMD [ "go", "run" , "main.go" ]

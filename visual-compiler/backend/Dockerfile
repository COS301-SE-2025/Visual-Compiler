FROM golang:1.24.3-alpine AS builder
WORKDIR /visual-compiler-backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o backend -ldflags="-s -w" main.go

FROM alpine:latest
WORKDIR /visual-compiler-backend

COPY --from=builder /visual-compiler-backend/backend .

ENV GOGC=50
ENV GOMAXPROCS=1

EXPOSE 8080

CMD [ "./backend" ]

#build
FROM golang:1.23-alpine as builder
WORKDIR /build
RUN pwd
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /main main.go
#final
FROM alpine:3
COPY config.json config.json
COPY .env .env
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]

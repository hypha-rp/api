FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o hypha-api cmd/app/main.go

FROM alpine:latest
COPY --from=builder /app/hypha-api /usr/local/bin/
CMD ["hypha-api"]
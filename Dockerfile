FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o hypha-api cmd/app/main.go

FROM alpine:3

LABEL org.opencontainers.image.source="https://github.com/hypha-rp/api"
LABEL org.opencontainers.image.description="Image for the Hypha backend API"
LABEL org.opencontainers.image.licenses="Apache-2.0"

ENV GIN_MODE=release

COPY --from=builder /app/hypha-api /usr/local/bin/

CMD ["hypha-api"]
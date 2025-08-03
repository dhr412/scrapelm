FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY src/scrapelm.go .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o ./scrapelm

FROM alpine:latest

COPY --from=builder /app/scrapelm ./scrapelm

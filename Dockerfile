FROM golang:alpine AS builder

WORKDIR /app

ADD go.mod go.sum ./
RUN go mod download

COPY . .

ENV GO111MODULE=on CGO_ENABLED=1 GOOS=linux

RUN go build -o /app/jandan-girl .

FROM alpine:latest

RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true
WORKDIR /app

COPY --from=builder /app/jandan-girl .

ENTRYPOINT ["./jandan-girl"]
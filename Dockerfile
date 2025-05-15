FROM golang:1.24.3-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

WORKDIR /app

COPY . ./

RUN go mod download
RUN go build -o /bin/dockwatch .

FROM alpine:3.21.3 AS latest

WORKDIR /app

COPY --from=builder /bin/dockwatch /app/dockwatch

ENTRYPOINT ["/app/dockwatch"]
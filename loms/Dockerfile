FROM golang:1.22.0-alpine AS build
RUN apk --no-cache add gcc g++ make

WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/loms ./cmd/server/main.go


FROM alpine:latest
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin

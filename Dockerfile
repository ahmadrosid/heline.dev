FROM golang:1.21.1-alpine as base
RUN apk add bash make

WORKDIR /go/src/app
COPY go.* .
RUN go mod download

COPY . .
# Replace the make build command with direct go build
RUN go build -o /go/bin/heline

FROM alpine:3.17.2
COPY --from=base /go/bin/heline /heline
ENV PORT=8000

EXPOSE 8000

CMD ["./heline", "server", "start"]
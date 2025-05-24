FROM golang:1.21.1-alpine as base
RUN apk add bash make

WORKDIR /go/src/app
COPY go.* .
RUN go mod download

COPY . .
RUN go build -o /go/bin/heline

FROM alpine:3.17.2
RUN apk add --no-cache bash make curl
COPY --from=base /go/bin/heline /heline

# Copy Makefile for make commands
COPY Makefile /app/Makefile
WORKDIR /app
ENV PORT=8000

EXPOSE 8000

CMD ["./heline", "server", "start"]

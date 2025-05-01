FROM golang:1.21.1-alpine as base
RUN apk add bash make

WORKDIR /go/src/app
COPY go.* .
RUN go mod download

COPY . .
RUN make build TARGET_DIR=/go/bin/app

FROM alpine:3.17.2
COPY --from=base /go/bin/app /app
ENV PORT=8000

EXPOSE 8000

# Command to run the executable
CMD ["./app", "server", "start"]

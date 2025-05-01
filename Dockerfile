FROM node:18 AS ui-builder

WORKDIR /app
COPY ui /app/ui
WORKDIR /app/ui

RUN npm install -g pnpm
RUN pnpm install
RUN pnpm build

FROM golang:1.24 AS go-builder

WORKDIR /app
COPY . /app
COPY --from=ui-builder /app/ui/dist /app/ui/dist

# Build the Go application
RUN go build -o heline

# Final stage
FROM golang:1.24-slim

WORKDIR /app

# Copy the built executable
COPY --from=go-builder /app/heline /app/
COPY --from=go-builder /app/http /app/http

# Create necessary directories
RUN mkdir -p /app/_build

# Expose the port the app runs on
EXPOSE 8000

# Command to run the executable
CMD ["./heline", "server", "start"]

FROM rust:1.70 as builder

WORKDIR /app
COPY . .

# Build the Rust application in release mode
RUN cargo build --release

# Final stage
FROM debian:bullseye-slim

WORKDIR /app

# Install dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    git \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Copy the built executable
COPY --from=builder /app/target/release/hli /app/
COPY --from=builder /app/repo.json /app/

# Create directories for repositories
RUN mkdir -p /app/repos

# Set the working directory
WORKDIR /app

# Command to run the indexer
CMD ["./hli", "repo.json"]

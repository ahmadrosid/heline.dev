<h1 align="center"> Heline.dev </h1>
<p align="center">
    Code search for development productivity.
</p>

<p align="center">
    <img src="https://raw.githubusercontent.com/ahmadrosid/heline.dev/refs/heads/main/demo.png" />
</p>

## About

Heline.dev is a modern code search tool built with a multi-language technology stack:

- **Rust**: Powers the code highlighting and indexing functionality, providing fast and efficient code processing (using [hl](https://github.com/ahmadrosid/hl) syntax highlighter)
- **Go**: Handles the backend API services and core application logic
- **Next.js**: Delivers a responsive and interactive frontend user interface

## Requirements

### Local Development
- rust
- golang
- nodejs >= 18
- java >= 8

### Docker Deployment
- Docker
- Docker Compose

## Local Development

Run script build to run dev mode locally.

```bash
bash build.sh
```

Install required dependencies.

```bash
bash scripts/bootstrap.sh
```

Start solr

```bash
bash scripts/solr.sh start
```

Prepare solr index

```bash
bash scripts/solr.sh prepare
```

Run production mode

```bash
bash scripts/run.sh production
```

Reset production mode - this will delete the ES data and run the indexer.

```bash
bash scripts/run.sh reset
```

## Docker Deployment

Heline.dev can be easily deployed using Docker Compose. The setup includes three services:

1. **Apache Solr** - Search engine running on port 8984
2. **Heline App** - Main Go application running on port 8000
3. **Heline Indexer** - Rust-based indexer service running on port 8080

### Running with Docker Compose

```bash
# Build and start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop all services
docker compose down
```

### Volumes

The Docker setup uses the following persistent volumes:

- `solr_data`: Stores Solr indexes and configuration
- `app_data`: Stores application build data
- `indexer_repos`: Stores repositories for indexing

### Service Details

- **Solr**: Runs on port 8984 with precreated cores for 'heline'
- **Heline App**: Connects to Solr and the indexer service
- **Heline Indexer**: Provides API for code indexing

All services are connected through the 'heline-network' bridge network.

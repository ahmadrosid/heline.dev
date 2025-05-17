<h1 align="center"> Heline.dev </h1>
<p align="center"> Code search for development productivity. </p>

<p align="center">
    <img src="https://raw.githubusercontent.com/ahmadrosid/heline.dev/refs/heads/main/demo.png" /> 
</p>

When I first started coding, finding the right code examples was always a struggle. That's why I built Heline.dev - a practical code search tool that actually makes sense for me.

## What is this?

Heline.dev combines three powerful technologies:

- **Rust**: Handles the heavy lifting for code highlighting and indexing (using my [hl](https://github.com/ahmadrosid/hl) syntax highlighter)
- **Go**: Powers the backend API - fast and reliable
- **Next.js**: Creates a smooth frontend experience

The hl syntax highlighter (https://github.com/ahmadrosid/hl) is something I built specifically for this project when I couldn't find a highlighting solution that worked exactly how I wanted.

## Getting Started

### For local development
You'll need:
- rust
- golang
- nodejs >= 18
- java >= 8

### For Docker deployment
Just have:
- Docker
- Docker Compose

## Local Development

Want to run it in dev mode? Just do:

```bash
bash build.sh
```

Need the dependencies?

```bash
bash scripts/bootstrap.sh
```

Starting Solr is easy:

```bash
bash scripts/solr.sh start
```

Prepare your index:

```bash
bash scripts/solr.sh prepare
```

For production mode:

```bash
bash scripts/run.sh production
```

Need to reset everything? This will clear ES data and reindex:

```bash
bash scripts/run.sh reset
```

## Docker Deployment

I made Docker deployment super simple. You get three services:

1. **Apache Solr**: Search engine on port 8984
2. **Heline App**: Main Go application on port 8000
3. **Heline Indexer**: Rust-based indexer on port 8080

### Running with Docker

```bash
# Start everything
docker compose up -d

# Check the logs
docker compose logs -f

# Stop when you're done
docker compose down
```

### What's stored where?

- `solr_data`: All your Solr indexes
- `app_data`: Application build files
- `indexer_repos`: Your repositories for searching

### How it connects

- **Solr** runs on port 8984 with a 'heline' core
- **Heline App** connects to both Solr and the indexer
- **Heline Indexer** provides the API for code indexing

Everything talks to each other through the 'heline-network' bridge.

<p align="center">
    <img src="/ui/public/favicon.png" />
</p>

<h1 align="center"> Heline </h1>
<p align="center">
    Code search for modern Developers.
</p>

## Requirements
- rust
- golang
- nodejs >= 16
- java >= 8

## Development
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

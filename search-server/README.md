# Heline Search Engine

This project is the search engine server for [heline.dev](https://heline.dev) written in rust. The engine utilize the tantivy for indexing etc.

## Index Design

So basically the goal of the engine is to return this field to the searcher.

```json
{
    "page": "1",
    "page_size": "20",
    "data": [
        {
            "id": "path/to/file/id",
            "repo": "username/repo-name",
            "path": "/some-file/path/to/file/id",
            "host": "(github.com,gitlab.com)",
            "lang": "JSON",
            "branch": "master",
            "contents": [
                "example <mark>function</mark> with html encoding",
                "example <mark>function</mark> with html encoding",
            ]
        },
        ...
    ]
}
```
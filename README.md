# api-search

[![Lines Of Code](https://tokei.rs/b1/github/kevinmichaelchen/api-search?category=code)](https://github.com/kevinmichaelchen/api-search)

A proof-of-concept search service, powered by
[**Meilisearch**](https://www.meilisearch.com/) and
[**meilisearch-go**](https://github.com/meilisearch/meilisearch-go).

## Project structure

| Directory                                    | Description                               |
|----------------------------------------------|-------------------------------------------|
| [`./cmd`](./cmd)                             | CLI for making gRPC requests              |
| [`./idl`](./idl)                             | Protobufs (Interface Definition Language) |
| [`./internal/app`](./internal/app)           | App dependency injection / initialization |
| [`./internal/idl`](./internal/idl)           | Auto-generated protobufs                  |
| [`./internal/service`](./internal/service)   | Service layer / Business logic            |

## Getting started
```bash
docker-compose up -d
go run main.go
```

## Usage

Check out the [full API](./idl/coop/drivers/search/v1beta1/api.proto).

### Generating fake documents

Generate a fake CSV file of drivers with

```bash
go run cmd/search/*.go generate --size 1000
```

You can view its contents with:
```bash
xsv sort --select last_name,first_name fake-drivers.csv | xsv table | bat --file-name fake-drivers.csv
```

### Ingestion

Index the CSV records into Meilisearch with:
```bash
go run cmd/search/*.go ingest drivers --file fake-drivers.csv
```

### Querying
Perform a search query with:

```bash
go run cmd/search/*.go query --query Nichole
```

Response would look like:
```json
{
  "hits": [
    {
      "driver": {
        "id": "c9q7k6vrirfhbdec6e00",
        "firstName": "Nichole",
        "lastName": "Bailey",
        "email": "Nichole.Bailey@gLGtalk.biz",
        "phone": "108-674-1932"
      }
    }
  ]
}
```
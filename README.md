# vroomberg

## Overview

A CLI tool for querying financial statements

## Usage

1. Initialize the DB, with optional query
2. Run a query against DB


```
  -init string
    	Initialize DB with JSON filepath
  -query string
    	Supported query
```

### Docker Hub

```
docker run --name vroomberg -v $(PWD):/go/src/app rreinold/vroomberg:dev -init ./init.json -query "TSLA *"
```

### Docker, Local

1. Initialize once
```
docker build -t rreinold/vroomberg:dev .
docker run --name vroomberg -v $(PWD):/go/src/app rreinold/vroomberg:dev -init ./init.json -query "TSLA *"
```
2. Then run queries against db
```
docker run --name vroomberg -v $(PWD):/go/src/app rreinold/vroomberg:dev -query "TSLA *"
```




### Bare Metal

```
CGO_ENABLED=1 go run main.go -init ./init.json -query "TSLA *"
CGO_ENABLED=1 go run main.go -query "TSLA *"
```




## Roadmap

1. Bulk INSERT for initializing DB
2. Separate build and production Docker images
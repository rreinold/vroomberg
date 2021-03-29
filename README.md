# vroomberg

## Overview

A CLI tool for querying S&P500 financial statements

## Usage

```
  --init string
    	Initialize DB (./statements.db) with JSON filepath
  --query string
    	Supported query
```

### Basic (Recommended)

#### Docker Hub, Preloaded DB

```
docker run rreinold/vroomberg:dev vroomberg --query "TSLA *"
```

### Advanced

### Docker, Initialize DB

1. Initialize once
```
docker build -t rreinold/vroomberg:dev .
docker run --name vroomberg -v $(PWD):/go/src/app rreinold/vroomberg:dev -init ./init.json -query "TSLA *"
```

### Bare Metal

Requires gcc to compile sqlite3 driver

```
go get
CGO_ENABLED=1 go install
vroomberg -init ./init.json -query "TSLA *"
vroomberg -query "TSLA *"
```

## Roadmap

1. Known Issue: Query Type #4 works only if same keys are used every year
2. Performance: Bulk INSERT for initializing DB
3. DevOps: Separate build and production Docker images
4. Feature: Update representation of floats into decimals to match financial domain
5. QA: Add more units tests
6. QA: Add system tests
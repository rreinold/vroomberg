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
docker run rreinold/vroomberg:1.0.0 vroomberg --query "TSLA *"
```

### Advanced

### Docker, Initialize DB

1. Initialize (takes several min)
```
docker build -t rreinold/vroomberg:dev .
docker run --name vroomberg -v $(PWD):/go/src/app rreinold/vroomberg:dev vroomberg --init ./init.json --query "TSLA *"
```
2. Run queries
```
docker run --name vroomberg -v $(PWD):/go/src/app rreinold/vroomberg:dev vroomberg --query "TSLA *"
```

### Bare Metal

Requires gcc to compile sqlite3 driver

1. Initialize (takes several min)

```
go get
CGO_ENABLED=1 go install
vroomberg --init ./init.json
```

2. Run queries

```
vroomberg --query "TSLA *"
```

## Query

Supported query types:


|Type|Desc|Syntax|Example|
|---|---|---|---|
|1|Returns all companies that meet criteria on most recent annual statement|`<KEY> [<,>] <VALUE>`|`NetIncomeLoss < -400000000`|
|2|Returns single value from a specific <COMPANY> on most recent annual statement|`<COMPANY> <KEY>`|`TSLA NetIncomeLoss`
|3|Generate ratios between two reported values from specific companies on most recent annual statement|`<COMPANY1> <KEY1> / <COMPANY2> <KEY2>`| `TSLA NetIncomeLoss / TSLA OperatingLeasePayments`
|4|Returns all values from a specific company on most recent annual statement|`<COMPANY> *`|`TSLA *`


## Roadmap

1. Known Issue: Query Type #4 works only if same keys are used every year
2. Performance: Bulk INSERT for initializing DB
3. DevOps: Separate build and production Docker images
4. Feature: Update representation of floats into decimals to match financial domain
5. QA: Add more units tests
6. QA: Add system tests
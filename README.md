[PostgRESTman](#postgrestman)

- [Motivation](#motivation)
- [Status](#status)
- [Getting started](#getting-started)

# PostgRESTman

## Motivation

Make rest api having db schema

### Status

Project is in pre-alpha status.

## Getting started

```
go get github.com/iamgoroot/postgrestman

postgrestman -out . -assignPort 8080 -conn "postgres://user:pass@localhost:5432/db_name?sslmode=disable"
```


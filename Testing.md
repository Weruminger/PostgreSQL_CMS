# Testing / Coverage / CI

## Local unit tests
```bash
cd portal
make test
make cover
make coverhtml   # opens ../work/coverage/coverage.html
```

## Local integration tests (PostgreSQL via docker-compose)
```bash
cd portal
make itest
make test-db-down   # stop & remove test db
```

## GitHub Actions
- Unit tests run in `jobs.unit`
- Integration tests use a postgres service and apply `db/*.sql` before `go test -tags=integration`

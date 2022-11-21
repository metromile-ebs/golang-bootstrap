# golang-bootstrap

# Gin (API Framework)
TODO
# Zap (Logging)
TODO

# Goose (DB Migrations)
Documentation: https://github.com/pressly/goose
### Install Goose
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```
or 
```
brew install goose
```
### Create Migration file
As of the time of writing this we are using raw sql for our migrations
```
goose --dir=./init/db/migrations create <MIGRATION_NAME> sql
```
### Run Migrations
```
goose --dir=./init/db/migrations postgres "host=<DB_HOST> user=<DB_USER> password=<DB_PASS> dbname=<DB_NAME> sslmode=disable" up
```

# Terraform (AWS Infrastructure)
TODO
# Helm Charts (Kubernetes Infrastructure)
TODO 
# Jenkins (CI/CD)
TODO 
# Docker
TODO

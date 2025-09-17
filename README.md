# System requirements
Postgresql version: 16
<br>
Golang version: 1.23.4

# Setup
## Installing dependencies
### [Install Golang](https://go.dev/dl/)
### [Install PostgreSQL](https://www.postgresql.org/download/)

## Install Dependencies
Ensure you have all the required dependencies installed. You can do this by running:
``` bash
# Install dependencies
go mod download

# Tidy up the module dependencies
go mod tidy
```

## [Set up database guide](https://drive.google.com/open?id=11YERYXUiPxKusKegHhlN2C42ltfAZA0l&usp=drive_fs)

## Set up environment variables:
Copy the .env.dist file to .env and fill in the necessary environment variables.
``` bash
cp .env.dist .env
```
Fill in the necessary environment variables in the `.env` file.

## Run the server:
Use the following command to start the server:
``` bash
go run cmd/server/main.go
```

## Run the worker:
Use the following command to start the worker:
``` bash
go run cmd/worker/main.go
```

## Run the cron:
Use the following command to start the cron:
``` bash
go run cmd/cron/main.go
```

## Run the migration:
Use the following command to start the migration:
``` bash
go run cmd/migration/main.go [up/down]
```

## Build exec files:
Use the following command to build the exec files:
``` bash
go build -o [output_file_name] cmd/[server/worker/cron/migration]/main.go
```
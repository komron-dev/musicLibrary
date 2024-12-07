# Music(Song) Library

This project implements an online song library service with features such as adding songs, retrieving song details, filtering songs with pagination, deleting songs, and updating song data. The enriched song information is stored in a PostgreSQL database. The service also integrates with an external API described by Swagger to fetch additional song details.

## Setup local development

### Install tools

### MacOS 
```bash
brew install golang-migrate
```

### Windows Using [scoop](https://scoop.sh/) 
```bash 
scoop install migrate
``` 
### Linux (*.deb package) 
```bash 
curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
apt-get update
apt-get install -y migrate
```


- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    brew install sqlc
    ```

- [Gomock](https://github.com/golang/mock)

    ``` bash
    go install github.com/golang/mock/mockgen
    ```

### Steps to run and use the project

- Create ```song_library``` database:

    ```bash
    make create_db
    ```
- Create a new db migration:

    ```bash
    make new_migration name=<migration_name>
    ```
- Run db migration:

    ```bash
    make up_migrate
    ```
- Generate Go code with sqlc:

    ```bash
    make sqlc
    ```
    
- Generate swagger docs:

    ```bash
    make swag-gen
    ```
- Install imported dependencies:
  ```bash
    make import
    ```
  In case you don't have go.mod file, initialize it:
  ```bash
    go mod init
    ```
- Run the program:

    ```bash
    make run
    ```


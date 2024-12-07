package main

import (
	"database/sql"
	"fmt"
	"github.com/komron-dev/musicLibrary/api"
	db "github.com/komron-dev/musicLibrary/db/sqlc"
	"github.com/komron-dev/musicLibrary/util"
	_ "github.com/lib/pq"
)

func main() {
	util.InitLogger()
	logger := util.Logger

	config, err := util.LoadConfigFrom(".")
	if err != nil {
		logger.Fatal("Failed to load config:", err)
	}

	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresDatabase,
	)
	conn, err := sql.Open(config.DBDriver, dbSource)
	if err != nil {
		logger.Fatal("Failed to connect to db:", err)
	}

	store := db.NewStore(conn, logger)
	server, err := api.NewServer(config, store, logger)
	if err != nil {
		logger.Fatal("Failed to create server: ", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		logger.Fatal("Failed to connect to server: ", err)
	}
}

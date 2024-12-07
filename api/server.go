package api

import (
	"github.com/gin-gonic/gin"
	"github.com/komron-dev/musicLibrary/api/docs"
	db "github.com/komron-dev/musicLibrary/db/sqlc"
	"github.com/komron-dev/musicLibrary/util"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
	config util.Config
	logger *logrus.Logger
}

// NewServer @title Song Library
// @version 1.0
// @description Song Library API in Go using Gin and Swagger
// @host localhost:8080
func NewServer(config util.Config, store *db.Store, logger *logrus.Logger) (*Server, error) {
	server := &Server{
		store:  store,
		config: config,
		logger: logger,
	}

	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/"

	router.POST("/songs", server.addSong)
	router.GET("/songs/info", server.getSong)
	router.GET("/songs", server.listSongs)
	router.GET("/songs/:id", server.getSongLyrics)
	router.PUT("/songs", server.updateSong)
	router.DELETE("/songs/:id", server.deleteSong)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.router = router

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

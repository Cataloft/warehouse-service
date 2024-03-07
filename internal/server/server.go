package server

import (
	"lamoda_task/internal/handlers/goods"
	"log/slog"

	"lamoda_task/internal/config"
	"lamoda_task/internal/handlers/warehouses"
	"lamoda_task/internal/middlewares"
	"lamoda_task/internal/storage/postgres"

	"github.com/gin-gonic/gin"
	requestid "github.com/sumit-tembe/gin-requestid"
)

type Server struct {
	router  *gin.Engine
	storage *postgres.Storage
	config  *config.Server
	logger  *slog.Logger
}

func New(db *postgres.Storage, cfg *config.Server, logger *slog.Logger) *Server {
	router := gin.New()

	router.Use(middlewares.LogMiddleware(logger))
	router.Use(gin.Recovery())
	router.Use(requestid.RequestID(nil))

	return &Server{
		router:  router,
		storage: db,
		config:  cfg,
		logger:  logger,
	}
}

func (s *Server) initHandlers() {
	api := s.router.Group("/api")
	api.PATCH("/goods", goods.Update(s.storage, s.logger))
	api.GET("/warehouses/:id", warehouses.GetOne(s.storage, s.logger))
}

func (s *Server) Start() error {
	s.initHandlers()

	return s.router.Run(s.config.Address)
}

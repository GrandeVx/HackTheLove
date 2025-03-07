package routes

import (
	"backend/config"
	"backend/internal/database"
	"backend/internal/logger"
	"backend/internal/server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var log = logger.GetLogger()

func InitRoutes(db database.Database) *gin.Engine {
	r := gin.Default()
	h := NewHandler(db)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}), middleware.RateLimiterMiddleware())

	registerPublicRoutes(r, h)

	protected := r.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	registerProtectedRoutes(protected, h)

	return r
}

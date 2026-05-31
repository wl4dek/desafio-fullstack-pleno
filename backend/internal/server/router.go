package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"backend/internal/auth"
	"backend/internal/children"
	"backend/internal/summary"
)

func SetupRouter(pool *pgxpool.Pool, jwtSecret string, allowOrigins []string) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(requestLogger())

	r.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	childRepo := children.NewChildRepository(pool)
	childService := children.NewChildService(childRepo)
	childHandler := children.NewChildHandler(childService)

	authService := auth.NewAuthService(jwtSecret, time.Hour)
	authHandler := auth.NewAuthHandler(authService)

	summaryService := summary.NewSummaryService(childRepo)
	summaryHandler := summary.NewSummaryHandler(summaryService)

	r.POST("/auth/token", authHandler.Token)

	authGroup := r.Group("/")
	authGroup.Use(auth.AuthMiddleware(jwtSecret))
	{
		authGroup.GET("/children", childHandler.List)
		authGroup.GET("/children/neighborhood", childHandler.ListNeighborhood)
		authGroup.GET("/children/:id", childHandler.GetByID)

		authGroup.PATCH("/children/:id/review", childHandler.MarkReviewed)
		authGroup.GET("/summary", summaryHandler.GetSummary)
	}

	return r
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Msg("request")
	}
}

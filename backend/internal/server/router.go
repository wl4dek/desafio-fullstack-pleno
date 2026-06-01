package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"backend/internal/auth"
	"backend/internal/children"
	"backend/internal/statistics"
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
	statisticRepo := statistics.NewStatisticsRepository(pool)
	childService := children.NewChildService(childRepo)
	childHandler := children.NewChildHandler(childService)

	authService := auth.NewAuthService(jwtSecret, time.Hour)
	authHandler := auth.NewAuthHandler(authService)

	statisticsService := statistics.NewStatisticsService(statisticRepo)
	statisticsHandler := statistics.NewStatisticsHandler(statisticsService)

	r.POST("/auth/token", authHandler.Token)
	r.GET("/auth/session", authHandler.Session)
	r.DELETE("/auth/session", authHandler.Logout)

	authGroup := r.Group("/api/v1")
	authGroup.Use(auth.AuthMiddleware(jwtSecret))
	{
		authGroup.GET("/children", childHandler.List)
		authGroup.GET("/children/neighborhood", childHandler.ListNeighborhood)
		authGroup.GET("/children/:id", childHandler.GetByID)

		authGroup.PATCH("/children/:id/review", childHandler.MarkReviewed)
		authGroup.GET("/summary", statisticsHandler.GetSummary)
		authGroup.GET("/statistics", statisticsHandler.GetStatistics)
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

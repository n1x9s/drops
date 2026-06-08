package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n1x9s/second-brain/backend/internal/api/handlers"
	"github.com/n1x9s/second-brain/backend/internal/infrastructure/security"
	"github.com/n1x9s/second-brain/backend/internal/platform/http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Handlers struct {
	Health       handlers.HealthHandler
	Auth         handlers.AuthHandler
	Memories     handlers.MemoryHandler
	Tasks        handlers.TaskHandler
	Reminders    handlers.ReminderHandler
	Search       handlers.SearchHandler
	Settings     handlers.SettingsHandler
	Integrations handlers.IntegrationHandler
}

func NewRouter(log *zap.Logger, tokens security.TokenManager, handlers Handlers) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.RequestID(), middleware.Logging(log), middleware.Recovery(log), cors())

	router.GET("/health/live", handlers.Health.Live)
	router.GET("/health/ready", handlers.Health.Ready)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/yaml; charset=utf-8", OpenAPIYAML)
	})
	router.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(SwaggerHTML))
	})

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		auth.POST("/register", handlers.Auth.Register)
		auth.POST("/login", handlers.Auth.Login)
		auth.POST("/refresh", handlers.Auth.Refresh)
		auth.POST("/logout", handlers.Auth.Logout)

		protected := api.Group("")
		protected.Use(middleware.Auth(tokens))
		{
			protected.POST("/memories", handlers.Memories.Create)
			protected.GET("/memories", handlers.Memories.List)
			protected.GET("/memories/:id", handlers.Memories.Get)
			protected.PATCH("/memories/:id", handlers.Memories.Update)
			protected.DELETE("/memories/:id", handlers.Memories.Delete)
			protected.GET("/memories/:id/similar", handlers.Memories.Similar)

			protected.POST("/tasks", handlers.Tasks.Create)
			protected.GET("/tasks", handlers.Tasks.List)
			protected.GET("/tasks/:id", handlers.Tasks.Get)
			protected.PATCH("/tasks/:id", handlers.Tasks.Update)
			protected.DELETE("/tasks/:id", handlers.Tasks.Delete)
			protected.POST("/tasks/:id/complete", handlers.Tasks.Complete)

			protected.POST("/reminders", handlers.Reminders.Create)
			protected.GET("/reminders", handlers.Reminders.List)

			protected.GET("/search", handlers.Search.Search)

			protected.GET("/settings", handlers.Settings.Get)
			protected.PUT("/settings", handlers.Settings.Update)
			protected.PUT("/settings/telegram", handlers.Settings.UpsertTelegram)
			protected.PUT("/settings/linear", handlers.Settings.UpsertLinear)

			protected.POST("/telegram/test", handlers.Integrations.TelegramTest)
			protected.POST("/linear/issues", handlers.Integrations.CreateLinearIssue)
			protected.GET("/linear/issues", handlers.Integrations.ListLinearIssues)
			protected.PATCH("/linear/issues/:id", handlers.Integrations.UpdateLinearIssue)
		}
	}

	return router
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

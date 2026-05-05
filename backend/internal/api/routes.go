package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lcs/devhelper/internal/repository"
	"github.com/lcs/devhelper/internal/service"
)

func SetupRoutes(
	r *gin.Engine,
	corsOrigins string,
	authSvc *service.AuthService,
	jsonSvc *service.JsonService,
	userRepo *repository.UserRepo,
	historyRepo *repository.HistoryRepo,
	schemaRepo *repository.SchemaRepo,
) {
	r.Use(CORSMiddleware(corsOrigins))

	authH := NewAuthHandler(authSvc)
	jsonH := NewJsonHandler(jsonSvc, historyRepo, schemaRepo)
	adminH := NewAdminHandler(userRepo)

	v1 := r.Group("/api/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.Refresh)
		auth.GET("/me", AuthMiddleware(), authH.Me)
	}

	// Protected routes
	protected := v1.Group("", AuthMiddleware())
	{
		// User profile
		protected.PUT("/user/profile", UpdateProfile(userRepo))
		protected.PUT("/user/password", ChangePassword(userRepo))
		protected.DELETE("/user/account", DeleteAccount(userRepo))

		// JSON operations
		j := protected.Group("/json")
		{
			j.POST("/validate", jsonH.Validate)
			j.POST("/format", jsonH.Format)
			j.POST("/minify", jsonH.Minify)
			j.POST("/convert", jsonH.Convert)
			j.POST("/parse", jsonH.Parse)
			j.POST("/schema/generate", jsonH.GenerateSchema)
			j.POST("/schema/validate", jsonH.ValidateSchema)
			j.POST("/diff", jsonH.Diff)
			j.POST("/query", jsonH.Query)
		}

		// History
		history := protected.Group("/history")
		{
			history.GET("", jsonH.GetHistory)
			history.POST("", jsonH.SaveHistory)
			history.DELETE("/:id", jsonH.DeleteHistory)
		}

		// Schemas
		schemas := protected.Group("/schemas")
		{
			schemas.GET("", jsonH.ListSchemas)
			schemas.POST("", jsonH.SaveSchema)
			schemas.GET("/:id", jsonH.GetSchema)
			schemas.PUT("/:id", jsonH.UpdateSchema)
			schemas.DELETE("/:id", jsonH.DeleteSchema)
		}

		// Admin routes
		admin := protected.Group("/admin", AdminMiddleware())
		{
			admin.GET("/users", adminH.ListUsers)
			admin.GET("/users/:id", adminH.GetUser)
			admin.PUT("/users/:id", adminH.UpdateUser)
			admin.DELETE("/users/:id", adminH.DeleteUser)
			admin.POST("/users/:id/reset-password", adminH.ResetPassword)
		}
	}
}

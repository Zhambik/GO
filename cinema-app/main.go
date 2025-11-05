package main

import (
	"cinema-app/database"
	"cinema-app/handlers"
	"cinema-app/middleware"

	"github.com/gin-gonic/gin"
)

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	database.InitDB()

	r := gin.Default()

	// Добавляем CORS middleware
	r.Use(CORSMiddleware())

	// Обслуживание статических файлов фронтенда
	r.Static("/static", "./frontend")
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/index.html", "./frontend/index.html")
	r.StaticFile("/login.html", "./frontend/login.html")
	r.StaticFile("/register.html", "./frontend/register.html")
	r.StaticFile("/movies.html", "./frontend/movies.html")
	r.StaticFile("/style.css", "./frontend/style.css")
	r.StaticFile("/app.js", "./frontend/app.js")

	// API routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Public movies routes (без авторизации)
	r.GET("/movies/public", handlers.GetMoviesPublic)

	// Protected movies routes (требуют авторизации)
	auth := r.Group("/movies")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("", handlers.GetMovies)
		auth.POST("", handlers.CreateMovie)
		auth.PUT("/:id", handlers.UpdateMovie)
		auth.DELETE("/:id", handlers.DeleteMovie)
	}

	r.Run(":8080")
}

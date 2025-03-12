package main

import (
	"fmt"
	"os"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No se pudo cargar el archivo .env, usando valores predeterminados")
	}

	// Inicializar configuración y base de datos
	bootstrap.LoadConfig()
	bootstrap.InitDB()

	// Inicializar router de Gin
	r := gin.Default()

	// Configurar CORS para permitir solicitudes del frontend
	r.Use(cors.Default())

	// Middleware global (logs estructurados y manejo de errores)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Ruta para verificar el estado del backend
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Rendalla backend funcionando correctamente 🚀",
		})
	})

	// Rutas públicas (sin autenticación)
	public := r.Group("/")
	{
		public.GET("/songs", handlers.GetAllSongsHandler)
		public.GET("/songs/:id", handlers.GetSongByIDHandler)
		public.GET("/songs/:id/documents", handlers.GetAllDocumentsBySongIDHandler)
		public.GET("/documents/:id", handlers.GetDocumentByIDHandler)
		public.GET("/songs/search", handlers.SearchSongsByTitleHandler)
		public.GET("/documents/search", handlers.SearchDocumentsByTitleHandler)
		public.GET("/documents/filter", handlers.FilterDocumentsByInstrumentHandler)
		public.POST("/auth/login", handlers.LoginHandler)
	}

	// Rutas protegidas (con autenticación JWT)
	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		// Canciones
		auth.POST("/songs", handlers.CreateSongHandler)
		auth.PUT("/songs/:id", handlers.UpdateSongHandler)
		auth.DELETE("/songs/:id", handlers.DeleteSongWithDocumentsHandler)

		// Documentos
		auth.POST("/songs/:id/documents", handlers.CreateDocumentHandler)
		auth.PUT("/documents/:id", handlers.UpdateDocumentHandler)
		auth.DELETE("/documents/:id", handlers.DeleteDocumentHandler)

		// Autenticación
		auth.GET("/auth/me", handlers.MeHandler)
	}

	// Obtener y validar el puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Infof("🚀 Rendalla backend corriendo en el puerto %s", port)
	r.Run(fmt.Sprintf(":%s", port))
}

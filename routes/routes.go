package routes

import (
	"time"

	"github.com/ekchills/go-resume-screener/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	Server *gin.Engine
}

func (s *Routes) RegisterRoutes() {
	uploadController := &controllers.UploadController{}
	ResumesController := &controllers.ResumeController{}
	s.Server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	groupedAuth := s.Server.Group("/api/v1/auth")
	groupedMain := s.Server.Group("/api/v1").Use(controllers.AuthMiddleware())
	groupedAuth.POST("/register", controllers.RegisterController)
	groupedAuth.POST("/login", controllers.LoginController)
	groupedAuth.POST("/refresh", controllers.RefreshToken)
	groupedMain.POST("/upload", uploadController.UploadResume)
	groupedMain.GET("/resumes", ResumesController.GetAllResumes)
	groupedMain.POST("/resumes/shortlist", ResumesController.ShortListResume)
	groupedMain.GET("/resumes/:id", ResumesController.GetResumeDetail)
}


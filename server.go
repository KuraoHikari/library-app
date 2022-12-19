package main

import (
	"github.com/KuraoHikari/library-app/config"
	"github.com/KuraoHikari/library-app/controller"
	"github.com/KuraoHikari/library-app/repository"
	"github.com/KuraoHikari/library-app/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.SetupDatabaseConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	jwtService service.JWTService = service.NewJWTService()
	authService service.AuthService = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
package router

import (
	"portfolio/internal/delivery/rest/controller"
	"portfolio/internal/domain"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterOption struct {
	Service *domain.UserService
}

// @Description Created by Otajonov Quvonchbek
// securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(option RouterOption) *gin.Engine {
	router := gin.New()

	userController := &controller.UserController{Service: *option.Service}

	api := router.Group("/v1")

	// Send
	api.POST("/send-code", userController.SendCode)

	// Sign-Up
	api.POST("/sign-up-phone", userController.SignUpUserByPhone)
	api.POST("/sign-up-email", userController.SignUpUserByEmail)

	// Sign-In
	api.POST("/sign-in-phone", userController.SignInPhone)
	api.POST("/sign-in-email", userController.SignInEmail)

	// Reset-Password
	api.POST("/check-user", userController.CheckUser)
	api.POST("/check-code", userController.CheckCodePhone)
	api.POST("/update-password-phone", userController.UpdatePasswordByPhone)
	api.POST("/update-password-email", userController.UpdatePasswordByEmail)



	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}

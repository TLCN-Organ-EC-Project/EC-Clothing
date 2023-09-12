package api

import (
	"fmt"

	db "github.com/XuanHieuHo/EC_Clothing/db/sqlc"
	"github.com/XuanHieuHo/EC_Clothing/token"
	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config     util.Config
	store      db.Stores
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Stores) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")

	// login and register

	api.POST("/register", server.createUser)
	api.POST("/login", server.loginUser)
	api.POST("/forgotpassword", server.sendResetPasswordToken)
	api.POST("/resetpassword", server.resetPassword)

	// -----------------------------------user--------------------------------
	authUserRoutes := api.Group("/").Use(authMiddleware(server.tokenMaker))
	authUserRoutes.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//user
	authUserRoutes.GET("/users/:username", server.getUserByUsername)
	authUserRoutes.PUT("/users/:username", server.updateUser)
	authUserRoutes.POST("/users/:username/check", server.checkPassword)
	authUserRoutes.POST("/users/:username/change", server.changePassword)
	//feedback
	authUserRoutes.POST("/users/:username/feedbacks/:product_commented", server.createFeedback)
	authUserRoutes.GET("/products/:id/feedbacks", server.listFeedbackByID)
	authUserRoutes.PUT("/users/:username/feedbacks/:product_commented/:id", server.updateFeedback)
	authUserRoutes.DELETE("/users/:username/feedbacks/:product_commented/:id", server.deleteFeedback)

	// -----------------------------------admin--------------------------------
	authAdminRoutes := api.Group("/admin").Use(authAdminMiddleware(server.tokenMaker, server.store))
	authAdminRoutes.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// user
	authAdminRoutes.GET("/users", server.listUser)
	authAdminRoutes.PUT("/users/:username", server.adminUpdateUser)
	authAdminRoutes.GET("/users/:username", server.adminGetUserByUsername)
	authAdminRoutes.DELETE("/users/:username", server.deleteUser)
	// feedback
	authAdminRoutes.GET("/products/:id/feedbacks", server.listFeedbackByID)
	authAdminRoutes.DELETE("/users/:username/feedbacks/:product_commented/:id", server.adminDeleteFeedback)

	server.router = router
}

// Start runs thes HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

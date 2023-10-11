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
		c.Header("Access-Control-Expose-Headers", "X-Session-ID, X-Access-Token, X-Access-Token-Expired-At, X-Refresh-Token, X-Refresh-Token-Expired-At")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

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
	api.POST("/tokens/renew_access", server.renewAccessToken)
	//category
	api.GET("/categories", server.listCategory)
	api.GET("/categories/:id", server.getCategory)
	// province
	api.GET("/provinces", server.listProvinces)
	api.GET("/provinces/:id", server.getProvinceByID)
	//product
	api.GET("/products/:id", server.getProductByID)
	api.GET("/products", server.listProduct)
	api.GET("/products/find", server.findProduct)
	// product in category
	api.GET("/categories/:id/products", server.getProductsInCategory)

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
	//promotion
	authUserRoutes.GET("/promotions/:title", server.getPromotionByTitle)
	authUserRoutes.GET("/promotions/", server.listPromotion)
	// order
	authUserRoutes.POST("/users/:username/orders", server.createOrder)
	authUserRoutes.GET("/users/:username/orders", server.listOrderByUser)
	authUserRoutes.GET("/users/:username/orders/:booking_id", server.getOrder)
	authUserRoutes.PUT("/users/:username/orders/:booking_id", server.updateOrder)
	authUserRoutes.PUT("/users/:username/orders/:booking_id/cancel", server.cancelOrder)
	// order details
	authUserRoutes.GET("/users/:username/orders/:booking_id/detail", server.getDetailOrderByBookingID)
	// carts
	authUserRoutes.POST("/users/:username/carts", server.createCart)
	authUserRoutes.GET("/users/:username/carts", server.listCartOfUser)
	authUserRoutes.PUT("/users/:username/carts/:cart_id", server.updateCart)
	authUserRoutes.DELETE("/users/:username/carts/:cart_id", server.deleteCart)

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
	// promotion
	authAdminRoutes.POST("/promotions", server.createPromotion)
	authAdminRoutes.GET("/promotions/:title", server.getPromotionByTitle)
	authAdminRoutes.GET("/promotions/", server.listPromotion)
	authAdminRoutes.PUT("/promotions/:title", server.updatePromotion)
	authAdminRoutes.DELETE("/promotions/:id", server.deletePromotion)
	// category
	authAdminRoutes.POST("/categories", server.adminCreateCategory)
	authAdminRoutes.PUT("/categories/:id", server.adminUpdateCategory)
	authAdminRoutes.DELETE("/categories/:id", server.adminDeleteCategory)
	// Products
	authAdminRoutes.POST("/products", server.adminCreateProduct)
	authAdminRoutes.PUT("/products/:id", server.adminUpdateProduct)
	// Images of Product
	authAdminRoutes.POST("/products/:id", server.adminAddImageProduct)

	// Products In Category
	authAdminRoutes.POST("/categories/:id/products", server.adminCreateProductInCategory)
	authAdminRoutes.DELETE("/categories/:id/products/:product_id", server.adminDeleteProductInCategory)

	// Add Product to Store
	authAdminRoutes.POST("/products/:id/store", server.adminAddProductToStore)
	authAdminRoutes.PUT("/products/:id/store", server.adminUpdateProductToStore)

	// Order
	authAdminRoutes.GET("/orders", server.adminListOrder)
	authAdminRoutes.GET("/orders/:booking_id", server.adminGetOrderByBookingID)
	authAdminRoutes.GET("/orders/users/:username", server.adminListOrderByUser)
	// Order Detail
	authAdminRoutes.GET("/orders/:booking_id/detail", server.adminGetDetailOrderByBookingID)
	// Province
	authAdminRoutes.POST("/provinces", server.createProvinces)
	server.router = router
}

// Start runs thes HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

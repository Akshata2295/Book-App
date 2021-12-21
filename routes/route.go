package routes

import (
	"Book-App/controllers"
	//_ "Book-App/docs"
	middlewares "Book-App/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

)

func PublicEndpoints(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// Generate public endpoints - [ signup] - api/v1/signup

	r.POST("/register", controllers.Register)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/login", authMiddleware.LoginHandler)
	// r.POST("/loginm", authMiddleware.LoginHandler)
	r.POST("/logout", authMiddleware.LogoutHandler)
}

func AuthenticatedEndpoints(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// Generate Authenticated endpoints - [] - api/v1/auth/
	r.Use(authMiddleware.MiddlewareFunc())

	r.POST("supervisor/create", controllers.CreateSupervisor)
	r.POST("admin/create", controllers.CreateAdmin)

	//books endpoints
	r.POST("books/create", controllers.CreateBook)
	r.GET("books/", controllers.ListAllBook)
	r.GET("books/:id", controllers.GetBook)
	r.PATCH("books/:id", controllers.UpdateBook)
	r.DELETE("books/delete/:id", controllers.DeleteBook)

	//category endpoints
	r.GET("category/", controllers.ListAllCategories)
	r.POST("category/create", controllers.CreateCategory)
	r.PATCH("category/:id", controllers.UpdateCategory)
	r.GET("category/:id", controllers.GetCategory)
	

	//cart endpoints
	r.POST("cart/add", controllers.AddToCart)
	r.GET("cart/view", controllers.ViewCart)
	r.DELETE("cart/delete/:id", controllers.DeleteFromCart)

}

func GetRouter(router chan *gin.Engine) {
	gin.ForceConsoleColor()
   //setting up gin
	r := gin.Default()  //go routines
	r.Use(cors.Default())
	r.Use(middlewares.RequestLogger)
	r.Use(gin.CustomRecovery(middlewares.LogFailedRequests))
	authMiddleware, _ := middlewares.GetAuthMiddleware()

	// Create a BASE_URL - /api/v1
	v1 := r.Group("/api/v1/")
	PublicEndpoints(v1, authMiddleware)
	AuthenticatedEndpoints(v1.Group("auth"), authMiddleware)
	router <- r //channels
}

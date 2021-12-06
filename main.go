
 package main

import (
	//"net/http"

	_ "Book-App/docs"
	"fmt"
	"Book-App/models"
	"Book-App/routes"
	"os"

	"github.com/gin-gonic/gin"
	//jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/joho/godotenv"
)

// @title Book API With JWT
// @version 1.0
// @Description Golang basic API with JWT Authentication and Authorization.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8081
// @BasePath /
func main() {
	// r := gin.Default()
	godotenv.Load()          // Load env variables
	models.ConnectDataBase() // load db
// We want to get the router in async, thus a channel is required to return the router instance.


	var router = make(chan *gin.Engine)
	go routes.GetRouter(router)
	var port string = os.Getenv("SERVER_PORT")
	server_addr := fmt.Sprintf(":%s", port)
	r := <-router
	
	r.Run(server_addr)
}

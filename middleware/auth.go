package middleware

import (
	"Book-App/controllers"
	"Book-App/models"
	"fmt"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// the jwt middleware

func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	var identityKey1 = "email"
	var identityKey2 = "Mobile"
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "test zone",
		SigningAlgorithm: "",
		Key:              []byte("secret key"),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		Authenticator:    controllers.Login,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*models.User); ok {
				return true
			}
			return false
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{identityKey1: v.Email, identityKey2: v.Mobile}
			}

			return jwt.MapClaims{}
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"code": code, "message": message})
		},
		// LoginResponse: func(*gin.Context, int, string, time.Time) {
		// },
		LogoutResponse: func(c *gin.Context, code int) {
			email := c.GetString("user_email")
			mobile := c.GetString("mobile")
			models.Rdb.Del(email)
			models.Rdb.Del(mobile)
			fmt.Println("Redis Cleared")
			c.JSON(code, gin.H{
				"message": "logged out successfully",
			})
		},
		RefreshResponse: func(*gin.Context, int, string, time.Time) {
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Email:  claims[identityKey1].(string),
				Mobile: claims[identityKey2].(string),
			}
		},
		//IdentityKey:   identityKey,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		// HTTPStatusMessageFunc: func(e error, c *gin.Context) string {
		// },
		PrivKeyFile:       "",
		PubKeyFile:        "",
		SendCookie:        true,
		SecureCookie:      false,
		CookieHTTPOnly:    true,
		CookieDomain:      "",
		SendAuthorization: true,
		DisabledAbort:     false,
		CookieName:        "",
	})
	if err != nil {
		return nil, err
	}
	return authMiddleware, nil
}

package controllers

import (
	"Book-App/models"
	"fmt"
	"html/template"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const SecretKey = "secret"

type tempUser struct {
	FirstName        string `json:"first_name" binding:"required"`
	LastName         string `json:"last_name" binding:"required"`
	Email            string `json:"email" binding:"required"`
	Password         string `json:"password" binding:"required"`
	Confirm_Password string `json:"confirm_password" binding:"required"`
}

type RedisCache struct {
	Id     int
	Email  string
	RoleId int
}

func ReturnParameterMissingError(c *gin.Context, parameter string) {
	var err = fmt.Sprintf("Required parameter %s missing.", parameter)
	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}

// @Summary Register endpoint is used for customer registeration. ( Supervisors/admin can be added only by admin. )
// @Description API Endpoint to register the user with the role of customer.
// @Router /api/v1/Register [post]
// @Tags auth
// @Accept json
// @Produce json
// @Param email formData string true "Email of the user"
// @Param first_name formData string true "First name of the user"
// @Param last_name formData string true "Last name of the user"
// @Param password formData string true "Password of the user"
// @Param confirm_password formData string true "Confirm password."
// @Success 200 {object} object
// @Failure 400,500 {object} object
func Register(c *gin.Context) {
	var tempUser tempUser
	var Role models.UserRole

	c.Request.ParseForm()
	paramList := []string{"email", "first_name", "last_name", "password", "confirm_password"}

	for _, param := range paramList {
		if c.PostForm(param) == "" {
			ReturnParameterMissingError(c, param)
		}
	}

	// if err := c.ShouldBindJSON(&tempUser); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	tempUser.Email = template.HTMLEscapeString(c.PostForm("email"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.Confirm_Password = template.HTMLEscapeString(c.PostForm("confirm_password"))

	if tempUser.Password != tempUser.Confirm_Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
		return
	}

	ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	if ispasswordstrong {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
		return
	}

	// Check if the user already exists.
	if DoesUserExist(tempUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	}

	encryptedPassword, error := HashPassword(tempUser.Password)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}

	err := models.DB.Where("role= ?", "customer").First(&Role).Error
	if err != nil {
		fmt.Println("err ", err.Error())
		return
	}

	SanitizedUser := models.User{
		FirstName:  tempUser.FirstName,
		LastName:   tempUser.LastName,
		Email:      tempUser.Email,
		Password:   encryptedPassword,
		UserRoleID: Role.Id, //This endpoint will be used only for customer registration.
		CreatedAt:  time.Now(),
		IsActive:   true,
	}

	errs := models.DB.Create(&SanitizedUser).Error
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})

}

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// redisClient := Redis.createclient()

// Login godoc
// @Summary Login endpoint is used by the user to login.
// @Description API Endpoint to register the user with the role of customer.
// @Router /api/v1/Login [post]
// @Tags auth
// @Accept json
// @Produce json
// @Param login formData login true "Credentials of the user"
// @Success 200 {object} object
// @Failure 400,500 {object} object
func Login(c *gin.Context) (interface{}, error) {
	var loginVals login

	// var User User
	var User models.User
	var count int64
	// var user models.User
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	email := loginVals.Email
	// First check if the user exist or not...
	models.DB.Where("email = ?", email).First(&User).Count(&count)
	if count == 0 {
		return nil, jwt.ErrFailedAuthentication
	}

	fmt.Println("set value ", loginVals.Email)
	err := models.Rdb.Set("email", loginVals.Email, 0).Err()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
	}

	if CheckCredentials(loginVals.Email, loginVals.Password, models.DB) {
		NewRedisCache(c, User)
		return &models.User{
			Email: email,
		}, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

// CreateSupervisor godoc
// @Summary CreateSupervisor endpoint is used by the admin role user to create a new admin or supervisor account.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/supervisor/create [post]
// @Tags supervisor
// @Accept json
// @Produce json
// @Param login formData tempUser true "Info of the user"
// @Success 200 {object} object
// @Failure 400,500 {object} object
func CreateSupervisor(c *gin.Context) {
	//fmt.Println("supervisor api hit")

	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	// Create a user with the role of supervisor.
	var tempUser tempUser
	var Role models.UserRole

	c.Request.ParseForm()
	paramList := []string{"first_name", "last_name", "email", "password", "confirm_password"}

	for _, param := range paramList {
		if c.PostForm(param) == "" {
			ReturnParameterMissingError(c, param)
		}
	}

	tempUser.Email = template.HTMLEscapeString(c.PostForm("email"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.Confirm_Password = template.HTMLEscapeString(c.PostForm("confirm_password"))

	//check if the password is strong and matches the password policy
	//length > 8, atleast 1 upper case, atleast 1 lower case, atleast 1 symbol
	// ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	// if ispasswordstrong == false {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
	// 	return
	// }

	if tempUser.Password != tempUser.Confirm_Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
	}

	// Check if the user already exists.
	if DoesUserExist(tempUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	}

	encryptedPassword, error := HashPassword(tempUser.Password)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occured."})
		return
	}

	err := models.DB.Where("role= ?", "supervisor").First(&Role).Error
	if err != nil {
		fmt.Println("err ", err.Error())
		return
	}

	SanitizedUser := models.User{
		FirstName:  tempUser.FirstName,
		LastName:   tempUser.LastName,
		Email:      tempUser.Email,
		Password:   encryptedPassword,
		UserRoleID: Role.Id, //This endpoint will be used only for customer registeration.
		CreatedAt:  time.Now(),
		IsActive:   true,
	}

	errs := models.DB.Create(&SanitizedUser).Error
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occured."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})
}

// CreateAdmin godoc
// @Summary CreateAdmin endpoint is used by the admin role user to create a new admin or supervisor account.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/admin/create [post]
// @Tags admin
// @Accept json
// @Produce json
// @Param login formData tempUser true "Info of the user"
// @Success 200 {object} object
// @Failure 400,500 {object} object
func CreateAdmin(c *gin.Context) {
	//var User models.User
	// if !IsAdmin(c) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// }
	// user_email := rdb.
	//fmt.Println("test line")
	//fmt.Println(user_email)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"error": "redis get not working",
	// 	})
	// }
	// //fmt.Println()

	// if err := models.DB.Where("email = ? AND user_role_id=1", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	var tempUser tempUser
	var Role models.UserRole

	c.Request.ParseForm()
	paramList := []string{"first_name", "last_name", "email", "password", "confirm_password"}

	for _, param := range paramList {
		if c.PostForm(param) == "" {
			ReturnParameterMissingError(c, param)
		}
	}

	tempUser.Email = template.HTMLEscapeString(c.PostForm("email"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.Confirm_Password = template.HTMLEscapeString(c.PostForm("confirm_password"))

	fmt.Println("debug start")
	fmt.Println(tempUser.FirstName, tempUser.LastName, tempUser.Password)
	fmt.Println("debug end")

	if tempUser.Password != tempUser.Confirm_Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
	}

	// check if the password is strong and matches the password policy
	// length > 8, atleast 1 upper case, atleast 1 lower case, atleast 1 symbol
	// ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	// if ispasswordstrong == false {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
	// 	return
	// }

	// Check if the user already exists.
	if DoesUserExist(tempUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	}

	encryptedPassword, error := HashPassword(tempUser.Password)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}

	err := models.DB.Where("role= ?", "admin").First(&Role).Error
	if err != nil {
		fmt.Println("err ", err.Error())
		return
	}

	SanitizedUser := models.User{
		FirstName:  tempUser.FirstName,
		LastName:   tempUser.LastName,
		Email:      tempUser.Email,
		Password:   encryptedPassword,
		UserRoleID: Role.Id, //This endpoint will be used only for customer registeration.
		CreatedAt:  time.Now(),
		IsActive:   true,
	}

	errs := models.DB.Create(&SanitizedUser).Error
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})
}

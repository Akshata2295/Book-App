package controllers

import (
	"Book-App/models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/go-redis/redis"

)

// CreateCategory godoc
// @Summary CreateCategory endpoint is used by admin to create category.
// @Description CreateCategory endpoint is used by admin to create category.
// @Router /api/v1/auth/category/create [post]
// @Tags category
// @Accept json
// @Produce json
// @Param name formData string true "name of the category"
// @Success 200 {object} models.Category
// @Failure 400,404 {object} object
func CreateCategory(c *gin.Context) {
	var existingCategory models.Category
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	// var User models.User
	// user_email, err := Rdb.HGet("user", "email").Result()
	// email := c.GetString("user_email")
	// fmt.Println(models.Rdb.HGetAll(email))
	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)
	fmt.Println(ID)
	roleId, _ := models.Rdb.HGet("user", "RoleID").Result()

	if roleId != "1" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be updated by admin user"})
		return
	}

	// // Check if the current user had admin role.
	// if err := models.DB.Where("email = ? AND user_role_id=1", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be added by admin user"})
	// }

	c.Request.ParseForm()
	var flag bool
	if c.PostForm("name") == "" {
		ReturnParameterMissingError(c, "name")
		flag = true
	}
	category_title := template.HTMLEscapeString(c.PostForm("name"))
	category_id :=template.HTMLEscapeString(c.PostForm("id"))
	categoryID,_ :=strconv.Atoi(category_id)

	if flag {
		return
	}

	// fmt.Println(category_title)
	// fmt.Println("category printed")
	// Check if the category already exists.
	
	err := models.DB.Where("title = ?", category_title).First(&existingCategory).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category already exists."})
		return
	}

	cat := models.Category{
		ID: categoryID,
		CategoryName: category_title,
		CreatedBy:    ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	fmt.Println(cat)

	err = models.DB.Create(&cat).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":   cat.ID,
		"name": cat.CategoryName,
	})


}

type ReturnedCategory struct {
	ID           int    `json:"id"`
	CategoryName string `json:"name"`
}

// ListAllCategories godoc
// @Summary ListAllCategories endpoint is used to list all the categories.
// @Description ListAllCategories endpoint is used to list all the categories.
// @Router /api/v1/auth/category/ [get]
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {array} models.Category
// @Failure 404 {object} object
func ListAllCategories(c *gin.Context) {

	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	var User models.User
	var Categories []models.Category
	var ExistingCategories []ReturnedCategory
	user_email, _ := models.Rdb.HGet("user", "email").Result()

	id, _ := models.Rdb.HGet("user", "RoleID").Result()
    
	CheckRedis(c)
	if id != "1" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Books can only be seen by Admin user"})
		return
	}

	if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	models.DB.Model(Categories).Find(&ExistingCategories)
	c.JSON(http.StatusOK, ExistingCategories)
}

// GetCategory godoc
// @Summary GetCategory endpoint is used to get info of a category..
// @Description GetCategory endpoint is used to get info of a category.
// @Router /api/v1/auth/category/:id/ [get]
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {object} models.Category
// @Failure 400,404 {object} object
func GetCategory(c *gin.Context) {
	var existingCategory models.Category

	// Check if the category already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingCategory).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category does not exists."})
		return
	}
	
	id, _ := models.Rdb.HGet("user", "RoleID").Result()

	if id != "1" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Books can only be Updated by Admin user"})
		return
	}
	// GET FROM CACHE FIRST
	c.JSON(http.StatusOK, gin.H{"category": existingCategory})
}

// UpdateCategory godoc
// @Summary UpdateCategory endpoint is used to get info of a category..
// @Description UpdateCategory endpoint is used to get info of a category.
// @Router /api/v1/auth/category/:id/ [PUT]
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {object} models.Category
// @Failure 400,404 {object} object
func UpdateCategory(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	//var User models.User
	var existingCategory models.Category
	var UpdateCategory models.Category
	//user_email, _ := Rdb.HGet("user", "email").Result()
	id, _ := models.Rdb.HGet("user", "RoleID").Result()

	if id != "1" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be updated by admin user"})
		return
	}

	// if err := models.DB.Where("email = ? AND user_role_id=1", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be updated by admin user"})
	// 	return
	// }
	// Check if the Category already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingCategory).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category does not exists."})
		return
	}

	if err := c.ShouldBindJSON(&UpdateCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&existingCategory).Updates(UpdateCategory)
}

func DeleteCategory(c *gin.Context) {
	var existingCategory models.Category
	// var User models.User
	// user_email, _ := Rdb.HGet("user", "email").Result()

	// if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Book can only be updated by supervisor user"})
	// 	return
	// }
	id, _ := models.Rdb.HGet("user", "RoleID").Result()

	if id != "1" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Books can only be deleted by admin"})
		return
	}
	// Check if the book already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingCategory).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category does not exists."})
		return
	}
	models.DB.Where("id = ?", c.Param("id")).Delete(&existingCategory)
	// GET FROM CACHE FIRST
	c.JSON(http.StatusOK, gin.H{"Success": "Category deleted"})
}





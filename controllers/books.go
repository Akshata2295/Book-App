package controllers

import (
	"Book-App/models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
)

// CreateBook godoc
// @Summary CreateBook endpoint is used by the supervisor role user to create a new book.
// @Description CreateBook endpoint is used by the supervisor role user to create a new book
// @Router /api/v1/auth/books/create [post]
// @Tags book
// @Accept json
// @Produce json
// @Param name formData string true "name of the book"
// @Param category_id formData string true "category_id of the book"
// @Success 200 {object} object
// @Failure 400 {object} object
func CreateBook(c *gin.Context) {

	var existingBook models.Book
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	//var User models.User
	var category models.Category
	
	// user_email, _ := Rdb.HGet("user", "email").Result()

	// // Check if the current user had admin role.
	// if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Book can only be added by supervisor user"})
	// 	return
	// }
	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)
	fmt.Println(ID)
	roleId, _ := models.Rdb.HGet("user", "RoleID").Result()

	if roleId != "2" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Books can only be added by supervisor"})
		return
	}

	c.Request.ParseForm()
	if c.PostForm("name") == "" {
		ReturnParameterMissingError(c, "name")	
	}
	
	if c.PostForm("category_id") == "" {
		ReturnParameterMissingError(c, "category_id")
	}

	if c.PostForm("price") == "" {
		ReturnParameterMissingError(c, "price")
	}
	
	title := template.HTMLEscapeString(c.PostForm("name"))
	category_id, _:= strconv.Atoi(template.HTMLEscapeString(c.PostForm("category_id")))
	book_id :=template.HTMLEscapeString(c.PostForm("id"))
    bookID,_ :=strconv.Atoi(book_id)
	price, _ := strconv.Atoi(template.HTMLEscapeString(c.PostForm("price")))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "can only convert string to int",
	// 	})
	// }

	//Check if the book already exists.
	err := models.DB.Where("title = ?").First(&existingBook).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book already exists."})
		return
	}

	// Check if the category exists
	err = models.DB.First(&category, category_id).Error
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "category does not exists."})
		return
	}

	book := models.Book{
		ID: bookID,
		Title:      title,
		CategoryId: category_id,
		Price:      price,
		CreatedBy:  ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	fmt.Println(book)

	err = models.DB.Create(&book).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":   book.ID,
		"name": book.Title,
		"price":book.Price,
		"category_id":book.CategoryId,
				
	})

}

// UpdateBook godoc
// @Summary UpdateBook endpoint is used by the supervisor role user to update a new book.
// @Description Updatebook endpoint is used by the supervisor role user to update a new book
// @Router /api/v1/auth/books/:id/ [PATCH]
// @Tags book
// @Accept json
// @Produce json
// @Success 200 {object} models.Book
// @Failure 400,404 {object} object
func UpdateBook(c *gin.Context) {
	var existingBook models.Book
	var updateBook models.Book
	// claims := jwt.ExtractClaims(c)
	
	// user_email, _ := claims["email"]
	//var User models.User
	// user_email, _ := Rdb.HGet("user", "email").Result()

	// // Check if the current user had admin role.
	// if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Book can only be updated by supervisor user"})
	// 	return
	// }

	id, _ := models.Rdb.HGet("user", "RoleID").Result()

	if id != "2" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Books can only be updated by supervisor"})
		return
	}

	// Check if the book already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingBook).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book does not exists."})
		return
	}

	if err := c.ShouldBindJSON(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&existingBook).Updates(updateBook)
}

type ReturnedBook struct {
	ID         int    `json:"id"`
	Title      string `json:"name"`
	CategoryId int    `json:"category_id"`
	   
}

// GetBook godoc
// @Summary GetBook endpoint is used to get info of a book..
// @Description GetBook endpoint is used to get info of a book.
// @Router /api/v1/auth/books/:id/ [get]
// @Tags book
// @Accept json
// @Produce json
// @Success 200 {object} models.Book
// @Failure 400,404 {object} object
func GetBook(c *gin.Context) {
	var existingBook models.Book

	// Check if the book already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingBook).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book does not exists."})
		return
	}

	// GET FROM CACHE FIRST
	c.JSON(http.StatusOK, gin.H{"book": existingBook})
}

// ListAllBook godoc
// @Summary ListAllBook endpoint is used to list all book.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/books/ [get]
// @Tags book
// @Accept json
// @Produce json
// @Success 200 {array} models.Book
// @Failure 404 {object} object
func ListAllBook(c *gin.Context) {

	// allBook := []models.Book{}
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	var User models.User
	var Book []models.Book
	var existingBook []ReturnedBook
	email := c.GetString("user_email")
	fmt.Println("c variable" + email)
	user_email, _ := models.Rdb.HGet("user", "email").Result()
	fmt.Println("user" + user_email)

	if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	models.DB.Model(Book).Find(&existingBook)
	c.JSON(http.StatusOK, existingBook)

}

// DeleteBook godoc
// @Summary DeleteBook endpoint is used to delete a book.
// @Description DeleteBook endpoint is used to delete a book.
// @Router /api/v1/auth/books/delete/:id/ [delete]
// @Tags book
// @Accept json
// @Produce json
// @Success 200 {object} models.Book
// @Failure 400,404 {object} object
func DeleteBook(c *gin.Context) {
	var existingBook models.Book
	// var User models.User
	// user_email, _ := Rdb.HGet("user", "email").Result()

	// if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Book can only be updated by supervisor user"})
	// 	return
	// }
	id, _ := models.Rdb.HGet("user", "RoleID").Result()

	if id != "2" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Books can only be deleted by supervisor"})
		return
	}
	// Check if the book already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingBook).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book does not exists."})
		return
	}
	models.DB.Where("id = ?", c.Param("id")).Delete(&existingBook)
	// GET FROM CACHE FIRST
	c.JSON(http.StatusOK, gin.H{"Success": "Book deleted"})
}




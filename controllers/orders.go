package controllers

import (
	"fmt"
	"Book-App/models"
	"html/template"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// AddToCart godoc
// @Summary AddToCart endpoint is used to add the book to the cart.
// @Description AddToCart endpoint is used to add the book to the cart.
// @Router /api/v1/auth/cart/add [post]
// @Tags book
// @Accept json
// @Produce json
// @Param title formData string true "title of the book"
// @Success 200 {object} models.Cart
// @Failure 400,404 {object} object
func AddToCart(c *gin.Context) {

	//var existingBook models.Book
	claims := jwt.ExtractClaims(c)
	user_email := claims["email"]
	var User models.User
	var Book models.Book

	// Check if the current user had admin role.
	if err := models.DB.Where("email = ? AND user_role_id=3", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book can only be added to cart by user"})
		return
	}

	c.Request.ParseForm()

	if c.PostForm("title") == "" {
		ReturnParameterMissingError(c, "title")
	}
	// if c.PostForm("category_id") == "" {
	// 	ReturnParameterMissingError(c, "category_id")
	// 	return
	// }

	title := template.HTMLEscapeString(c.PostForm("title"))
	//fmt.Println(name)
	// category_id := template.HTMLEscapeString(c.PostForm("category_id"))
	// price, err := strconv.Atoi(template.HTMLEscapeString(c.PostForm("price")))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "can only convert string to int",
	// 	})
	// }

	//Check if the book already exists.
	err := models.DB.Where("title = ?", title).First(&Book).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book does not exist."})
		return
	}
	fmt.Println(Book)
	cart := models.Cart{
		// BookName: Book.Title,
		// Price:    Book.Price,
		User: User,
		Book: Book,
	}

	err = models.DB.Create(&cart).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"name":    title,
	})

}

// ViewCart godoc
// @Summary ViewCart endpoint is used to list all book.
// @Description API Endpoint to view the cart items.
// @Router /api/v1/auth/cart/view [get]
// @Tags book
// @Accept json
// @Produce json
// @Success 200 {object} models.Cart
// @Failure 400,404 {object} object
func ViewCart(c *gin.Context) {

	// allBook := []models.Book{}
	claims := jwt.ExtractClaims(c)
	user_email := claims["email"]
	var User models.User
	//var tempBook models.Book
	//bookID := Cart.BookID
	//var Book []models.Book
	var Cart []models.Cart

	if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userid := User.ID
	
	if err := models.DB.Where("user_id = ?", userid).Find(&Cart).Error; err != nil {
		c.JSON(http.StatusFound, gin.H{
			"message": "Cart empty",
		})
	}

		var title string
		var price int
//JOINING TABLE USING RAW QUERY
	// rows, err := models.DB.Raw("SELECT title,price FROM books INNER JOIN carts on books.id=carts.book_id").Rows()
	// if err != nil {
	// 	c.JSON(404, gin.H{
	// 		"error": err,
	// 	})
	// }
	// for rows.Next() {
	// 	rows.Scan(&title, &price)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"title": title,
	// 		"price": price,
	// 	})
	// }

	//JOINING TABLE USING GORM JOINS
	rows, err := models.DB.Table("books").Select("books.title", "books.price").Joins("inner join carts on carts.book_id = books.id and carts.user_id=? ", userid).Rows()
	if err != nil {
		c.JSON(404, gin.H{
			"error": err,
		})
	}

	// if rows == nil {
	// 	c.JSON(http.StatusNoContent, gin.H{
	// 		"message": "cart is empty",
	// 	})
	// }

	for rows.Next() {
		rows.Scan(&title, &price)
		c.JSON(http.StatusOK, gin.H{
			"title": title,
			"price": price,
		})
	}

	// INDIVIDUAL QUERY SEARCHING FOR EACH CART ELEMENT
	// for _, cart := range Cart {
	// 	var tempBook models.Book
	// 	bookId := cart.BookID
	// 	//fmt.Println(bookId)
	// 	//bookid := append(bookId, cart.BookID)
	// 	if err := models.DB.Where("ID = ?", bookId).Find(&tempBook).Error; err != nil {
	// 		c.JSON(http.StatusFound, gin.H{
	// 			"message": "error searching book",
	// 		})
	// 	}
	// 	Book = append(Book, tempBook)
	// }

	// // //models.DB.Model(Cart).Find(&existingBook)
	// c.JSON(http.StatusOK, Book)
}



// DeleteFromCart godoc
// @Summary DeleteFromCart endpoint is used to delete book from cart.
// @Description DeleteFromCart endpoint is used to delete book from cart.
// @Router /api/v1/auth/cart/delete/:id/ [delete]
// @Tags book
// @Accept json
// @Produce json
// @Success 200 {object} models.Cart
// @Failure 400,404 {object} object
func DeleteFromCart(c *gin.Context) {
	// var existingBook models.Book
	// var updateBook models.Book
	var Cart models.Cart
	claims := jwt.ExtractClaims(c)
	user_email := claims["email"]
	var User models.User

	if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userId := User.ID
	// Check if the book already exists.
	err := models.DB.Where("book_id = ? AND user_id = ?", c.Param("id"), userId).Delete(&Cart).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book not in cart."})
		return
	}
	//err = models.DB.Delete(&Cart).Error
	// models.DB.Where("book_id = ?", c.paaram)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"Success": "Book removed from cart",
		})
	}
	// if err := c.ShouldBindJSON(&updateBook); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	//models.DB.Model(&existingBook).Updates(updateBook)

}

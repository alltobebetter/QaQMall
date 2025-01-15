package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"qaqmall/handlers"
	"qaqmall/models"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/qaqmall?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{}, &models.CartItem{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	userHandler := handlers.NewUserHandler(db)
	productHandler := handlers.NewProductHandler(db)
	cartHandler := handlers.NewCartHandler(db)

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	r.GET("/products", productHandler.ListProducts)
	r.GET("/products/search", productHandler.SearchProducts)
	r.GET("/products/:id", productHandler.GetProduct)
	r.POST("/products", productHandler.CreateProduct)
	r.PUT("/products/:id", productHandler.UpdateProduct)
	r.DELETE("/products/:id", productHandler.DeleteProduct)

	r.GET("/cart", cartHandler.GetCart)
	r.POST("/cart/add", cartHandler.AddItem)
	r.PUT("/cart/items/:id", cartHandler.UpdateItem)
	r.DELETE("/cart/items/:id", cartHandler.RemoveItem)
	r.DELETE("/cart", cartHandler.EmptyCart)

	if err := r.Run(":8888"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

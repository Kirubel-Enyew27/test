package main

import (
	"shopping-cart/config"
	"shopping-cart/data"
	"shopping-cart/db"
	"shopping-cart/handler"
	"shopping-cart/service"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	config.ConnectDB()
	config.CreateTables()
}

func main() {
	dbQueries := db.New(config.DB)
	cartRepository := data.NewCartRepo(dbQueries)
	cartService := service.NewCartService(cartRepository)
	cartHandler := handler.NewCartHandler(cartService)

	router := gin.Default()

	router.POST("/products/add", cartHandler.AddProduct)
	router.POST("/cart/add", cartHandler.AddItem)
	router.DELETE("/cart/remove/:itemID", cartHandler.RemoveItem)
	router.PUT("/cart/update", cartHandler.UpdateItemQuantity)
	router.POST("/cart/discount", cartHandler.ApplyDiscount)
	router.GET("/cart", cartHandler.ViewCart)
	router.POST("/cart/checkout", cartHandler.Checkout)

	router.Run(":8080")
}

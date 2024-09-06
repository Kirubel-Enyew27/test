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
	cartrHandler := handler.NewCartHandler(cartService)

	router := gin.Default()

	router.POST("/cart/add", cartrHandler.AddItem)
	// router.DELETE("/cart/remove/:productID", handler.RemoveItem)
	// router.PUT("/cart/update", handler.UpdateItemQuantity)
	// router.POST("/cart/discount", handler.ApplyDiscount)
	// router.GET("/cart", handler.ViewCart)
	// router.POST("/cart/checkout", handler.Checkout)

	router.Run(":8080")
}

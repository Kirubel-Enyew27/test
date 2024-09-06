package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    router.POST("/cart/add", handler.AddItem)
    router.DELETE("/cart/remove/:productID", handler.RemoveItem)
    router.PUT("/cart/update", handler.UpdateItemQuantity)
    router.POST("/cart/discount", handler.ApplyDiscount)
    router.GET("/cart", handler.ViewCart)
    router.POST("/cart/checkout", handler.Checkout)

    router.Run(":8080")
}



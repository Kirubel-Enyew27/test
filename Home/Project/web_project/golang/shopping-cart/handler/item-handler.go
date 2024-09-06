package handler

import (
	"context"
	"net/http"
	"shopping-cart/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service service.CartServiceInterface
}

func NewCartHandler(service service.CartServiceInterface) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) AddItem(c *gin.Context) {
	var req struct {
		ItemID   int32   `json:"item_id"`
		ItemName string  `json:"item_name"`
		Price    float64 `json:"price"`
		Quantity int32   `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.service.AddItem(context.Background(), req.ItemID, req.ItemName, req.Price, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
}

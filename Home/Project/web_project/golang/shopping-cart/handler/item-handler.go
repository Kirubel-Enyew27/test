package handler

import (
	"context"
	"net/http"
	"shopping-cart/service"
	"strconv"

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

	items, err := h.service.AddItem(context.Background(), req.ItemID, req.ItemName, req.Price, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart", "total unique items added": items})
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = h.service.RemoveItem(context.Background(), itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

func (h *CartHandler) UpdateItemQuantity(c *gin.Context) {
	var req struct {
		ItemID   int32 `json:"item_id"`
		Quantity int32 `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.service.UpdateItemQuantity(context.Background(), req.ItemID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item quantity updated"})
}

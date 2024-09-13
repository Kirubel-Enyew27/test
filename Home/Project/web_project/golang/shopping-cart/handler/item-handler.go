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

func (h *CartHandler) AddProduct(c *gin.Context) {
	var req struct {
		ProductID   int32   `json:"product_id"`
		ProductName string  `json:"product_name"`
		Price       float64 `json:"price"`
		Stock       int32   `json:"stock"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.service.AddProduct(context.Background(), req.ProductID, req.ProductName, req.Price, req.Stock)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added successfully"})
}

func (h *CartHandler) AddItem(c *gin.Context) {
	var req struct {
		ProductID int32 `json:"product_id"`
		Quantity  int32 `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	items, err := h.service.AddItem(context.Background(), req.ProductID, req.Quantity)
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

func (h *CartHandler) ApplyDiscount(c *gin.Context) {
	var req struct {
		Discount float64 `json:"discount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.service.ApplyDiscount(context.Background(), req.Discount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Discount applied"})
}

func (h *CartHandler) ViewCart(c *gin.Context) {
	items, err := h.service.ViewCart(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cart": items})
}

func (h *CartHandler) Checkout(c *gin.Context) {
	err := h.service.Checkout(context.Background())
	if err != nil {
		if err.Error() == "cart is empty" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Checkout successful"})
}

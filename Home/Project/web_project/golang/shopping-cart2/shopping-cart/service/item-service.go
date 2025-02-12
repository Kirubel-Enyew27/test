package service

import (
	"context"
	"errors"
	"fmt"
	"shopping-cart/data"
	"shopping-cart/db"
)

type DiscountType string

const (
	Percentage DiscountType = "percentage"
	FlatRate   DiscountType = "flat_rate"
)

type CartServiceInterface interface {
	AddProduct(ctx context.Context, productID int32, productName string, price float64, stock int32) error
	AddItem(ctx context.Context, productID int32, quantity int32) (int, error)
	RemoveItem(ctx context.Context, itemID int) error
	UpdateItemQuantity(ctx context.Context, itemID int32, quantity int32) error
	ApplyDiscount(ctx context.Context, discount float64, discountType DiscountType) error
	ViewCart(ctx context.Context) ([]db.CartItem, float64, error)
	Checkout(ctx context.Context) error
}

type CartService struct {
	repo data.CartRepoInterface
}

func NewCartService(repo data.CartRepoInterface) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) AddProduct(ctx context.Context, productID int32, productName string, price float64, stock int32) error {
	if price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if productName == "" {
		return errors.New("product name must not be empty")
	}

	return s.repo.AddProduct(ctx, productID, productName, price, stock)
}

func (s *CartService) AddItem(ctx context.Context, productID int32, quantity int32) (int, error) {
	if quantity <= 0 {
		return 0, errors.New("quantity must be greater than zero")
	}

	items, err := s.repo.AddCartItem(ctx, productID, quantity)
	if err != nil {
		return items, err
	}

	return items, nil
}

func (s *CartService) RemoveItem(ctx context.Context, itemID int) error {
	isItemExist, err := s.repo.FindItem(ctx, itemID)
	if err != nil {
		return err
	}
	if !isItemExist {
		return errors.New("item not found")
	}

	err = s.repo.RemoveItem(ctx, itemID)
	if err != nil {
		return err
	}

	return nil

}

func (s *CartService) UpdateItemQuantity(ctx context.Context, itemID int32, quantity int32) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	return s.repo.UpdateItemQuantity(ctx, itemID, quantity)
}

func (s *CartService) ApplyDiscount(ctx context.Context, discount float64, discountType DiscountType) error {
	if discountType == Percentage {
		if discount <= 0 || discount > 100 {
			return errors.New("percentage discount must be between 0 and 100")
		}
		err := s.repo.ApplyPercentageDiscountToCart(ctx, discount)
		if err != nil {
			return err
		}
	} else if discountType == FlatRate {
		if discount < 0 {
			return errors.New("flat rate discount must be non-negative")
		}
		err := s.repo.ApplyFlatRateDiscountToCart(ctx, discount)
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid discount type")
	}

	return nil
}

func (s *CartService) ViewCart(ctx context.Context) ([]db.CartItem, float64, error) {
	items, err := s.repo.ViewCart(ctx)
	totlalPrice := 0.0

	for _, item := range items {
		totlalPrice += float64(item.Quantity) * item.Price
	}

	return items, totlalPrice, err
}

func (s *CartService) Checkout(ctx context.Context) error {
	items, err := s.repo.ViewCart(ctx)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}
	return s.repo.Checkout(ctx)
}

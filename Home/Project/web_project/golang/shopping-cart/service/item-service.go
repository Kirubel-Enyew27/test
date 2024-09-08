package service

import (
	"context"
	"errors"
	"shopping-cart/data"
)

type CartServiceInterface interface {
	AddItem(ctx context.Context, itemID int32, itemName string, price float64, quantity int32) (int, error)
	RemoveItem(ctx context.Context, itemID int) error
	UpdateItemQuantity(ctx context.Context, itemID int32, quantity int32) error
	ApplyDiscount(ctx context.Context, discount float64) error
}

type CartService struct {
	repo data.CartRepoInterface
}

func NewCartService(repo data.CartRepoInterface) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) AddItem(ctx context.Context, itemID int32, itemName string, price float64, quantity int32) (int, error) {
	if quantity <= 0 {
		return 0, errors.New("quantity must be greater than zero")
	}

	items, err := s.repo.AddCartItem(ctx, itemID, itemName, price, quantity)
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

func (s *CartService) ApplyDiscount(ctx context.Context, discount float64) error {
	if discount <= 0 || discount > 100 {
		return errors.New("invalid discount value")
	}
	return s.repo.ApplyDiscount(ctx, discount)
}

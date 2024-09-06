package service

import (
	"context"
	"errors"
	"shopping-cart/data"
)

type CartServiceInterface interface {
	AddItem(ctx context.Context, itemID int32, itemName string, price float64, quantity int32) error
	RemoveItem(ctx context.Context, itemID int) error
}

type CartService struct {
	repo data.CartRepoInterface
}

func NewCartService(repo data.CartRepoInterface) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) AddItem(ctx context.Context, itemID int32, itemName string, price float64, quantity int32) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	err := s.repo.AddCartItem(ctx, itemID, itemName, price, quantity)
	if err != nil {
		return err
	}

	return nil
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

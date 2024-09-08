package data

import (
	"context"
	"errors"
	"shopping-cart/db"
)

type CartRepoInterface interface {
	AddCartItem(ctx context.Context, itemID int32, itemName string, price float64, quantity int32) (int, error)
	RemoveItem(ctx context.Context, itemID int) error
	RemoveAllItem(ctx context.Context) error
	FindItem(ctx context.Context, itemID int) (bool, error)
	UpdateItemQuantity(ctx context.Context, itemID int32, quantity int32) error
	ApplyDiscount(ctx context.Context, discount float64) error
}

type CartRepo struct {
	dbQueries *db.Queries
}

func NewCartRepo(dbQueries *db.Queries) *CartRepo {
	return &CartRepo{
		dbQueries: dbQueries,
	}
}

func (r *CartRepo) AddCartItem(ctx context.Context, itemID int32, itemName string, price float64, quantity int32) (int, error) {
	count, err := r.dbQueries.CountUniqueItemsInCart(ctx)
	if err != nil {
		return int(count), err
	}

	if count >= 10 {
		return int(count), errors.New("cannot add more than 10 unique items to the cart")

	}

	params := db.AddCartItemParams{
		ItemID:   itemID,
		ItemName: itemName,
		Price:    price,
		Quantity: quantity,
	}

	err = r.dbQueries.AddCartItem(ctx, params)
	if err != nil {
		return int(count), err
	}

	return int(count + 1), nil
}

func (r *CartRepo) RemoveItem(ctx context.Context, itemID int) error {
	err := r.dbQueries.RemoveItem(ctx, int32(itemID))
	if err != nil {
		return errors.New("failed to remove item from the cart")
	}
	return nil
}

func (r *CartRepo) RemoveAllItem(ctx context.Context) error {
	err := r.dbQueries.RemoveAllItem(ctx)
	if err != nil {
		return errors.New("failed to remove item from the cart")
	}
	return nil
}

func (r *CartRepo) FindItem(ctx context.Context, itemID int) (bool, error) {
	exists, err := r.dbQueries.FindItemInCart(ctx, int32(itemID))
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *CartRepo) UpdateItemQuantity(ctx context.Context, itemID int32, quantity int32) error {
	params := db.UpdateItemQuantityParams{
		ItemID:   itemID,
		Quantity: quantity,
	}
	err := r.dbQueries.UpdateItemQuantity(ctx, params)
	if err != nil {
		return errors.New("failed to update item quantity")
	}
	return nil
}

func (r *CartRepo) ApplyDiscount(ctx context.Context, discount float64) error {
	err := r.dbQueries.ApplyDiscountToCart(ctx, discount)
	if err != nil {
		return errors.New("failed to apply discount to the cart")
	}
	return nil
}

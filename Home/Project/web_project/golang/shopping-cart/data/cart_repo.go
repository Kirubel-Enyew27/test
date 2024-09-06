package data

import (
	"context"
	"errors"
	"shopping-cart/db"
)

type CartRepoInterface interface {
	AddCartItem(ctx context.Context, itemID int32, ItemName string, Price float64, quantity int32) error
}

type CartRepo struct {
	dbQueries *db.Queries
}

func NewCartRepo(dbQueries *db.Queries) *CartRepo {
	return &CartRepo{
		dbQueries: dbQueries,
	}
}

func (r CartRepo) AddCartItem(ctx context.Context, itemID int32, itemName string, price float64, quantity int32) error {
	count, err := r.dbQueries.CountUniqueItemsInCart(ctx)
	if err != nil {
		return err
	}

	if count >= 10 {
		return errors.New("cannot add more than 10 unique items to the cart")
	}

	params := db.AddCartItemParams{
		ItemID:   itemID,
		ItemName: itemName,
		Price:    price,
		Quantity: quantity,
	}

	err = r.dbQueries.AddCartItem(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

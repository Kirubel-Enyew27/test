package data

import (
	"context"
	"errors"
	"shopping-cart/db"
)

type CartRepoInterface interface {
	AddProduct(ctx context.Context, productID int32, productName string, price float64, stock int32) error
	AddCartItem(ctx context.Context, productID int32, quantity int32) (int, error)
	RemoveItem(ctx context.Context, itemID int) error
	RemoveAllItem(ctx context.Context) error
	FindItem(ctx context.Context, itemID int) (bool, error)
	UpdateItemQuantity(ctx context.Context, itemID int32, quantity int32) error
	ApplyDiscount(ctx context.Context, discount float64) error
	ViewCart(ctx context.Context) ([]db.CartItem, error)
	Checkout(ctx context.Context) error
}

type CartRepo struct {
	dbQueries *db.Queries
}

func NewCartRepo(dbQueries *db.Queries) *CartRepo {
	return &CartRepo{
		dbQueries: dbQueries,
	}
}

func (r *CartRepo) AddProduct(ctx context.Context, productID int32, productName string, price float64, stock int32) error {
	params := db.AddProductParams{
		ProductID:   productID,
		ProductName: productName,
		Price:       price,
		Stock:       stock,
	}

	err := r.dbQueries.AddProduct(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *CartRepo) GetProduct(ctx context.Context, productID int32) (*db.Product, error) {
	product, err := r.dbQueries.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *CartRepo) AddCartItem(ctx context.Context, productID int32, quantity int32) (int, error) {

	product, err := r.dbQueries.GetProductByID(ctx, productID)
	if err != nil {
		return 0, errors.New("product does not exist")
	}

	count, err := r.dbQueries.CountUniqueItemsInCart(ctx)
	if err != nil {
		return int(count), err
	}

	if count >= 10 {
		return int(count), errors.New("cannot add more than 10 unique items to the cart")

	}

	if quantity > product.Stock {
		return int(count), errors.New("quantity exceeds available stock")
	}

	params := db.AddCartItemParams{
		ItemID:   product.ProductID,
		ItemName: product.ProductName,
		Price:    product.Price,
		Quantity: quantity,
	}

	err = r.dbQueries.AddCartItem(ctx, params)
	if err != nil {
		return int(count), err
	}

	stock := product.Stock - quantity
	productParams := db.UpdateProductStockParams{
		ProductID: productID,
		Stock:     stock,
	}

	err = r.dbQueries.UpdateProductStock(ctx, productParams)
	if err != nil {
		return int(count + 1), err
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

func (r *CartRepo) RemoveAllProduct(ctx context.Context) error {
	err := r.dbQueries.RemoveAllProduct(ctx)
	if err != nil {
		return errors.New("failed to remove product")
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

func (r *CartRepo) ViewCart(ctx context.Context) ([]db.CartItem, error) {
	items, err := r.dbQueries.ViewCart(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve cart items")
	}
	return items, nil
}

func (r *CartRepo) Checkout(ctx context.Context) error {
	err := r.dbQueries.CheckoutCart(ctx)
	if err != nil {
		return errors.New("failed to checkout cart")
	}
	return nil
}

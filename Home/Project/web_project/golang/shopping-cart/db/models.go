// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
)

type CartDiscount struct {
	DiscountID int32
	Percentage sql.NullString
	FlatRate   sql.NullString
}

type CartItem struct {
	ItemID    int32
	ProductID sql.NullInt32
	Quantity  int32
}

type Product struct {
	ProductID int32
	Name      string
	Price     string
	Stock     int32
}

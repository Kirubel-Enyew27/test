package steps

import (
	"context"
	"fmt"
	"shopping-cart/config"
	"shopping-cart/data"
	"shopping-cart/db"
)

func ClearDB() {
	cartRepo := data.NewCartRepo(db.New(config.DB))

	err := cartRepo.RemoveAllProduct(context.Background())
	if err != nil {
		fmt.Printf("failed to clean up products, %v", err)
	}

	err = cartRepo.RemoveAllItem(context.Background())
	if err != nil {
		fmt.Printf("failed to clean up cart items, %v", err)
	}
}

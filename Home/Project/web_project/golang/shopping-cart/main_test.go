package main

import (
	"context"
	"fmt"
	"os"
	"shopping-cart/config"
	"shopping-cart/data"
	"shopping-cart/db"
	"shopping-cart/steps"
	"testing"

	"github.com/cucumber/godog"
)

func init() {
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

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",
	}

	status := godog.TestSuite{
		Name:                "Shopping Cart",
		ScenarioInitializer: steps.InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

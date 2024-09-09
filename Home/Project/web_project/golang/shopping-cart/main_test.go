package main

import (
	"os"
	"shopping-cart/steps"
	"testing"

	"github.com/cucumber/godog"
)

func init() {
	steps.ClearDB()
}

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",
	}

	status := godog.TestSuite{
		Name: "Shopping Cart",
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			steps.InitializeAddItemScenario(ctx)
			steps.InitializeRemoveItemScenario(ctx)
		},
		Options: &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

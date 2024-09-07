package main

import (
	"os"
	"shopping-cart/steps"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",
		Output: colors.Colored(os.Stdout),
	}

	suite := godog.TestSuite{
		Name: "user-registration",
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			steps.InitializeScenario(ctx)
		},
		Options: &opts,
	}

	status := suite.Run()

	os.Exit(status)
}

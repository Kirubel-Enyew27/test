package steps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"shopping-cart/config"
	"shopping-cart/data"
	"shopping-cart/db"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
)

func iApplyADiscountOfToTheCart(discount_value float64) error {
	reqBody := struct {
		Discount float64 `json:"discount"`
	}{
		Discount: discount_value,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req := httptest.NewRequest(http.MethodPost, "/discount", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	cartHandler.ApplyDiscount(c)

	lastResponse = recorder

	if recorder.Code != http.StatusOK {
		return fmt.Errorf("expected status code 200 but got %d: %s", recorder.Code, recorder.Body.String())
	}

	return nil
}

func theExpectedPriceShouldBe(expected_price float64) error {
	cartRepo := data.NewCartRepo(db.New(config.DB))
	cart_items, err := cartRepo.ViewCart(context.Background())
	if err != nil {
		return fmt.Errorf("failed to find the item, %v", err)
	}

	for _, item := range cart_items {
		if item.Price != expected_price {
			return fmt.Errorf("expected price %v, but got %v", expected_price, item.Price)
		}
	}

	return nil
}

func InitializeApplyDiscountScenario(ctx *godog.ScenarioContext) {
	SetupScenario()
	ctx.Step(`^I have added a product with ID "([^"]*)", name "([^"]*)", price "([^"]*)", and stock "([^"]*)" is available$`, iHaveAddeAProductWithIDNamePriceAndStockIsAvailable)
	ctx.Step(`^I have added "([^"]*)" of product "([^"]*)" to the cart$`, iHaveAddedOfProductToTheCart)
	ctx.Step(`^I apply a discount of ([\d.]+)% to the cart$`, iApplyADiscountOfToTheCart)
	ctx.Step(`^the expected price should be ([\d.]+)$`, theExpectedPriceShouldBe)
}

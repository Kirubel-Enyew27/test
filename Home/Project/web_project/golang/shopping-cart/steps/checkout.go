package steps

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
)

func iCheckoutTheItemsInTheCart() error {
	req := httptest.NewRequest(http.MethodPost, "/checkout", nil)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	cartHandler.Checkout(c)

	lastResponse = recorder

	if recorder.Code != http.StatusOK {
		return fmt.Errorf("expected status code 200 but got %d: %s", recorder.Code, recorder.Body.String())
	}

	return nil
}

func theItemsShouldNoLongerBeInTheCart() error {
	req := httptest.NewRequest(http.MethodPost, "/checkout", nil)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	cartHandler.Checkout(c)

	lastResponse = recorder

	if lastResponse.Code != http.StatusBadRequest {
		return fmt.Errorf("expected status code 400 but got %d: %s", lastResponse.Code, lastResponse.Body.String())
	}

	return nil
}

func InitializeCheckoutScenario(ctx *godog.ScenarioContext) {
	SetupScenario()
	ctx.Step(`^I have added a product with ID "([^"]*)", name "([^"]*)", price "([^"]*)", and stock "([^"]*)" is available$`, iHaveAddeAProductWithIDNamePriceAndStockIsAvailable)
	ctx.Step(`^I have added "([^"]*)" of product "([^"]*)" to the cart$`, iHaveAddedOfProductToTheCart)
	ctx.Step(`^I checkout the items in the cart$`, iCheckoutTheItemsInTheCart)
	ctx.Step(`^the items should no longer be in the cart$`, theItemsShouldNoLongerBeInTheCart)
}

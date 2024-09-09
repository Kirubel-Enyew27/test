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
	"strconv"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
)

func iUpdateOfProductTo(productID, newQuantity string) error {
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		return err
	}

	quantityInt, err := strconv.Atoi(newQuantity)
	if err != nil {
		return err
	}

	reqBody := struct {
		ProductID int32 `json:"product_id"`
		Quantity  int32 `json:"quantity"`
	}{
		ProductID: int32(productIDInt),
		Quantity:  int32(quantityInt),
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req := httptest.NewRequest(http.MethodPut, "/update-item", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	cartHandler.UpdateItemQuantity(c)

	lastResponse = recorder

	if recorder.Code != http.StatusOK {
		return fmt.Errorf("expected status code 200 but got %d: %s", recorder.Code, recorder.Body.String())
	}

	return nil
}

func theQuantityOfTheProductInTheCartShouldBeUpdatedTo(newQuantity string) error {

	newQuantityInt, err := strconv.Atoi(newQuantity)
	if err != nil {
		return err
	}

	updatedQuantity := int32(newQuantityInt)

	cartRepo := data.NewCartRepo(db.New(config.DB))
	cart_items, err := cartRepo.ViewCart(context.Background())
	if err != nil {
		return fmt.Errorf("failed to find the item, %v", err)
	}

	for _, item := range cart_items {
		if item.Quantity != updatedQuantity {
			return fmt.Errorf("expected quantity %v, but got %v", updatedQuantity, item.Quantity)
		}
	}

	return nil
}

func InitializeUpdateItemScenario(ctx *godog.ScenarioContext) {
	SetupScenario()
	ctx.Step(`^I have added a product with ID "([^"]*)", name "([^"]*)", price "([^"]*)", and stock "([^"]*)" is available$`, iHaveAddeAProductWithIDNamePriceAndStockIsAvailable)
	ctx.Step(`^I have added "([^"]*)" of product "([^"]*)" to the cart$`, iHaveAddedOfProductToTheCart)
	ctx.Step(`^I update the quantity of product "([^"]*)" to "([^"]*)"$`, iUpdateOfProductTo)
	ctx.Step(`^the quantity of the product in the cart should be updated to "([^"]*)"$`, theQuantityOfTheProductInTheCartShouldBeUpdatedTo)
}

package steps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
)

func iViewTheCart() error {
	req := httptest.NewRequest(http.MethodGet, "/view", nil)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	lastResponse = recorder

	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	cartHandler.ViewCart(c)

	if recorder.Code != http.StatusOK {
		return fmt.Errorf("expected status code 200 but got %d: %s", recorder.Code, recorder.Body.String())
	}

	return nil
}

func iShouldSeeTheItemWithIDNamePriceAndQuantity(itemID int32, itemName string, price float64, quantity int32) error {
	var response map[string]interface{}
	if err := json.Unmarshal(lastResponse.Body.Bytes(), &response); err != nil {
		return err
	}

	cartItems, ok := response["cart"].([]interface{})
	if !ok {
		return fmt.Errorf("couldn't find items in the cart")
	}

	for _, item := range cartItems {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("item is not a valid map")
		}

		idFloat, ok := itemMap["ItemID"].(float64)
		if !ok {
			return fmt.Errorf("itemID is not a float64")
		}
		id := int32(idFloat)

		name, ok := itemMap["ItemName"].(string)
		if !ok {
			return fmt.Errorf("itemName is not a string")
		}

		priceValueFloat, ok := itemMap["Price"].(float64)
		if !ok {
			return fmt.Errorf("price is not a float64")
		}
		priceValue := priceValueFloat

		quantityFloat, ok := itemMap["Quantity"].(float64)
		if !ok {
			return fmt.Errorf("quantity is not a float64")
		}
		quantityValue := int32(quantityFloat)

		if id != itemID {
			return fmt.Errorf("expected id %v but got %v", itemID, id)
		} else if name != itemName {
			return fmt.Errorf("expected name %v but got %v", itemName, name)
		} else if priceValue != price {
			return fmt.Errorf("expected price %v but got %v", price, priceValue)
		} else if quantityValue != quantity {
			return fmt.Errorf("expected quantity %v but got %v", quantity, quantityValue)
		}
	}

	return nil
}

func InitializeViewItemsScenario(ctx *godog.ScenarioContext) {
	SetupScenario()
	ctx.Step(`^I have added a product with ID "([^"]*)", name "([^"]*)", price "([^"]*)", and stock "([^"]*)" is available$`, iHaveAddeAProductWithIDNamePriceAndStockIsAvailable)
	ctx.Step(`^I have added "([^"]*)" of product "([^"]*)" to the cart$`, iHaveAddedOfProductToTheCart)
	ctx.Step(`^I view the cart$`, iViewTheCart)
	ctx.Step(`^I should see the item with ID (\d+), name "([^"]*)", price (\d+), and quantity (\d+) :$`, iShouldSeeTheItemWithIDNamePriceAndQuantity)

}

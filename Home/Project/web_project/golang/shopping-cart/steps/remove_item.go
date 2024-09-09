package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
)

func iHaveAddeAProductWithIDNamePriceAndStockIsAvailable(productID, productName, price, stock string) error {
	ClearDB()
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		return err
	}

	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return err
	}

	stockInt, err := strconv.Atoi(stock)
	if err != nil {
		return err
	}

	reqBody := struct {
		ProductID   int32   `json:"product_id"`
		ProductName string  `json:"product_name"`
		Price       float64 `json:"price"`
		Stock       int32   `json:"stock"`
	}{
		ProductID:   int32(productIDInt),
		ProductName: productName,
		Price:       priceFloat,
		Stock:       int32(stockInt),
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req := httptest.NewRequest(http.MethodPost, "/add-product", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	cartHandler.AddProduct(c)

	if recorder.Code != http.StatusOK {
		return fmt.Errorf("expected status code 200 but got %d: %s", recorder.Code, recorder.Body.String())
	}

	return nil
}

func iHaveAddedOfProductToTheCart(quantity, productID string) error {
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		return err
	}

	quantityInt, err := strconv.Atoi(quantity)
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

	req := httptest.NewRequest(http.MethodPost, "/add-item", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	cartHandler.AddItem(c)

	lastResponse = recorder

	if recorder.Code != http.StatusOK {
		return fmt.Errorf("expected status code 200 but got %d: %s", recorder.Code, recorder.Body.String())
	}

	return nil
}

func iRemoveProductFromTheCart(productID string) error {
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		return err
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/remove-item/%d", productIDInt), nil)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	c.Params = gin.Params{
		{Key: "itemID", Value: productID},
	}

	cartHandler.RemoveItem(c)

	lastResponse = recorder

	if recorder.Code != http.StatusOK {
		return fmt.Errorf("expected status code 200 but got %d: %s", recorder.Code, recorder.Body.String())
	}

	return nil
}

func theProductShouldNoLongerBeInTheCart(productID string) error {
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		return err
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/remove-item/%d", productIDInt), nil)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	c.Params = gin.Params{
		{Key: "itemID", Value: productID},
	}

	cartHandler.RemoveItem(c)

	lastResponse = recorder

	if recorder.Code != http.StatusInternalServerError {
		return fmt.Errorf("expected status code 500 but got %d: %s", recorder.Code, recorder.Body.String())
	}

	ClearDB()
	return nil
}

func InitializeRemoveItemScenario(ctx *godog.ScenarioContext) {
	SetupScenario()
	ctx.Step(`^I have added a product with ID "([^"]*)", name "([^"]*)", price "([^"]*)", and stock "([^"]*)" is available$`, iHaveAddeAProductWithIDNamePriceAndStockIsAvailable)
	ctx.Step(`^I have added "([^"]*)" of product "([^"]*)" to the cart$`, iHaveAddedOfProductToTheCart)
	ctx.Step(`^I remove product "([^"]*)" from the cart$`, iRemoveProductFromTheCart)
	ctx.Step(`^the product "([^"]*)" should no longer be in the cart$`, theProductShouldNoLongerBeInTheCart)

}

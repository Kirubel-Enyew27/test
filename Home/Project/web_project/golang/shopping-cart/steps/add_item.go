package steps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"shopping-cart/config"
	"shopping-cart/data"
	"shopping-cart/db"
	"shopping-cart/handler"
	"shopping-cart/service"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
)

var cartHandler *handler.CartHandler
var lastResponse *httptest.ResponseRecorder

func SetupScenario() {
	cartRepo := data.NewCartRepo(db.New(config.DB))
	cartService := service.NewCartService(cartRepo)
	cartHandler = handler.NewCartHandler(cartService)

	gin.SetMode(gin.TestMode)
}

func aProductWithIDNamePriceAndStockIsAvailable(productID, productName, price, stock string) error {
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

func iAddOfProductToTheCart(quantity, productID string) error {
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

	cartRepo := db.New(config.DB)
	count, err := cartRepo.CountUniqueItemsInCart(context.Background())
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

	if count >= 10 {
		if recorder.Code != http.StatusBadRequest {
			return fmt.Errorf("expected status code 400 but got %d: %s", recorder.Code, recorder.Body.String())
		}
	} else if recorder.Code != http.StatusOK {
		return fmt.Errorf("expected status code 200 but got %d: %s", recorder.Code, recorder.Body.String())
	}

	return nil
}

func theTotalUniqueItemsInTheCartShouldBe(expectedUniqueItems string) error {
	expected, err := strconv.Atoi(expectedUniqueItems)

	if err == nil {

		var responseBody map[string]interface{}
		if err := json.Unmarshal(lastResponse.Body.Bytes(), &responseBody); err != nil {
			return err
		}

		uniqueItems, ok := responseBody["total unique items added"].(float64)
		if !ok {
			return fmt.Errorf("could not find 'total unique items added' in response")
		}

		if int(uniqueItems) != expected {
			return fmt.Errorf("expected %d unique items, but got %d", expected, int(uniqueItems))
		}

	} else {

		if lastResponse.Code != http.StatusBadRequest {
			return fmt.Errorf("expected status code 400 but got %d: %s", lastResponse.Code, lastResponse.Body.String())

		}

		if !(strings.Contains(lastResponse.Body.String(), "cannot add more than 10 unique items")) {
			return fmt.Errorf("expected error, %s got, %s", expectedUniqueItems, lastResponse.Body.String())
		}
	}

	return nil
}

func InitializeAddItemScenario(ctx *godog.ScenarioContext) {
	SetupScenario()

	ctx.Step(`^a product with ID "([^"]*)", name "([^"]*)", price "([^"]*)", and stock "([^"]*)" is available$`, aProductWithIDNamePriceAndStockIsAvailable)
	ctx.Step(`^I add "([^"]*)" of product "([^"]*)" to the cart$`, iAddOfProductToTheCart)
	ctx.Step(`^the total unique items in the cart should be "([^"]*)", and error should be returned when attempted to add more than 10 unique items$`, theTotalUniqueItemsInTheCartShouldBe)

}

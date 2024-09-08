package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"shopping-cart/config"
	"shopping-cart/data"
	"shopping-cart/db"
	"shopping-cart/handler"
	"shopping-cart/service"
	"strings"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
)

type CartSteps struct {
	cartHandler *handler.CartHandler
	cartRepo    *data.CartRepo
	recorder    *httptest.ResponseRecorder
	router      *gin.Engine
}

func (s *CartSteps) Setup() {
	s.router = gin.New()
	s.router.POST("/cart/add", s.cartHandler.AddItem)
	s.recorder = httptest.NewRecorder()
}

func (s *CartSteps) cleanup() {
	err := s.cartRepo.RemoveAllItem(context.Background())
	if err != nil {
		fmt.Printf("failed to clean up cart after test, %v", err)
	}
}

func (s *CartSteps) AddAnItemWithIDNamePriceAndQuantity(itemID int, itemName string, price float64, quantity int) error {

	req, _ := http.NewRequest(http.MethodPost, "/cart/add", strings.NewReader(fmt.Sprintf(`{"item_id":"%v","item_name":"%v","price":"%v","quantity":"%v"}`, itemID, itemName, price, quantity)))
	req.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(s.recorder, req)

	return nil

}

func (s *CartSteps) theItemShouldBeAddedSuccessfully() error {

	if s.recorder.Code != http.StatusOK {
		return fmt.Errorf("expected status 200 but got %d", s.recorder.Code)
	}

	want := "Item added to cart"
	got := s.recorder.Body.String()

	if !strings.Contains(got, want) {
		return fmt.Errorf("want: %s, got: %s", want, got)
	}

	return nil

}

func (s *CartSteps) theTotalNumberOfUniqueItemsShouldBe(uniqueItemCount int) error {
	if s.recorder.Code != http.StatusOK {
		return fmt.Errorf("expected status 200 but got %d", s.recorder.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(s.recorder.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	got, ok := response["total unique items added"].(float64)
	if !ok {
		return fmt.Errorf("response does not contain 'total unique items added' field or it's not a number")
	}

	if int(got) != uniqueItemCount {
		return fmt.Errorf("expected %d unique items but got %d", uniqueItemCount, int(got))
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	dbQueries := db.New(config.DB)
	cartRepo := data.NewCartRepo(dbQueries)
	cartService := service.NewCartService(cartRepo)
	cartHandler := handler.NewCartHandler(cartService)

	steps := &CartSteps{
		cartHandler: cartHandler,
		cartRepo:    cartRepo,
	}
	steps.cleanup()
	steps.Setup()
	ctx.Step(`^I add an item with ID (\d+), name "([^"]*)", price (\d+\.?\d*), and quantity (\d+)$`, steps.AddAnItemWithIDNamePriceAndQuantity)
	ctx.Step(`^the item should be added successfully$`, steps.theItemShouldBeAddedSuccessfully)
	ctx.Step(`^the total number of unique items should be (\d+)$`, steps.theTotalNumberOfUniqueItemsShouldBe)
}

package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/leandroberetta/stoqr/stoqr-api/mocks"
	"github.com/leandroberetta/stoqr/stoqr-api/models"
)

func TestCreateItemOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockItemRepository := mocks.NewMockItemRepository(ctrl)

	mockItemRepository.
		EXPECT().
		CreateItem(gomock.AssignableToTypeOf(&models.Item{})).
		DoAndReturn(func(item *models.Item) {
			createFakeItem(item)
		}).
		Return(nil)

	itemService := NewItemService(mockItemRepository)

	req, err := http.NewRequest("POST", "/api/items", bytes.NewReader(createFakeJSONItem()))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(itemService.CreateItem)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestCreateItemBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockItemRepository := mocks.NewMockItemRepository(ctrl)

	itemService := NewItemService(mockItemRepository)

	req, err := http.NewRequest("POST", "/api/items", bytes.NewReader([]byte{}))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(itemService.CreateItem)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestCreateItemInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockItemRepository := mocks.NewMockItemRepository(ctrl)

	mockItemRepository.
		EXPECT().
		CreateItem(gomock.AssignableToTypeOf(&models.Item{})).
		Return(errors.New("error"))

	itemService := NewItemService(mockItemRepository)

	req, err := http.NewRequest("POST", "/api/items", bytes.NewReader(createFakeJSONItem()))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(itemService.CreateItem)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestWithdrawItemOK(t *testing.T) {
	cases := []struct {
		name string
		item models.Item
	}{
		{name: "actual1", item: models.Item{ID: 1, Name: "Test", Desired: 1, Actual: 1}},
		{name: "actual0", item: models.Item{ID: 1, Name: "Test", Desired: 1, Actual: 0}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockItemRepository := mocks.NewMockItemRepository(ctrl)

			mockItemRepository.
				EXPECT().
				ReadItem(gomock.Any()).
				Return(c.item, nil)

			mockItemRepository.
				EXPECT().
				UpdateItem(gomock.Any(), gomock.Any()).
				Return(nil)

			itemService := NewItemService(mockItemRepository)

			req, err := http.NewRequest("GET", "/api/items/withdraw/1", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/api/items/withdraw/{itemId}", itemService.WithdrawItem)
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			updatedItem := &models.Item{}
			err = json.Unmarshal(rr.Body.Bytes(), updatedItem)

			if updatedItem.Actual != 0 {
				t.Errorf("wrong actual value: got %v want %v", updatedItem.Actual, 0)
			}
		})
	}
}

func TestWithdrawItemNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockItemRepository := mocks.NewMockItemRepository(ctrl)

	mockItemRepository.
		EXPECT().
		ReadItem(gomock.Any()).
		Return(models.Item{}, errors.New("error"))

	itemService := NewItemService(mockItemRepository)

	req, err := http.NewRequest("GET", "/api/items/withdraw/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/items/withdraw/{itemId}", itemService.WithdrawItem)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestWithdrawItemBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockItemRepository := mocks.NewMockItemRepository(ctrl)
	itemService := NewItemService(mockItemRepository)

	req, err := http.NewRequest("GET", "/api/items/withdraw/wrong", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/items/withdraw/{itemId}", itemService.WithdrawItem)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestWithdrawItemInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockItemRepository := mocks.NewMockItemRepository(ctrl)
	item := &models.Item{}
	createFakeItem(item)

	mockItemRepository.
		EXPECT().
		ReadItem(gomock.Any()).
		Return(*item, nil)

	mockItemRepository.
		EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		Return(errors.New("error"))

	itemService := NewItemService(mockItemRepository)

	req, err := http.NewRequest("GET", "/api/items/withdraw/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/items/withdraw/{itemId}", itemService.WithdrawItem)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func createFakeItem(item *models.Item) {
	item.ID = 1
	item.Name = "Test"
	item.Desired = 1
	item.Actual = 1
}

func createFakeJSONItem() []byte {
	item := &models.Item{}
	createFakeItem(item)
	bytes, _ := json.Marshal(item)
	return bytes
}

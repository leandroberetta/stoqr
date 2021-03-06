package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/leandroberetta/stoqr/stoqr-api/mocks"
	"github.com/leandroberetta/stoqr/stoqr-api/models"
)

func TestCreateItem(t *testing.T) {
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

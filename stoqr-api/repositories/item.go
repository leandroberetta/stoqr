package repositories

import (
	"fmt"

	"github.com/leandroberetta/stoqr/stoqr-api/models"
	"gorm.io/gorm"
)

// ItemRepository interface define the methods to persist items
type ItemRepository interface {
	CreateItem(item *models.Item) error
	ReadItem(id int) (models.Item, error)
	UpdateItem(id int, item models.Item) error
	DeleteItem(id int) error
	ReadItems(filter string) ([]models.Item, error)
}

// ItemRepositorySQL persist items into a SQL database
type ItemRepositorySQL struct {
	*gorm.DB
}

// CreateItem persists an item into a database
func (db *ItemRepositorySQL) CreateItem(item *models.Item) error {
	result := db.Create(item)
	return result.Error
}

// ReadItem gets an item from a database
func (db *ItemRepositorySQL) ReadItem(id int) (models.Item, error) {
	var item models.Item
	result := db.First(&item, id)
	return item, result.Error
}

// UpdateItem updates an item and persists it into a database
func (db *ItemRepositorySQL) UpdateItem(id int, updatedItem models.Item) error {
	var item models.Item
	result := db.Find(&item, id)
	if result.Error != nil {
		return result.Error
	}
	item.Name = updatedItem.Name
	item.Desired = updatedItem.Desired
	item.Actual = updatedItem.Actual
	db.Save(&item)
	return nil
}

// DeleteItem removes an item from a database
func (db *ItemRepositorySQL) DeleteItem(id int) error {
	result := db.Delete(&models.Item{}, id)
	return result.Error
}

// ReadItems gets items from a database optionally filtering by name
func (db *ItemRepositorySQL) ReadItems(filter string) ([]models.Item, error) {
	var items []models.Item
	if filter != "" {
		db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", filter)).Find(&items)
	} else {
		db.Find(&items)
	}
	return items, nil
}

// NewItemRepositorySQL returns a new ItemRepositorySQL instance
func NewItemRepositorySQL(db *gorm.DB) ItemRepository {
	return &ItemRepositorySQL{db}
}

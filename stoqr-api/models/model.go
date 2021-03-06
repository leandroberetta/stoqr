package models

// Item is the model of the item object
type Item struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Desired int    `json:"desired"`
	Actual  int    `json:"actual"`
}

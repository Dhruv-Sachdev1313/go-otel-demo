package models

// ShoppingCart represents a user's shopping cart
type ShoppingCart struct {
	UserID string   `json:"user_id"`
	Items  []string `json:"items"`
}

package handler

import "time"

type Location struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type,omitempty"`
	OperatingHours string `json:"operating_hours,omitempty"`
	ImageURL       string `json:"image_url,omitempty"`
}

type Tenant struct {
	ID          int64  `json:"id"`
	LocationID  int64  `json:"location_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
	IsOpen      bool   `json:"is_open"`
}

type Menu struct {
	ID          int64   `json:"id"`
	TenantID    int64   `json:"tenant_id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	Category    string  `json:"category,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
	IsAvailable bool    `json:"is_available"`
}

type User struct {
	ID           int64     `json:"id,omitempty"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number,omitempty"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type Order struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	TenantID        int64     `json:"tenant_id"`
	TotalPrice      float64   `json:"total_price"`
	OrderType       string    `json:"order_type"` 
	Status          string    `json:"status"`
	DeliveryAddress string    `json:"delivery_address,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	// Untuk join data
	OrderItems []OrderItem `json:"order_items,omitempty"`
	TenantInfo Tenant      `json:"tenant_info,omitempty"`
}

type OrderItem struct {
	ID           int64   `json:"id"`
	OrderID      int64   `json:"order_id"`
	MenuID       int64   `json:"menu_id"`
	Quantity     int     `json:"quantity"`
	PricePerItem float64 `json:"price_per_item"`
	// Untuk join data
	MenuInfo Menu `json:"menu_info,omitempty"`
}
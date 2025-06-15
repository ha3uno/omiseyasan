package models

import "time"

// HistoryEntry represents a single history record
type HistoryEntry struct {
	ID           int     `json:"id"`
	Timestamp    string  `json:"timestamp"`
	Description  string  `json:"description"`
	EffortHours  float64 `json:"effortHours"`
	ClaudePrompt string  `json:"claudePrompt"`
}

// User represents a user account
type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        int     `json:"id" db:"id"`
	OrderID   int     `json:"orderId" db:"order_id"`
	ProductID int     `json:"productId" db:"product_id"`
	Name      string  `json:"name" db:"product_name"`
	Price     float64 `json:"price" db:"price"`
	Quantity  int     `json:"quantity" db:"quantity"`
	Subtotal  float64 `json:"subtotal"` // 計算値なのでdbタグなし
}

// ShippingInfo represents shipping information
type ShippingInfo struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
}

// Order represents a customer order
type Order struct {
	ID           int          `json:"id" db:"id"`
	UserID       *int         `json:"userId,omitempty" db:"user_id"` // Nullable for guest orders
	Timestamp    string       `json:"timestamp"`                     // フォーマット済みの文字列
	OrderedAt    time.Time    `json:"-" db:"ordered_at"`             // データベース用のタイムスタンプ
	Items        []OrderItem  `json:"items"`                         // JOINで取得
	TotalAmount  float64      `json:"totalAmount" db:"total_amount"`
	ShippingInfo ShippingInfo `json:"shippingInfo"` // 個別のフィールドからマッピング
}
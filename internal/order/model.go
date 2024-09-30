package order

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusCancelled  Status = "cancelled"
)

func NewOrderStatus(status string) Status {
	return Status(status)
}

type Order struct {
	ID             uuid.UUID `json:"id"`
	CustomerID     uuid.UUID `json:"customer_id"`
	ProductSKU     string    `json:"product_sku"`
	DeliveryOption string    `json:"delivery_option"`
	Quantity       int       `json:"quantity"`
	Price          int       `json:"price"`
	Status         Status    `json:"status"`
	CreationTime   time.Time `json:"creation_time"`
}

func (o *Order) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func FromJSON(raw []byte) (Order, error) {
	var o Order
	if err := json.Unmarshal(raw, &o); err != nil {
		return Order{}, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return o, nil
}

package order

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

var sku = []string{"sku1", "sku2", "sku3", "sku4", "sku5"}
var deliveryOptions = []string{"standard", "express", "next-day"}
var status = []Status{StatusPending, StatusProcessing, StatusCompleted, StatusCancelled}

func Random() Order {

	rand.Seed(uint64(time.Now().UnixNano()))

	randSKU := sku[rand.Intn(len(sku))]
	randDeliveryOption := deliveryOptions[rand.Intn(len(deliveryOptions))]
	randStatus := status[rand.Intn(len(status))]

	return Order{
		ID:             uuid.New(),
		CustomerID:     uuid.New(),
		ProductSKU:     randSKU,
		DeliveryOption: randDeliveryOption,
		Quantity:       rand.Intn(10),
		Price:          rand.Intn(999999),
		Status:         randStatus,
		CreationTime:   generateRandomTime(),
	}
}

func generateRandomTime() time.Time {
	start := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC)

	delta := end.Sub(start)

	randDuration := time.Duration(rand.Int63n(int64(delta)))

	return start.Add(randDuration)
}

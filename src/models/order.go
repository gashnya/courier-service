package models

type OrderDto struct {
	OrderId       int64    `json:"order_id"`
	Weight        float32  `json:"weight"`
	Regions       int32    `json:"regions"`
	DeliveryHours []string `json:"delivery_hours"`
	Cost          int32    `json:"cost"`
	CompletedTime *string  `json:"completed_time,omitempty"`
}

type CreateOrderDto struct {
	Cost          int32    `json:"cost" validate:"gte=1"`
	DeliveryHours []string `json:"delivery_hours" validate:"gt=0,dive,isValidHours"`
	Regions       int32    `json:"regions" validate:"gte=1"`
	Weight        float32  `json:"weight" validate:"gt=0"`
}

type CompleteOrderDto struct {
	CourierId    int64  `json:"courier_id" validate:"gte=1"`
	OrderId      int64  `json:"order_id" validate:"gte=1"`
	CompleteTime string `json:"complete_time" validate:"isValidTimestamp"`
}

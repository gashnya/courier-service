package models

type CourierType string

const (
	Auto CourierType = "AUTO"
	Bike             = "BIKE"
	Foot             = "FOOT"
)

type CourierDto struct {
	CourierId    int64       `json:"courier_id"`
	CourierType  CourierType `json:"courier_type"`
	Regions      []int32     `json:"regions"`
	WorkingHours []string    `json:"working_hours"`
}

type CreateCourierDto struct {
	CourierType  CourierType `json:"courier_type" validate:"isValidCourierType"`
	Regions      []int32     `json:"regions" validate:"gt=0,dive,gte=1"`
	WorkingHours []string    `json:"working_hours" validate:"gt=0,dive,isValidHours"`
}

type CourierMetaInfoDto struct {
	CourierId    int64       `json:"courier_id"`
	CourierType  CourierType `json:"courier_type"`
	WorkingHours []string    `json:"working_hours"`
	Regions      []int32     `json:"regions"`
	Rating       *int32      `json:"rating,omitempty"`
	Earnings     *int32      `json:"earnings,omitempty"`
}

func GetCourierRatingCoefficient(courier *CourierDto) int32 {
	switch courier.CourierType {
	case Auto:
		return 1
	case Bike:
		return 2
	case Foot:
		return 3
	default:
		return 0
	}
}

func GetCourierEarningsCoefficient(courier *CourierDto) int32 {
	switch courier.CourierType {
	case Auto:
		return 4
	case Bike:
		return 3
	case Foot:
		return 2
	default:
		return 0
	}
}

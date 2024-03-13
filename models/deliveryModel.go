package models

type Delivery struct {
	DeliveryID int
	DeliveryVolume int
	PricePerLiterNetto float32
	PricePerLiterBrutto float32
	TotalPriceNetto float32
	TotalPriceBrutto float32
	DeliveryDate string
}

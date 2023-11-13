package orderType

type ORDER_TYPE string

const (
	PickAndDelivery  ORDER_TYPE = "PICK_AND_DELIVERY"
	PickDeliveryPark ORDER_TYPE = "PICK_DELIVERY_PARK"
	PickOnly         ORDER_TYPE = "PICK_ONLY"
	DeliveryOnly     ORDER_TYPE = "DELIVERY_ONLY"
	DeliveryPark     ORDER_TYPE = "DELIVERY_PARK"
)

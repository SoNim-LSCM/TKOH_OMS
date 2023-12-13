package processingStatus

type PROCESSING_STATUS string

const (
	planningToStartLocation       PROCESSING_STATUS = "PLANNING_TO_START_LOCATION"
	goingToStartLocation          PROCESSING_STATUS = "GOING_TO_START_LOCATION"
	queuingAtStartBay             PROCESSING_STATUS = "QUEUING_AT_START_BAY"
	arrivedStartLocation          PROCESSING_STATUS = "ARRIVED_START_LOCATION"
	planningToDestinationLocation PROCESSING_STATUS = "PLANNING_TO_DESTINATION_LOCATION"
	goingToDestinationLocation    PROCESSING_STATUS = "GOING_TO_DESTINATION_LOCATION"
	queuingAtDestinationBay       PROCESSING_STATUS = "QUEUING_AT_DESTINATION_BAY"
	arrivedDestinationLocation    PROCESSING_STATUS = "ARRIVED_DESTINATION_LOCATION"
	planningToParking             PROCESSING_STATUS = "PLANNING_TO_PARKING"
	goingToParking                PROCESSING_STATUS = "GOING_TO_PARKING"
	parking                       PROCESSING_STATUS = "PARKING"
)

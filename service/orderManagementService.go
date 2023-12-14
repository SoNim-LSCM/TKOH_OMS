package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	apiHandler "github.com/SoNim-LSCM/TKOH_OMS/api"
	"github.com/SoNim-LSCM/TKOH_OMS/constants/orderStatus"
	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"
	"github.com/SoNim-LSCM/TKOH_OMS/models/orderManagement"
	"github.com/SoNim-LSCM/TKOH_OMS/models/rfms"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"
	"gorm.io/gorm"
)

func FindOrders(filterFields string, filterValues ...interface{}) ([]db_models.Orders, error) {
	var orders []db_models.Orders
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		query := "SELECT *, C.location_name as start_location_name, D.location_name as end_location_name FROM tkoh_oms.orders LEFT JOIN tkoh_oms.locations C ON tkoh_oms.orders.start_location_id = C.location_id  LEFT JOIN tkoh_oms.locations D ON tkoh_oms.orders.end_location_id = D.location_id WHERE " + filterFields + " ORDER BY case when processing_status like 'UNLOADING' then 1 when processing_status like 'ARRIVED%' then 2 when processing_status like 'QUEUING%' then 3 when processing_status like 'MOVING_TO_LAYBY_AREA' then 4 when processing_status like 'GOING%' then 5 when processing_status like 'PLANNING%' then 6 else 7 end asc"
		if err := FindRecordsWithRaw(tx, &orders, query, filterValues...); err != nil {
			return errors.New("Failed to search: " + err.Error())
		}
		return nil
	})
	return orders, err
}

func FindRoutines(filterFields string, filterValues ...interface{}) ([]db_models.Routines, error) {
	var routines []db_models.Routines
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		query := "SELECT *, C.location_name as start_location_name, D.location_name as end_location_name FROM tkoh_oms.routines LEFT JOIN tkoh_oms.locations C ON tkoh_oms.routines.start_location_id = C.location_id  LEFT JOIN tkoh_oms.locations D ON tkoh_oms.routines.end_location_id = D.location_id WHERE " + filterFields
		if err := FindRecordsWithRaw(tx, &routines, query, filterValues...); err != nil {
			return errors.New("Failed to search: " + err.Error())
		}
		return nil
	})
	return routines, err
}

func AddOrders(orderRequest dto.AddDeliveryOrderDTO, userId int) (orderManagement.OrderList, error) {
	var orderList orderManagement.OrderList
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var lastSchedule db_models.Schedules

		// create new schedule
		schedulesList := []db_models.Schedules{{ScheduleID: 0, ScheduleStatus: "CREATED", ScheduleCraeteTime: utils.GetTimeNowString(), OrderType: orderRequest.OrderType, NumberOfAmrRequire: orderRequest.NumberOfAmrRequire, LastUpdateTime: utils.GetTimeNowString()}}
		if err := AddRecords(tx, schedulesList); err != nil {
			return err
		}
		// check no of schedules
		if err := tx.Table("schedules").Last(&lastSchedule).Error; err != nil {
			return err
		}
		// translate order request to orders
		orders, err := OrderRequestToOrders(orderRequest, lastSchedule.ScheduleID, userId)
		if err != nil {
			return err
		}
		// create new orders
		if err := AddRecords(tx, orders); err != nil {
			return err
		}
		// translate new orders to order response
		orderList, err = OrderListToOrderResponse(orders)
		if err != nil {
			return err
		}
		return nil
	})
	return orderList, err
}

func AddRoutines(routineRequest dto.AddRoutineDTO, userId int) (orderManagement.RoutineOrderList, error) {
	database.CheckDatabaseConnection()
	var routineOrderList orderManagement.RoutineOrderList
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// translate order request to routines
		routines, err := RoutineRequestToRoutines(routineRequest, userId)
		if err != nil {
			return err
		}
		// create new orders
		if err := AddRecords(tx, routines); err != nil {
			return err
		}
		// translate new routine to routine response
		routineOrderList, err = RoutineListToRoutineResponse(routines)
		if err != nil {
			return err
		}
		return nil
	})

	return routineOrderList, err
}

func TriggerOrderOrderIds(orderId []int) (orderManagement.OrderList, error) {
	var orders []db_models.Orders
	var orderList orderManagement.OrderList
	updateFields := []string{"processing_status"}
	updateMap1 := utils.CreateMap(updateFields, string("ARRIVED_START_LOCATION"))
	updateMap2 := utils.CreateMap(updateFields, string("ARRIVED_END_LOCATION"))
	// updateValues := []string{string(orderStatus.Processing), timeNow}

	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := UpdateRecords(tx, &orders, "orders", updateMap1, "order_id IN ? AND order_status = ? AND processing_status IN ?", orderId, orderStatus.Processing, "QUEUING_AT_START_BAY")
		// if err != nil {
		// 	return err
		// }
		err = UpdateRecords(tx, &orders, "orders", updateMap2, "order_id IN ? AND order_status = ? AND processing_status IN ?", orderId, orderStatus.Processing, "QUEUING_AT_END_BAY")
		// if err != nil {
		// 	return err
		// }
		orderList, err = OrderListToOrderResponse(orders)
		if err != nil {
			return err
		}
		return nil
	})
	return orderList, err
}

func TriggerOrderScheduleId(scheduleId int) (orderManagement.OrderList, error) {
	var orders []db_models.Orders
	var orderList orderManagement.OrderList
	updateFields := []string{"order_status", "order_start_time"}
	timeNow := utils.GetTimeNowString()
	updateMap := utils.CreateMap(updateFields, string(orderStatus.Processing), timeNow)

	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := UpdateRecords(tx, &orders, "orders", updateMap, "schedule_id = ?", scheduleId)
		if err != nil {
			return err
		}
		orderList, err = OrderListToOrderResponse(orders)
		if err != nil {
			return err
		}
		return nil
	})
	return orderList, err
}

func UpdateOrders(userId int, request dto.UpdateDeliveryOrderDTO) (orderManagement.OrderList, error) {
	var orderList orderManagement.OrderList

	schedules := []db_models.Schedules{}
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if FindRecords(tx, &schedules, "schedules", "schedule_id = ?", request.ScheduleID) != nil {
			return errors.New("Failed to find schedule with schedule id")
		}

		orders := []db_models.Orders{}
		if FindRecords(tx, &orders, "orders", "schedule_id = ? AND order_status = ?", request.ScheduleID, orderStatus.Created) != nil {
			return errors.New("Failed to find order with schedule id")
		}

		if len(orders) != schedules[0].NumberOfAmrRequire {
			return errors.New("Update fail, some orders already started")
		}
		expectedStartTime, err := StringToDatetime(request.ExpectedStartTime)
		if err != nil {
			return errors.New("Fail translate expectedStartTime")
		}
		expectedDeliveryTime, err := StringToDatetime(request.ExpectedDeliveryTime)
		if err != nil {
			return errors.New("Fail translate expectedDeliveryTime")
		}

		updateMap := utils.CreateMap([]string{"number_of_amr_require"}, request.NumberOfAmrRequire)
		err = AddSchedulesLogs(tx, userId, "schedule_id = ?", request.ScheduleID)
		if err != nil {
			return errors.New("Failed to create log")
		}
		err = UpdateRecords(tx, &[]db_models.Schedules{}, "schedules", updateMap, "schedule_id = ?", request.ScheduleID)
		if err != nil {
			return errors.New("Failed to update schedule table")
		}
		for i := 0; i < utils.Max(request.NumberOfAmrRequire, schedules[0].NumberOfAmrRequire); i++ {
			// cancel orders
			if i < (schedules[0].NumberOfAmrRequire - request.NumberOfAmrRequire) {
				updatedOrderList := []db_models.Orders{}
				updateMap := utils.CreateMap([]string{"order_status"}, orderStatus.Canceled)
				if UpdateRecords(tx, &updatedOrderList, "orders", updateMap, "order_id = ?", orders[i].OrderID) != nil {
					return errors.New("Failed to translate orders to order response")
				}
				updatedOrderResponse, err := OrderListToOrderResponse(updatedOrderList)
				if err != nil {
					return errors.New("Failed to translate orders to order response")
				}
				orderList = append(orderList, updatedOrderResponse[0])
				// change orders
			} else if i < schedules[0].NumberOfAmrRequire {
				updatedOrderList := []db_models.Orders{}
				updateMap := utils.CreateMap([]string{"schedule_id", "start_location_id", "end_location_id", "expected_start_time", "expected_delivery_time"}, request.ScheduleID, request.StartLocationID, request.EndLocationID, expectedStartTime, expectedDeliveryTime)
				err = AddOrdersLogs(tx, userId, "order_id = ?", orders[i].OrderID)
				if err != nil {
					return errors.New("Failed to create log")
				}
				err = UpdateRecords(tx, &updatedOrderList, "orders", updateMap, "order_id = ?", orders[i].OrderID)
				if err != nil {
					return errors.New("Failed to update order")
				}
				updatedOrderResponse, err := OrderListToOrderResponse(updatedOrderList)
				if err != nil {
					return errors.New("Failed to translate orders to order response")
				}
				orderList = append(orderList, updatedOrderResponse[0])
				// add orders
			} else {
				bJson, err := json.Marshal(request)
				var orderRequest dto.AddDeliveryOrderDTO
				json.Unmarshal(bJson, &orderRequest)
				orderRequest.NumberOfAmrRequire = 1
				// translate order request to uploadOrders
				uploadOrders, err := OrderRequestToOrders(orderRequest, request.ScheduleID, 6)
				if err != nil {
					return err
				}
				uploadOrders[0].OrderType = schedules[0].OrderType
				// create new orders
				if err := AddRecords(tx, uploadOrders); err != nil {
					return err
				}
				updatedOrderResponse, err := OrderListToOrderResponse(uploadOrders)
				if err != nil {
					return errors.New("Failed to translate orders to order response")
				}
				orderList = append(orderList, updatedOrderResponse[0])
			}
		}

		return nil
	})
	return orderList, err
}

func CancelOrders(scheduleId int) (orderManagement.OrderList, error) {
	var orderList orderManagement.OrderList
	var schedules []db_models.Schedules
	var orders []db_models.Orders
	var updatedOrders []db_models.Orders

	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := FindRecords(tx, &schedules, "schedules", "schedule_id = ?", scheduleId); err != nil {
			return err
		}

		if err := FindRecords(tx, &orders, "orders", "schedule_id = ?", scheduleId); err != nil {
			return err
		}

		if len(orders) == 0 || len(schedules) == 0 {
			return errors.New("Orders not found")
		}

		var amrs = schedules[0].NumberOfAmrRequire

		for _, order := range orders {
			if (order.OrderStatus != string(orderStatus.Created)) && (order.OrderStatus != string(orderStatus.Canceled)) {
				return errors.New("Cancel failed, order started")
			} else if order.OrderStatus == string(orderStatus.Created) {
				amrs -= 1
			}
		}

		if amrs != 0 {
			return errors.New("Cancel failed, amr number not match")
		}

		updateMap := utils.CreateMap([]string{"schedule_status", "number_of_amr_require"}, orderStatus.Canceled, 0)
		err := UpdateRecords(tx, &[]db_models.Schedules{}, "schedules", updateMap, "schedule_id = ?", scheduleId)
		if err != nil {
			return err
		}

		updateMap = utils.CreateMap([]string{"order_status"}, orderStatus.Canceled)
		err = UpdateRecords(tx, &updatedOrders, "orders", updateMap, "schedule_id = ?", scheduleId)
		if err != nil {
			return err
		}
		orderList, err = OrderListToOrderResponse(updatedOrders)
		if err != nil {
			return err
		}
		return nil
	})
	return orderList, err
}

func UpdateRoutineOrders(userId int, request dto.UpdateRoutineDeliveryOrderDTO) (orderManagement.RoutineOrderList, error) {
	var updatedList orderManagement.RoutineOrderList
	var routinesList = []db_models.Routines{}
	expectedStartTime, err := RoutineResponseTimeToString(request.ExpectedStartTime)
	if err != nil {
		return updatedList, err
	}
	expectedDeliveryTime, err := RoutineResponseTimeToString(request.ExpectedDeliveryTime)
	if err != nil {
		return updatedList, err
	}
	routinePattern, err := RoutinePatternToString(request.RoutinePattern)
	if err != nil {
		return updatedList, err
	}
	updateMap := utils.CreateMap([]string{"routine_pattern", "number_of_amr_require", "start_location_id", "end_location_id", "expected_start_time", "expected_delivery_time"}, routinePattern, request.NumberOfAmrRequire, request.StartLocationID, request.EndLocationID, expectedStartTime, expectedDeliveryTime)
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		err := UpdateRecords(tx, &routinesList, "routines", updateMap, "routine_id = ?", request.RoutineID)
		if err != nil {
			return err
		}
		err = AddRoutinesLogs(tx, userId, "routine_id = ?", request.RoutineID)
		if err != nil {
			return err
		}
		updatedList, err = RoutineListToRoutineResponse(routinesList)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return updatedList, err
	}

	return updatedList, nil
}
func OrderRequestToOrders(orderRequest dto.AddDeliveryOrderDTO, scheduleNo int, userId int) ([]db_models.Orders, error) {
	var orders []db_models.Orders
	for i := 0; i < orderRequest.NumberOfAmrRequire; i++ {
		var err error
		var order db_models.Orders
		order.ScheduleID = scheduleNo
		// order.OrderID = i + orderNo
		order.OrderID = 0
		order.OrderType = orderRequest.OrderType
		order.OrderCreatedType = "ADHOC"
		order.OrderCreatedBy = userId
		order.OrderStatus = "TO_BE_CREATED"
		order.OrderStartTime, err = StringToDatetime("")
		if err != nil {
			return orders, err
		}
		order.ActualArrivalTime, err = StringToDatetime("")
		if err != nil {
			return orders, err
		}
		order.StartLocationID = orderRequest.StartLocationID
		order.StartLocationName = orderRequest.StartLocationName
		order.EndLocationID = orderRequest.EndLocationID
		order.EndLocationName = orderRequest.EndLocationName
		order.ExpectedStartTime, err = StringToDatetime(orderRequest.ExpectedStartTime)
		if err != nil {
			return orders, err
		}
		order.ExpectedDeliveryTime, err = StringToDatetime(orderRequest.ExpectedDeliveryTime)
		if err != nil {
			return orders, err
		}
		order.ExpectedArrivalTime, err = StringToDatetime("")
		if err != nil {
			return orders, err
		}
		order.ProcessingStatus = ""
		order.LastUpdateTime = utils.GetTimeNowString()
		order.LastUpdateBy = userId
		orders = append(orders, order)
	}
	return orders, nil
}

func OrderListToOrderResponse(orderList []db_models.Orders) (orderManagement.OrderList, error) {
	var orderListResponse orderManagement.OrderList
	// roomList, err := FindAllDutyRooms()
	// if err != nil {
	// 	return orderListResponse, err
	// }
	jsonString, err := json.Marshal(orderList)
	if err != nil {
		return orderListResponse, err
	}
	json.Unmarshal(jsonString, &orderListResponse)
	for i, order := range orderList {
		var err error
		orderListResponse[i].StartTime, err = StringToResponseTime(order.OrderStartTime)
		if err != nil {
			return orderListResponse, err
		}
		orderListResponse[i].ActualArrivalTime, err = StringToResponseTime(order.ActualArrivalTime)
		if err != nil {
			return orderListResponse, err
		}
		orderListResponse[i].ExpectedStartTime, err = StringToResponseTime(order.ExpectedStartTime)
		if err != nil {
			return orderListResponse, err
		}
		orderListResponse[i].ExpectedDeliveryTime, err = StringToResponseTime(order.ExpectedDeliveryTime)
		if err != nil {
			return orderListResponse, err
		}
		orderListResponse[i].ExpectedArrivalTime, err = StringToResponseTime(order.ExpectedArrivalTime)
		if err != nil {
			return orderListResponse, err
		}
		// orderListResponse[i].StartLocationName = roomList[orderListResponse[i].StartLocationID-1].LocationName
		// orderListResponse[i].EndLocationName = roomList[orderListResponse[i].EndLocationID-1].LocationName
	}
	return orderListResponse, nil
}

func RoutineRequestToRoutines(routinesRequest dto.AddRoutineDTO, userId int) ([]db_models.Routines, error) {
	var routinesList []db_models.Routines
	var routines db_models.Routines
	bJson, err := json.Marshal(routinesRequest)
	if err != nil {
		return routinesList, err
	}
	json.Unmarshal(bJson, &routines)
	if err != nil {
		return routinesList, err
	}
	routines.LastUpdateBy = userId
	routines.RoutineCreatedBy = userId
	routines.LastUpdateTime = utils.GetTimeNowString()
	routines.RoutinePattern, err = RoutinePatternToString(routinesRequest.RoutinePattern)
	if err != nil {
		return routinesList, nil
	}
	routines.ExpectedDeliveryTime, err = StringToDatetime(routinesRequest.ExpectedDeliveryTime)
	if err != nil {
		return routinesList, err
	}
	routines.ExpectedStartTime, err = StringToDatetime(routinesRequest.ExpectedStartTime)
	if err != nil {
		return routinesList, err
	}
	routinesList = append(routinesList, routines)
	return routinesList, nil
}

func RoutineListToRoutineResponse(routineList []db_models.Routines) (orderManagement.RoutineOrderList, error) {
	var routineOrderListResponse orderManagement.RoutineOrderList
	jsonString, err := json.Marshal(routineList)
	if err != nil {
		return routineOrderListResponse, err
	}
	json.Unmarshal(jsonString, &routineOrderListResponse)
	if err != nil {
		return routineOrderListResponse, err
	}
	for i, routine := range routineList {
		var err error
		routinePattern, err := StringToRoutinePattern(routine.RoutinePattern)
		if err != nil {
			return routineOrderListResponse, err
		}
		routineOrderListResponse[i].RoutinePattern = routinePattern
		routineOrderListResponse[i].NextDeliveryDate, err = GetNextDeliveryDate(routinePattern)
		if err != nil {
			return routineOrderListResponse, err
		}
		routineOrderListResponse[i].ExpectedDeliveryTime, err = StringToRoutineResponseTime(routine.ExpectedDeliveryTime)
		if err != nil {
			return routineOrderListResponse, err
		}
		routineOrderListResponse[i].ExpectedStartTime, err = StringToRoutineResponseTime(routine.ExpectedStartTime)
		if err != nil {
			return routineOrderListResponse, err
		}
	}
	return routineOrderListResponse, nil
}

func BackgroundInitOrderToRFMS() error {

	if database.CheckDatabaseConnection() {
		err := database.DB.Transaction(func(tx *gorm.DB) error {

			orders := []db_models.Orders{}

			if err := FindRecords(tx, &orders, "orders", "order_status = ?", "TO_BE_CREATED"); err != nil {
				return err
			}

			if len(orders) == 0 {
				return nil
			}

			for _, order := range orders {

				jobNatures := []string{}
				locations := []int{}
				statusLocation := []string{}
				// orderTypes := strings.Split(order.OrderType, "_")
				// jobNatures = append(jobNatures, orderTypes[0])
				switch order.OrderType {
				case "PICK_AND_DELIVERY":
					jobNatures = append(jobNatures, "PICK")
					locations = append(locations, 1)
					statusLocation = append(statusLocation, "START_LOCATION")

					jobNatures = append(jobNatures, "DELIVERY")
					locations = append(locations, 1)
					statusLocation = append(statusLocation, "END_LOCATION")
				case "PICK_DELIVERY_PARK":
					jobNatures = append(jobNatures, "PICK")
					locations = append(locations, 1)
					statusLocation = append(statusLocation, "START_LOCATION")

					jobNatures = append(jobNatures, "DELIVERY")
					locations = append(locations, 1)
					statusLocation = append(statusLocation, "END_LOCATION")

					jobNatures = append(jobNatures, "PARK")
					locations = append(locations, 1)
					statusLocation = append(statusLocation, "PARKING")
				case "PICK_ONLY":
					jobNatures = append(jobNatures, "PICK")
					locations = append(locations, 1)
					statusLocation = append(statusLocation, "END_LOCATION")
				case "DELIVERY_ONLY":
					jobNatures = append(jobNatures, "DELIVERY")
					locations = append(locations, 1)
					statusLocation = append(statusLocation, "END_LOCATION")
				case "DELIVERY_PARK":
					jobNatures = append(jobNatures, "DELIVERY")
					locations = append(locations, 1)
					statusLocation = append(statusLocation, "END_LOCATION")

					jobNatures = append(jobNatures, "PARK")
					locations = append(locations, 1)
					statusLocation = append(statusLocation, "PARKING")
				case "PARK_ONLY":
					jobNatures = append(jobNatures, "PARK")
					locations = append(locations, 1)
					statusLocation = append(statusLocation, "PARKING")
				default:
					return errors.New("Unknown Order Type")
				}

				param := rfms.CreateJobRequest{JobNature: jobNatures[0], LocationID: 15}
				response := apiHandler.POST("/createJob", param)

				updateJobStatus := dto.ReportJobStatusResponseDTO{}
				err := json.Unmarshal(response, &updateJobStatus)
				if err != nil {
					return err
				}
				if updateJobStatus.ResponseMessage == "FAILED" {
					return errors.New("Create Job Failed with Reason: " + updateJobStatus.FailReason)
				}
				est, err := StringToDatetime(updateJobStatus.Body.Est)
				if err != nil {
					return err
				}
				eta, err := StringToDatetime(updateJobStatus.Body.Eta)
				if err != nil {
					return err
				}
				lastUpdateTime, err := StringToDatetime(updateJobStatus.Body.MessageTime)
				if err != nil {
					return err
				}
				newJob := db_models.Jobs{OrderID: order.OrderID, JobID: updateJobStatus.Body.JobID, JobType: jobNatures[0], JobStatus: updateJobStatus.Body.Status, ProcessingStatus: updateJobStatus.Body.ProcessingStatus, JobStartTime: est, ExpectedArrivalTime: eta, EndLocationID: updateJobStatus.Body.LocationId, FailedReason: updateJobStatus.FailReason, LastUpdateTime: lastUpdateTime}
				newJobs := []db_models.Jobs{newJob}
				// log.Print(newJob)
				err = AddRecords(tx, newJobs)
				if err != nil {
					return err
				}
				for i, jobNature := range jobNatures {
					if i > 0 {
						defaultTime, err := StringToDatetime("")
						if err != nil {
							return err
						}
						newJob := db_models.Jobs{OrderID: order.OrderID, JobID: 0, JobType: jobNature, JobStatus: "TO_BE_CREATED", ProcessingStatus: "UNKNOWN", JobStartTime: defaultTime, ExpectedArrivalTime: defaultTime, EndLocationID: locations[i], FailedReason: updateJobStatus.FailReason, LastUpdateTime: defaultTime, StatusLocation: statusLocation[0]}
						newJobs := []db_models.Jobs{newJob}
						err = AddRecords(tx, newJobs)
						if err != nil {
							return err
						}
					}
				}

				updatedList := []db_models.Orders{}
				updateFields := []string{"order_status", "job_id"}
				updateMap := utils.CreateMap(updateFields, "CREATED", updateJobStatus.Body.JobID)
				err = UpdateRecords(tx, &updatedList, "orders", updateMap, "order_id = ?", order.OrderID)
				if err != nil {
					return err
				}

				// _, updateMap, err := GetUpdateJobFields(response)
				// if err != nil {
				// 	return errors.New("Failed to phrase create job response (2)")
				// }
				// var updatedList = []db_models.Orders{}
				// err = UpdateRecords(tx, &updatedList, "orders", updateMap, "order_id = ?", order.OrderID)
				// if err != nil {
				// 	return err
				// }
			}

			return nil
		})
		if err != nil {
			return err
		}
	} else {
		return errors.New("Database not initialized")
	}

	return nil

}

func UpdateOrderFromRFMS(request dto.ReportJobStatusResponseDTO) (orderManagement.OrderList, error) {
	newJobStatus := request.Body
	orderList := orderManagement.OrderList{}
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		est, err := StringToDatetime(newJobStatus.Est)
		if err != nil {
			return err
		}
		eta, err := StringToDatetime(newJobStatus.Eta)
		if err != nil {
			return err
		}
		lastUpdateTime, err := StringToDatetime(newJobStatus.MessageTime)
		if err != nil {
			return err
		}
		updateFields := []string{"job_status", "processing_status", "job_start_time", "expected_arrival_time", "end_location_id", "failed_reason", "last_update_time"}
		updateMap := utils.CreateMap(updateFields, newJobStatus.Status, newJobStatus.ProcessingStatus, est, eta, newJobStatus.LocationId, "", lastUpdateTime)
		var updatedJobList = []db_models.Jobs{}
		err = UpdateRecords(tx, &updatedJobList, "jobs", updateMap, "job_id = ?", newJobStatus.JobID)
		if err != nil {
			return err
		}

		currentJobId := newJobStatus.JobID
		currentJobStatus := newJobStatus.ProcessingStatus
		statusLocation := ""

		if newJobStatus.Status == "COMPLETED" {
			nextJobs := []db_models.Jobs{}
			FindRecords(tx, &nextJobs, "jobs", "order_id = ? and processing_status = ? order by create_id", updatedJobList[0].OrderID, "UNKNOWN")
			if len(nextJobs) > 0 {
				param := rfms.CreateJobRequest{JobNature: nextJobs[0].JobType, LocationID: nextJobs[0].EndLocationID}
				response := apiHandler.POST("/createJob", param)

				updateJobStatus := dto.ReportJobStatusResponseDTO{}
				err := json.Unmarshal(response, &updateJobStatus)
				if err != nil {
					return err
				}
				if updateJobStatus.ResponseMessage == "FAILED" {
					return errors.New("Create Job Failed")
				}
				est, err := StringToDatetime(updateJobStatus.Body.Est)
				if err != nil {
					return err
				}
				eta, err := StringToDatetime(updateJobStatus.Body.Eta)
				if err != nil {
					return err
				}
				lastUpdateTime, err := StringToDatetime(updateJobStatus.Body.MessageTime)
				if err != nil {
					return err
				}
				updateFields := []string{"job_status", "processing_status", "job_start_time", "expected_arrival_time", "end_location_id", "failed_reason", "last_update_time", "job_id"}
				updateMap := utils.CreateMap(updateFields, updateJobStatus.Body.Status, updateJobStatus.Body.ProcessingStatus, est, eta, updateJobStatus.Body.LocationId, "", lastUpdateTime, updateJobStatus.Body.JobID)
				var updatedJobList = []db_models.Jobs{}
				err = UpdateRecords(tx, &updatedJobList, "jobs", updateMap, "create_id = ?", nextJobs[0].CreateID)
				if err != nil {
					return err
				}
				jobsLogList, err := JobsToJobsLogs(updatedJobList)
				if err != nil {
					return err
				}
				err = AddRecords(tx, jobsLogList)
				if err != nil {
					return err
				}

				currentJobId = updateJobStatus.Body.JobID
				currentJobStatus = updateJobStatus.Body.ProcessingStatus
				statusLocation = updatedJobList[0].StatusLocation
			}

			processingStatus := getProcessingStatusFromJob(currentJobStatus, statusLocation)
			orderStatus := "PROCESSING"
			if currentJobStatus == "ARRIVED" && statusLocation == "END_LOCATION" {
				orderStatus = "COMPLETED"
			}

			updatedOrderList := []db_models.Orders{}
			updateOrderFields := []string{"order_status", "processing_status", "job_id"}
			updateOrderMap := utils.CreateMap(updateOrderFields, orderStatus, processingStatus, currentJobId)
			UpdateRecords(tx, &updatedOrderList, "orders", updateOrderMap, "order_id = ? and order_status = ?", updatedJobList[0].OrderID, "CREATED")
			ordersLogList, err := OrdersToOrdersLogs(0, updatedOrderList)
			if err != nil {
				return err
			}
			err = AddRecords(tx, ordersLogList)
			if err != nil {
				return err
			}
			orderList, err = OrderListToOrderResponse(updatedOrderList)
			if err != nil {
				return err
			}
		}
		// err = UpdateRecords(tx, &updatedOrderList, "orders", updateOrderMap, "order_id = ?", order.OrderID)
		// if err != nil {
		// 	return err
		// }
		return nil
	})
	if err != nil {
		return orderList, err
	}
	return orderList, nil
}

func getProcessingStatusFromJob(jobStatus string, statusLocation string) string {
	switch jobStatus {
	case "PLANNING":
		return "PLANNING_TO_" + statusLocation
	case "GOING_TO_LOCATION":
		return "GOING_TO_" + statusLocation
	case "QUEUING":
		return "QUEUING_AT_" + strings.Replace(statusLocation, "LOCATION", "BAY", -1)
	case "ARRIVING":
		return "ARRIVING_TO_" + statusLocation
	case "ARRIVED":
		return "ARRIVED_TO_" + statusLocation
	case "MOVING_TO_LAYBY":
		return "MOVING_TO_LAYBY"
	default:
		return "UNKNOWN"
	}
}

func GetUpdateJobFields(rawResponse []byte) (dto.ReportJobStatusDTO, map[string]interface{}, error) {
	rawString := string(rawResponse)
	updateJobStatus := dto.ReportJobStatusResponseDTO{}
	err := json.Unmarshal(rawResponse, &updateJobStatus)
	if err != nil {
		return updateJobStatus.Body, nil, errors.New("Failed to phrase create job response")
	}
	log.Println(string(rawResponse))
	log.Println(updateJobStatus)
	var updateMap = make(map[string]interface{})
	if strings.Contains(rawString, "status") {
		if updateJobStatus.Body.Status == "FAILED" {
			return updateJobStatus.Body, nil, errors.New("RFMS Returned Fail")
		}
		updateMap["job_status"] = updateJobStatus.Body.Status
	}
	if strings.Contains(rawString, "est") {
		updateMap["expected_start_time"] = updateJobStatus.Body.Est
	}
	if strings.Contains(rawString, "eta") {
		updateMap["expected_arrival_time"] = updateJobStatus.Body.Eta
	}
	if strings.Contains(rawString, "processingStatus") {
		updateMap["processing_status"] = updateJobStatus.Body.ProcessingStatus
	}
	// if (strings.Contains(rawString, "payloadId")){
	// 	updateFields = append(updateFields, "location")
	// }

	updateMap["last_update_time"] = utils.GetTimeNowString()
	// log.Println(updateMap)
	return updateJobStatus.Body, updateMap, nil
}

func OrderIdsToIntArray(orderIds string) ([]int, error) {
	stringArray := strings.Split(orderIds, ",")
	var ret []int
	for _, value := range stringArray {
		v, err := strconv.Atoi(value)
		if err != nil {
			return ret, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func OrderDtoToOrderList(dto.UpdateDeliveryOrderDTO) ([]db_models.Orders, error) {
	orderList := []db_models.Orders{}

	return orderList, nil
}

func JobsToJobsLogs(jobs []db_models.Jobs) ([]db_models.JobsLogs, error) {
	var jobsLogs []db_models.JobsLogs

	bJson, err := json.Marshal(jobs)
	if err != nil {
		return jobsLogs, err
	}
	err = json.Unmarshal(bJson, &jobsLogs)
	if err != nil {
		return jobsLogs, err
	}

	for _, jobsLog := range jobsLogs {
		if jobsLog.ExpectedArrivalTime == "" {
			jobsLog.ExpectedArrivalTime = utils.TimeInt64ToString(0)
		}
		if jobsLog.LastUpdateTime == "" {
			jobsLog.LastUpdateTime = utils.TimeInt64ToString(0)
		}
	}

	return jobsLogs, nil
}

func OrdersToOrdersLogs(userId int, orders []db_models.Orders) ([]db_models.OrdersLogs, error) {
	var ordersLogs []db_models.OrdersLogs

	bJson, err := json.Marshal(orders)
	if err != nil {
		return ordersLogs, err
	}
	err = json.Unmarshal(bJson, &ordersLogs)
	if err != nil {
		return ordersLogs, err
	}

	for i, _ := range ordersLogs {
		ordersLogs[i].LastUpdateBy = userId
		if ordersLogs[i].OrderStartTime == "" {
			ordersLogs[i].OrderStartTime = utils.TimeInt64ToString(0)
		}
		if ordersLogs[i].ActualArrivalTime == "" {
			ordersLogs[i].ActualArrivalTime = utils.TimeInt64ToString(0)
		}
		if ordersLogs[i].ExpectedStartTime == "" {
			ordersLogs[i].ExpectedStartTime = utils.TimeInt64ToString(0)
		}
		if ordersLogs[i].ExpectedDeliveryTime == "" {
			ordersLogs[i].ExpectedDeliveryTime = utils.TimeInt64ToString(0)
		}
		if ordersLogs[i].ExpectedArrivalTime == "" {
			ordersLogs[i].ExpectedArrivalTime = utils.TimeInt64ToString(0)
		}
		if ordersLogs[i].LastUpdateTime == "" {
			ordersLogs[i].LastUpdateTime = utils.TimeInt64ToString(0)
		}
	}

	return ordersLogs, nil
}

func SchedulesToSchedulesLogs(userId int, schedules []db_models.Schedules) ([]db_models.SchedulesLogs, error) {
	var schedulesLogs []db_models.SchedulesLogs

	bJson, err := json.Marshal(schedules)
	if err != nil {
		return schedulesLogs, err
	}
	err = json.Unmarshal(bJson, &schedulesLogs)
	if err != nil {
		return schedulesLogs, err
	}

	for i, _ := range schedulesLogs {
		schedulesLogs[i].LastUpdateBy = userId
		if schedulesLogs[i].ScheduleCraeteTime == "" {
			schedulesLogs[i].ScheduleCraeteTime = utils.TimeInt64ToString(0)
		}
		if schedulesLogs[i].LastUpdateTime == "" {
			schedulesLogs[i].LastUpdateTime = utils.TimeInt64ToString(0)
		}
	}

	return schedulesLogs, nil
}

func RoutinesToRoutinesLogs(userId int, routines []db_models.Routines) ([]db_models.RoutinesLogs, error) {
	var routinesLogs []db_models.RoutinesLogs

	bJson, err := json.Marshal(routines)
	if err != nil {
		return routinesLogs, err
	}
	err = json.Unmarshal(bJson, &routinesLogs)
	if err != nil {
		return routinesLogs, err
	}

	for i, _ := range routinesLogs {
		routinesLogs[i].LastUpdateBy = userId
		if routinesLogs[i].LastUpdateTime == "" {
			routinesLogs[i].LastUpdateTime = utils.TimeInt64ToString(0)
		}
	}

	return routinesLogs, nil
}

func RoutinePatternToString(routinePattern orderManagement.RoutinePattern) (string, error) {
	var patternString string
	bJson, err := json.Marshal(routinePattern)
	if err != nil {
		return patternString, err
	}
	patternString = string(bJson)
	return patternString, nil
}

func StringToRoutinePattern(patternString string) (orderManagement.RoutinePattern, error) {
	var routinePattern orderManagement.RoutinePattern
	bJson := []byte(patternString)
	err := json.Unmarshal(bJson, &routinePattern)
	if err != nil {
		return routinePattern, err
	}
	return routinePattern, nil
}

func GetNextDeliveryDate(routinePattern orderManagement.RoutinePattern) (string, error) {

	// var patternString string
	// patternString = utils.GetTimeNowString()
	timeNow := time.Now()
	var nextRoutineDate = timeNow
	if routinePattern.Week != nil {
		targetDate := timeNow
		if routinePattern.Month != nil {
			currentMonth := timeNow.Month()
			var monthsDiffMin = 12
			for _, month := range routinePattern.Month {
				monthsDiff := month - int(currentMonth)
				if monthsDiff < 0 {
					monthsDiff += 12
				}
				if monthsDiff < monthsDiffMin {
					monthsDiffMin = monthsDiff
				}
			}
			targetDate = timeNow.AddDate(0, monthsDiffMin, 1-timeNow.Day())
		}
		fmt.Printf("targetDate: %s\n", targetDate)
		currentWeekday := targetDate.Weekday()
		var daysDiffMin = 7
		for _, week := range routinePattern.Week {
			daysDiff := week - int(currentWeekday)
			if daysDiff < 0 {
				daysDiff += 7
			}
			if daysDiff < daysDiffMin {
				daysDiffMin = daysDiff
			}
		}
		nextRoutineDate = targetDate.AddDate(0, 0, daysDiffMin)
	} else if routinePattern.Day != nil {
		currentMonth := timeNow.Month()
		var monthsDiffMin = 12
		var nextRoutineDay = 0

		if routinePattern.Month != nil {
			for _, month := range routinePattern.Month {
				monthsDiff := month - int(currentMonth)
				if monthsDiff < 0 {
					monthsDiff += 12
				}
				if monthsDiff < monthsDiffMin {
					monthsDiffMin = monthsDiff
				}
			}
			fmt.Printf("monthsDiffMin: %d\n", monthsDiffMin)
		} else {
			monthsDiffMin = 0
		}

		if monthsDiffMin == 0 {
			daysDiffMin := 31
			for _, day := range routinePattern.Day {
				currentDay := timeNow.Day()
				daysDiff := day - int(currentDay)
				if daysDiff > 0 && daysDiff < daysDiffMin {
					daysDiffMin = daysDiff
				}
			}
			fmt.Printf("daysDiffMin: %d\n", daysDiffMin)
			if daysDiffMin == 31 {
				monthsDiffMin = 1
				nextRoutineDay = routinePattern.Day[0]
			} else {
				nextRoutineDay = timeNow.Day() + daysDiffMin
			}
		} else {
			nextRoutineDay = routinePattern.Day[0]
		}
		nextRoutineYear := int(timeNow.Year())
		nextRoutineMonth := int(timeNow.Month()) + monthsDiffMin

		nextRoutineDate = time.Time.AddDate(time.Unix(0, 0), nextRoutineYear-1970, nextRoutineMonth-1, nextRoutineDay-1)
	} else if routinePattern.Month != nil {
		currentMonth := timeNow.Month()
		var monthsDiffMin = 12
		for _, month := range routinePattern.Month {
			monthsDiff := month - int(currentMonth)
			if monthsDiff < 0 {
				monthsDiff += 12
			}
			if monthsDiff < monthsDiffMin {
				monthsDiffMin = monthsDiff
			}
		}
		nextRoutineDate = timeNow.AddDate(0, monthsDiffMin, 1-timeNow.Day())
	}

	return nextRoutineDate.Format("20060102"), nil
}

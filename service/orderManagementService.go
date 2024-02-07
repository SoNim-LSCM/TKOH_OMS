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
	"github.com/SoNim-LSCM/TKOH_OMS/constants/scheduleStatus"
	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	"github.com/SoNim-LSCM/TKOH_OMS/models"
	dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"
	"github.com/SoNim-LSCM/TKOH_OMS/models/orderManagement"
	"github.com/SoNim-LSCM/TKOH_OMS/models/rfms"
	ws_model "github.com/SoNim-LSCM/TKOH_OMS/models/websocket"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"
	"github.com/SoNim-LSCM/TKOH_OMS/websocket"
	"gorm.io/gorm"
)

func FindOrdersForFrontPage(filterFields string, locationId int, filterValues ...interface{}) ([]db_models.Orders, error) {
	var orders []db_models.Orders
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		query := "SELECT *, C.location_name as start_location_name, D.location_name as end_location_name FROM tkoh_oms.orders LEFT JOIN tkoh_oms.locations C ON tkoh_oms.orders.start_location_id = C.location_id  LEFT JOIN tkoh_oms.locations D ON tkoh_oms.orders.end_location_id = D.location_id WHERE " + filterFields + " ORDER BY case when start_location_id = " + fmt.Sprint(locationId) + " then 1 else 99 end , case when end_location_id = " + fmt.Sprint(locationId) + " then 1 else 99 end ,  if (start_location_id = " + fmt.Sprint(locationId) + ", case when processing_status like 'UNLOADING' then 3 when processing_status like 'ARRIVED_START%' then 4 when processing_status like 'QUEUING_AT_START%' then 5 when processing_status like 'MOVING_TO_LAYBY_AREA' then 6 when processing_status like 'GOING_TO_START%' then 7 when processing_status like 'PLANNING_TO_START%' then 8 when processing_status = 'UNKNOWN' then 10 when processing_status = '' then 98 else 99 end , if (end_location_id = " + fmt.Sprint(locationId) + ", case when processing_status like 'UNLOADING' then 3 when processing_status like 'ARRIVED_END%' then 4 when processing_status like 'QUEUING_AT_END%' then 5 when processing_status like 'MOVING_TO_LAYBY_AREA' then 6 when processing_status like 'PLANNING_TO_END%' then 7 when processing_status like 'GOING_TO_END%' then 8 when processing_status like '%START%' then 9 when processing_status = 'UNKNOWN' then 10 when processing_status = '' then 98 else 99 end , 999) ) asc"
		if err := FindRecordsWithRaw(tx, &orders, query, filterValues...); err != nil {
			return errors.New("Failed to search: " + err.Error())
		}
		return nil
	})
	return orders, err
}

func FindOrders(filterFields string, filterValues ...interface{}) ([]db_models.Orders, error) {
	var orders []db_models.Orders
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		query := "SELECT *, C.location_name as start_location_name, D.location_name as end_location_name FROM tkoh_oms.orders LEFT JOIN tkoh_oms.locations C ON tkoh_oms.orders.start_location_id = C.location_id  LEFT JOIN tkoh_oms.locations D ON tkoh_oms.orders.end_location_id = D.location_id WHERE " + filterFields
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

func AddOrders(orderRequests []dto.AddDeliveryOrderDTO, userId int, orderCreatedType string) (orderManagement.OrderList, error) {
	var orderList orderManagement.OrderList
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var lastSchedule db_models.Schedules

		for _, orderRequest := range orderRequests {
			// create new schedule
			schedulesList := []db_models.Schedules{{ScheduleID: 0, ScheduleStatus: string(scheduleStatus.Created), ScheduleCraeteTime: utils.GetTimeNowString(), OrderType: orderRequest.OrderType, OrderCreatedType: orderCreatedType, NumberOfAmrRequire: orderRequest.NumberOfAmrRequire, RoutineID: orderRequest.RoutineID, LastUpdateTime: utils.GetTimeNowString()}}
			if err := AddRecords(tx, schedulesList); err != nil {
				return err
			}
			// check no of schedules
			if err := tx.Table("schedules").Last(&lastSchedule).Error; err != nil {
				return err
			}
			// translate order request to orders
			orders, err := OrderRequestToOrders(orderRequest, lastSchedule.ScheduleID, userId, orderCreatedType)
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
		}

		return nil
	})
	return orderList, err
}

func AddRoutines(routineRequest dto.AddRoutineDTO, userId int) (orderManagement.RoutineOrderList, error) {
	database.CheckDatabaseConnection()
	var routineOrderList orderManagement.RoutineOrderList
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

func TriggerOrderOrderIds(orderIds string) (orderManagement.OrderList, error) {
	var orderList orderManagement.OrderList
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		jobList := []db_models.Jobs{}
		err := FindRecordsWithRaw(tx, &jobList, "SELECT * FROM jobs WHERE order_id IN ("+orderIds+") and job_status = ? ", "PROCESSING")
		if err != nil {
			return err
		}
		jobIdList := []int{}
		for _, orders := range jobList {
			jobIdList = append(jobIdList, orders.JobID)
		}
		param := rfms.TriggerHandlingJobRequest{JobIdList: jobIdList}
		response, err := apiHandler.Post("/triggerHandlingJob", param)
		if err != nil {
			return err
		}
		triggerOrderResponse := models.FailResponseHeader{}
		err = json.Unmarshal(response, &triggerOrderResponse)
		if err != nil {
			return err
		}
		if triggerOrderResponse.ResponseCode != 200 {
			return errors.New("RFMS with Fail Response")
		}
		return nil
	})
	return orderList, err
}

func TriggerOrderScheduleId(scheduleId int) (orderManagement.OrderList, error) {
	var orderList orderManagement.OrderList
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		jobList := []db_models.Jobs{}
		err := FindRecordsWithRaw(tx, &jobList, "SELECT * FROM jobs WHERE schedule_id = ? and job_status = ? ", scheduleId, "PROCESSING")
		if err != nil {
			return err
		}
		jobIdList := append([]int{}, jobList[0].JobID)
		param := rfms.TriggerHandlingJobRequest{JobIdList: jobIdList}
		response, err := apiHandler.Post("/triggerHandlingJob", param)
		if err != nil {
			return err
		}
		triggerOrderResponse := models.FailResponse{}
		err = json.Unmarshal(response, &triggerOrderResponse)
		if err != nil {
			return err
		}
		if triggerOrderResponse.Header.ResponseCode != 200 {
			return errors.New("RFMS with Fail Response")
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
		if FindRecords(tx, &schedules, "schedules", "schedule_id = ? AND schedule_status <> ?", request.ScheduleID, "CANCELED") != nil {
			return errors.New("Failed to find schedule with schedule id")
		}

		orders := []db_models.Orders{}
		if FindRecords(tx, &orders, "orders", "schedule_id = ? AND order_status = ?", request.ScheduleID, orderStatus.ToBeCreated) != nil {
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
				uploadOrders, err := OrderRequestToOrders(orderRequest, request.ScheduleID, 6, schedules[0].OrderCreatedType)
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

		if err := FindRecords(tx, &orders, "orders", "schedule_id = ? AND order_status = ?", scheduleId, orderStatus.ToBeCreated); err != nil {
			return err
		}

		if len(orders) == 0 || len(schedules) == 0 {
			return errors.New("Orders not found")
		}

		if len(orders) != schedules[0].NumberOfAmrRequire {
			return errors.New("Cancel fail, some orders already started")
		}

		updateMap := utils.CreateMap([]string{"schedule_status"}, orderStatus.Canceled)
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
	updateMap := utils.CreateMap([]string{"routine_pattern", "number_of_amr_require", "start_location_id", "end_location_id", "expected_start_time", "expected_delivery_time", "is_active"}, routinePattern, request.NumberOfAmrRequire, request.StartLocationID, request.EndLocationID, expectedStartTime, expectedDeliveryTime, request.IsActive)
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

func OrderRequestToOrders(orderRequests dto.AddDeliveryOrderDTO, scheduleNo int, userId int, orderCreatedType string) ([]db_models.Orders, error) {
	var orders []db_models.Orders
	for i := 0; i < orderRequests.NumberOfAmrRequire; i++ {
		var err error
		var order db_models.Orders
		order.ScheduleID = scheduleNo
		// order.OrderID = i + orderNo
		order.OrderID = 0
		order.OrderType = orderRequests.OrderType
		order.OrderCreatedType = orderCreatedType
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
		order.StartLocationID = orderRequests.StartLocationID
		order.StartLocationName = orderRequests.StartLocationName
		order.EndLocationID = orderRequests.EndLocationID
		order.EndLocationName = orderRequests.EndLocationName
		order.ExpectedStartTime, err = StringToDatetime(orderRequests.ExpectedStartTime)
		if err != nil {
			return orders, err
		}
		order.ExpectedDeliveryTime, err = StringToDatetime(orderRequests.ExpectedDeliveryTime)
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

func BackgroundRoutinesToSchedules() error {
	routines, err := FindRoutines("is_active = ?", true)
	if err != nil {
		return err
	}
	if len(routines) == 0 {
		return nil
	}
	log.Print(routines)
	addDeliveryOrderDTO, err := RoutinesToAddDeliveryOrderDTO(routines)
	if err != nil {
		return err
	}
	log.Print(addDeliveryOrderDTO)
	_, err = AddOrders(addDeliveryOrderDTO, 0, "ROUTINE")
	if err != nil {
		return err
	}

	return nil
}

func BackgroundInitOrderToRFMS() error {

	if database.CheckDatabaseConnection() {

		jobs := []db_models.Jobs{}

		err := database.DB.Transaction(func(tx *gorm.DB) error {

			if err := FindRecords(tx, &jobs, "jobs", "job_status = ?", "TO_BE_CREATED"); err != nil {
				return err
			}

			return nil
		})

		for _, job := range jobs {
			err := database.DB.Transaction(func(tx *gorm.DB) error {
				param := rfms.CreateJobRequest{JobNature: job.JobType, LocationID: job.EndLocationID, RobotID: job.RobotID, PayloadID: job.PayloadID}
				if job.JobType == "PARK" {
					param = rfms.CreateJobRequest{JobNature: job.JobType, RobotID: job.RobotID, PayloadID: job.PayloadID}
				}
				response, err := apiHandler.Post("/createJob", param)
				if err != nil {
					return err
				}

				updateJobStatus := dto.ReportJobStatusResponseDTO{}
				err = json.Unmarshal(response, &updateJobStatus)
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

				updateFields := []string{"job_id", "job_status", "processing_status", "job_start_time", "expected_arrival_time", "last_update_time"}
				updateMap := utils.CreateMap(updateFields, updateJobStatus.Body.JobID, updateJobStatus.Body.Status, updateJobStatus.Body.ProcessingStatus, est, eta, lastUpdateTime)
				var updatedJobList = []db_models.Jobs{}
				err = UpdateRecords(tx, &updatedJobList, "jobs", updateMap, "create_id = ?", job.CreateID)
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

				return nil
			})
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

		orders := []db_models.Orders{}

		err = database.DB.Transaction(func(tx *gorm.DB) error {

			if err := FindRecordsWithRaw(tx, &orders, "SELECT * FROM orders WHERE order_status = ? order by expected_start_time asc", "TO_BE_CREATED"); err != nil {
				return err
			}

			if len(orders) == 0 {
				return nil
			}

			return nil
		})

		for _, order := range orders {
			expectedStartTimeString, err := StringToDatetime(order.ExpectedStartTime)
			if err != nil {

			}
			expectedStartTime, err := time.Parse("2006-01-02 15:04:05", expectedStartTimeString)
			if err != nil {

			}
			timeNowUTCForCompare := time.Now().Add(8 * time.Hour)
			// log.Printf("expectedStartTime: %s, time now: %s", expectedStartTime, timeNowUTCForCompare)
			if expectedStartTime.After(timeNowUTCForCompare) {
				break
			}
			err = database.DB.Transaction(func(tx *gorm.DB) error {
				jobNatures := []string{}
				locations := []int{}
				statusLocation := []string{}
				// orderTypes := strings.Split(order.OrderType, "_")
				// jobNatures = append(jobNatures, orderTypes[0])
				switch order.OrderType {
				case "PICK_AND_DELIVERY":
					jobNatures = append(jobNatures, "PICK")
					locations = append(locations, order.StartLocationID)
					statusLocation = append(statusLocation, "START_LOCATION")

					jobNatures = append(jobNatures, "DELIVERY")
					locations = append(locations, order.EndLocationID)
					statusLocation = append(statusLocation, "END_LOCATION")
				case "PICK_DELIVERY_PARK":
					jobNatures = append(jobNatures, "PICK")
					locations = append(locations, order.StartLocationID)
					statusLocation = append(statusLocation, "START_LOCATION")

					jobNatures = append(jobNatures, "DELIVERY")
					locations = append(locations, order.EndLocationID)
					statusLocation = append(statusLocation, "END_LOCATION")

					jobNatures = append(jobNatures, "PARK")
					locations = append(locations, order.EndLocationID)
					statusLocation = append(statusLocation, "PARKING")
				case "PICK_ONLY":
					jobNatures = append(jobNatures, "PICK")
					locations = append(locations, order.StartLocationID)
					statusLocation = append(statusLocation, "END_LOCATION")
				case "DELIVERY_ONLY":
					jobNatures = append(jobNatures, "DELIVERY")
					locations = append(locations, order.EndLocationID)
					statusLocation = append(statusLocation, "END_LOCATION")
				case "DELIVERY_PARK":
					jobNatures = append(jobNatures, "DELIVERY")
					locations = append(locations, order.EndLocationID)
					statusLocation = append(statusLocation, "END_LOCATION")

					jobNatures = append(jobNatures, "PARK")
					locations = append(locations, order.EndLocationID)
					statusLocation = append(statusLocation, "PARKING")
				case "PARK_ONLY":
					jobNatures = append(jobNatures, "PARK")
					locations = append(locations, order.EndLocationID)
					statusLocation = append(statusLocation, "PARKING")
				default:
					return errors.New("Unknown Order Type")
				}

				param := rfms.CreateJobRequest{JobNature: jobNatures[0], LocationID: locations[0]}
				response, err := apiHandler.Post("/createJob", param)
				if err != nil {
					return err
				}

				updateJobStatus := dto.ReportJobStatusResponseDTO{}
				err = json.Unmarshal(response, &updateJobStatus)
				if err != nil {
					return err
				}
				// log.Printf("/createJob response: %f", updateJobStatus)
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
				newJob := db_models.Jobs{OrderID: order.OrderID, JobID: updateJobStatus.Body.JobID, JobType: jobNatures[0], JobStatus: updateJobStatus.Body.Status, ProcessingStatus: updateJobStatus.Body.ProcessingStatus, JobStartTime: est, ExpectedArrivalTime: eta, EndLocationID: updateJobStatus.Body.LocationId, FailedReason: updateJobStatus.FailReason, LastUpdateTime: lastUpdateTime, StatusLocation: statusLocation[0]}
				newJobs := append([]db_models.Jobs{}, newJob)
				for i, jobNature := range jobNatures {
					if i > 0 {
						defaultTime, err := StringToDatetime("")
						if err != nil {
							return err
						}
						newJob := db_models.Jobs{OrderID: order.OrderID, JobID: 0, JobType: jobNature, JobStatus: "WAIT_FOR_PREVIOUS_JOB_END", ProcessingStatus: "UNKNOWN", JobStartTime: defaultTime, ExpectedArrivalTime: defaultTime, EndLocationID: locations[i], FailedReason: updateJobStatus.FailReason, LastUpdateTime: defaultTime, StatusLocation: statusLocation[i]}
						newJobs = append(newJobs, newJob)
					}
				}

				// log.Print(newJobs)
				err = AddRecords(tx, newJobs)
				if err != nil {
					return err
				}

				startTime := utils.GetTimeNowString()
				if err != nil {
					return errors.New("Fail translate arrivalTime")
				}

				updatedList := []db_models.Orders{}
				updateFields := []string{"order_status", "job_id", "order_start_time"}
				updateMap := utils.CreateMap(updateFields, orderStatus.Created, updateJobStatus.Body.JobID, startTime)
				err = UpdateRecords(tx, &updatedList, "orders", updateMap, "order_id = ?", order.OrderID)
				log.Print(updatedList)
				if err != nil {
					return err
				}
				ordersLogList, err := OrdersToOrdersLogs(0, updatedList)
				if err != nil {
					return err
				}
				err = AddRecords(tx, ordersLogList)
				if err != nil {
					return err
				}
				websocket.SendBoardcastMessage(ws_model.GetUpdateOrderStatusResponse(updatedList[0].OrderID, updatedList[0].OrderStatus, newJob.PayloadID, updatedList[0].ProcessingStatus, append([]string{}, newJob.RobotID), updatedList[0].ScheduleID))

				return nil
			})
		}

		if err != nil {
			return err
		}
	} else {
		return errors.New("Database not initialized")
	}

	return nil
}

func RoutinesToAddDeliveryOrderDTO(routines []db_models.Routines) ([]dto.AddDeliveryOrderDTO, error) {
	addDeliveryOrderDTO := []dto.AddDeliveryOrderDTO{}
	today := time.Now().Format("20060102")
	for _, routine := range routines {
		pattern := orderManagement.RoutinePattern{}
		// bJson, err := json.Marshal(routine.RoutinePattern)
		// if err != nil {
		// 	return addDeliveryOrderDTO, err
		// }
		err := json.Unmarshal([]byte(routine.RoutinePattern), &pattern)
		if err != nil {
			return addDeliveryOrderDTO, err
		}
		nextDeliveryDate, err := GetNextDeliveryDate(pattern)
		if err != nil {
			return addDeliveryOrderDTO, err
		}
		if nextDeliveryDate == today {
			expectedStartTime := strings.Replace(routine.ExpectedStartTime, "1970-01-01", strings.Split(utils.GetTimeNowString(), " ")[0], 1)
			expectedDeliveryTime := strings.Replace(routine.ExpectedDeliveryTime, "1970-01-01", strings.Split(utils.GetTimeNowString(), " ")[0], 1)
			addDeliveryOrderDTO = append(addDeliveryOrderDTO, dto.AddDeliveryOrderDTO{OrderType: routine.OrderType, NumberOfAmrRequire: routine.NumberOfAmrRequire, StartLocationID: routine.StartLocationID, StartLocationName: routine.StartLocationName, ExpectedStartTime: expectedStartTime, EndLocationID: routine.EndLocationID, EndLocationName: routine.EndLocationName, RoutineID: routine.RoutineID, ExpectedDeliveryTime: expectedDeliveryTime})
		}
	}
	return addDeliveryOrderDTO, nil
}

func UpdateOrderFromRFMS(request dto.ReportJobStatusDTO) (orderManagement.OrderList, error) {
	newJobStatus := request
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
		updateFields := []string{"job_status", "processing_status", "job_start_time", "expected_arrival_time", "end_location_id", "failed_reason", "last_update_time", "robot_id", "payload_id"}
		updateMap := utils.CreateMap(updateFields, newJobStatus.Status, newJobStatus.ProcessingStatus, est, eta, newJobStatus.LocationId, "", lastUpdateTime, newJobStatus.RobotID, newJobStatus.PayloadID)
		var updatedJobList = []db_models.Jobs{}
		err = UpdateRecords(tx, &updatedJobList, "jobs", updateMap, "job_id = ?", newJobStatus.JobID)
		if err != nil {
			return err
		}

		currentJobId := newJobStatus.JobID
		currentJobStatus := newJobStatus.ProcessingStatus
		statusLocation := updatedJobList[0].StatusLocation
		orderStatus := "PROCESSING"

		if newJobStatus.Status == "COMPLETED" {
			nextJobs := []db_models.Jobs{}
			FindRecords(tx, &nextJobs, "jobs", "order_id = ? and job_status = ? order by create_id", updatedJobList[0].OrderID, "WAIT_FOR_PREVIOUS_JOB_END")
			if len(nextJobs) > 0 {
				updateFields := []string{"job_status", "robot_id", "payload_id"}
				updateMap := utils.CreateMap(updateFields, "TO_BE_CREATED", newJobStatus.RobotID, newJobStatus.PayloadID)
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
			} else {
				orderStatus = "COMPLETED"
			}
		}
		arrivalTime := ""
		if orderStatus == "COMPLETED" {
			arrivalTime = utils.GetTimeNowString()
		}
		arrivalTime, err = StringToDatetime(arrivalTime)
		if err != nil {
			return errors.New("Fail translate arrivalTime")
		}
		log.Printf("currentJobStatus: %s, statusLocation: %s", currentJobStatus, statusLocation)
		processingStatus := getProcessingStatusFromJob(currentJobStatus, statusLocation)

		updatedOrderList := []db_models.Orders{}
		updateOrderFields := []string{"order_status", "processing_status", "job_id", "actual_arrival_time"}
		updateOrderMap := utils.CreateMap(updateOrderFields, orderStatus, processingStatus, currentJobId, arrivalTime)
		UpdateRecords(tx, &updatedOrderList, "orders", updateOrderMap, "order_id = ?", updatedJobList[0].OrderID)
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
	case "GOING_TO":
		return "GOING_TO_" + statusLocation
	case "QUEUING":
		return "QUEUING_AT_" + strings.Replace(statusLocation, "LOCATION", "BAY", -1)
	case "ARRIVING":
		return "ARRIVING_TO_" + statusLocation
	case "ARRIVED":
		return "ARRIVED_" + statusLocation
	case "MOVING_TO_LAYBY":
		return "MOVING_TO_" + statusLocation + "_LAYBY"
	case "PARKING":
		return "PARKING"
	case "UNLOADING":
		return "UNLOADING"
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

	for _, routinesLog := range routinesLogs {
		routinesLog.LastUpdateBy = userId
		if routinesLog.LastUpdateTime == "" {
			routinesLog.LastUpdateTime = utils.TimeInt64ToString(0)
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
	if routinePattern.Week != nil && len(routinePattern.Week) > 0 {
		targetDate := timeNow
		if routinePattern.Month != nil && len(routinePattern.Month) > 0 {
			// currentMonth := timeNow.Format("01")
			currentMonth := timeNow.Month()
			var monthsDiffMin = 12
			for _, month := range routinePattern.Month {
				monthsDiff := (month - int(currentMonth))
				if monthsDiff < 0 {
					monthsDiff += 12
				}
				if monthsDiff < monthsDiffMin {
					monthsDiffMin = monthsDiff
				}
			}
			if monthsDiffMin > 0 {
				targetDate = timeNow.AddDate(0, monthsDiffMin, 1-timeNow.Day())
			}
		}
		log.Printf("targetDate: %s\n", targetDate)
		currentWeekday := targetDate.Weekday()
		var daysDiffMin = 7
		for _, week := range routinePattern.Week {
			daysDiff := (week - int(currentWeekday))
			if daysDiff < 0 {
				daysDiff += 7
			}
			if daysDiff < daysDiffMin {
				daysDiffMin = daysDiff
			}
		}
		nextRoutineDate = targetDate.AddDate(0, 0, daysDiffMin)
	} else if routinePattern.Day != nil && len(routinePattern.Day) > 0 {
		// timeNow.Date(timeNow.Year(), timeNow.Month(), routinePattern.Day)
		currentMonth := timeNow.Month()
		var monthsDiffMin = 12
		// var nextRoutineDay = 0

		if routinePattern.Month != nil && len(routinePattern.Month) > 0 {
			for _, month := range routinePattern.Month {
				monthsDiff := (month - int(currentMonth))
				if monthsDiff < 0 {
					monthsDiff += 12
				}
				if monthsDiff < monthsDiffMin {
					monthsDiffMin = monthsDiff
				}
			}
		} else {
			monthsDiffMin = 0
		}

		if monthsDiffMin == 0 {
			added := false
			for _, day := range routinePattern.Day {
				if day > timeNow.Day() {
					nextRoutineDate = time.Date(timeNow.Year(), timeNow.Month(), day, 0, 0, 0, 0, time.Now().Location())
					added = true
					break
				}
			}
			if !added {
				tempDate := timeNow.AddDate(0, monthsDiffMin, 0)
				nextRoutineDate = time.Date(tempDate.Year(), tempDate.Month()+1, routinePattern.Day[0], 0, 0, 0, 0, time.Now().Location())
			}
		} else {
			tempDate := timeNow.AddDate(0, monthsDiffMin, 0)
			nextRoutineDate = time.Date(tempDate.Year(), tempDate.Month(), routinePattern.Day[0], 0, 0, 0, 0, time.Now().Location())
		}

	} else if routinePattern.Month != nil && len(routinePattern.Month) > 0 {
		currentMonth := timeNow.Month()
		var monthsDiffMin = 12
		for _, month := range routinePattern.Month {
			monthsDiff := (month - int(currentMonth))
			if monthsDiff < 0 {
				monthsDiff += 12
			}
			if monthsDiff < monthsDiffMin {
				monthsDiffMin = monthsDiff
			}
		}
		if monthsDiffMin != 0 {
			tempDate := timeNow.AddDate(0, monthsDiffMin, 0)
			nextRoutineDate = time.Date(tempDate.Year(), tempDate.Month(), 1, 0, 0, 0, 0, time.Now().Location())
		} else {
			nextRoutineDate = timeNow
		}
		// nextRoutineDate = timeNow.AddDate(0, monthsDiffMin, 1-timeNow.Day())
	}

	return nextRoutineDate.Format("20060102"), nil
}

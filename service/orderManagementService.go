package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/SoNim-LSCM/TKOH_OMS/constants/orderStatus"
	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"
	"github.com/SoNim-LSCM/TKOH_OMS/models/orderManagement"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"
	"gorm.io/gorm"
)

func FindOrders(filterFields interface{}, filterValues ...interface{}) ([]db_models.Orders, error) {
	var orders []db_models.Orders
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := FindRecords(tx, &orders, "orders", filterFields, filterValues...); err != nil {
			return errors.New("Failed to search: " + err.Error())
		}
		return nil
	})
	return orders, err

}

func AddOrders(orderRequest dto.AddDeliveryOrderDTO, userId int) (orderManagement.OrderList, error) {
	database.CheckDatabaseConnection()
	log.Printf("mysql query: AddOrder\n")
	var orderList orderManagement.OrderList
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var lastSchedule db_models.Schedules

		// create new schedule
		schedulesList := []db_models.Schedules{{ScheduleID: 0, ScheduleStatus: "CREATED", ScheduleCraeteTime: utils.GetTimeNowString(), OrderType: orderRequest.OrderType, NumberOfAmrRequire: orderRequest.NumberOfAmrRequire, LastUpdateTime: utils.GetTimeNowString()}}
		fmt.Println(schedulesList)
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

func AddRoutines(orderRequest dto.AddRoutineDTO, userId int) (orderManagement.RoutineOrderList, error) {
	database.CheckDatabaseConnection()
	log.Printf("mysql query: AddRoutineOrder\n")
	var orderList orderManagement.RoutineOrderList
	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// translate order request to routines
		routines, err := RoutineRequestToRoutines(orderRequest, userId)
		if err != nil {
			return err
		}
		// create new orders
		if err := AddRecords(tx, routines); err != nil {
			return err
		}
		// translate new routine to routine response
		orderList, err = RoutineListToRoutineResponse(GetRoutines())
		if err != nil {
			return err
		}

		return nil
	})

	return orderList, err
}

func TriggerOrderOrderIds(orderId []int) (orderManagement.OrderList, error) {
	var orders []db_models.Orders
	var orderList orderManagement.OrderList
	updateFields := []string{"order_status", "order_start_time"}
	timeNow := utils.GetTimeNowString()
	updateMap := utils.CreateMap(updateFields, string(orderStatus.Processing), timeNow)
	// updateValues := []string{string(orderStatus.Processing), timeNow}

	database.CheckDatabaseConnection()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := UpdateRecords(tx, &orders, "orders", updateMap, "order_id IN ?", orderId)
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
				updateMap := utils.CreateMap([]string{"schedule_id", "start_location_id", "end_location_id", "expected_start_time", "expected_delivey_time"}, request.ScheduleID, request.StartLocationID, request.EndLocationID, expectedStartTime, expectedDeliveryTime)
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
				fmt.Println(orderRequest)
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

func OrderRequestToOrders(orderRequest dto.AddDeliveryOrderDTO, scheduleNo int, userId int) ([]db_models.Orders, error) {
	var orders []db_models.Orders
	log.Printf("mysql query: OrderRequestToOrders\n")
	for i := 0; i < orderRequest.NumberOfAmrRequire; i++ {
		var err error
		var order db_models.Orders
		order.ScheduleID = scheduleNo
		// order.OrderID = i + orderNo
		order.OrderID = 0
		order.OrderType = orderRequest.OrderType
		order.OrderCreatedType = "ADHOC"
		order.OrderCreatedBy = userId
		order.OrderStatus = "CREATED"
		order.OrderStartTime, err = StringToDatetime("")
		if err != nil {
			return orders, err
		}
		order.ActualArrivalTime, err = StringToDatetime("")
		if err != nil {
			return orders, err
		}
		order.StartLocationID = orderRequest.StartLocationID
		order.EndLocationID = orderRequest.EndLocationID
		order.ExpectedStartTime, err = StringToDatetime(order.ExpectedStartTime)
		if err != nil {
			return orders, err
		}
		order.ExpectedDeliveryTime, err = StringToDatetime(order.ExpectedDeliveryTime)
		if err != nil {
			return orders, err
		}
		order.ExpectedArrivalTime, err = StringToDatetime("")
		if err != nil {
			return orders, err
		}
		order.ProcessingStatus = "PLANNING_TO_START_LOCATION"
		order.LastUpdateTime = utils.GetTimeNowString()
		order.LastUpdateBy = userId
		orders = append(orders, order)
	}
	return orders, nil
}

func OrderListToOrderResponse(orderList []db_models.Orders) (orderManagement.OrderList, error) {
	log.Printf("mysql query: OrderListToOrderResponse\n")
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
	log.Println(orderListResponse)
	return orderListResponse, nil
}

func RoutineRequestToRoutines(routinesRequest dto.AddRoutineDTO, userId int) ([]db_models.Routines, error) {
	var routinesList []db_models.Routines
	var routines db_models.Routines
	log.Printf("mysql query: OrderRequestToOrders\n")
	bJson, err := json.Marshal(routinesRequest)
	if err != nil {
		return routinesList, err
	}
	json.Unmarshal(bJson, &routines)
	if err != nil {
		return routinesList, err
	}
	routines.LastUpdateBy = userId
	routines.LastUpdateTime = utils.GetTimeNowString()
	routines.RoutinePattern, err = RoutinePatternToString(routinesRequest.RoutinePattern)
	if err != nil {
		return routinesList, nil
	}
	routines.ExpectedDeliveryTime = utils.TimeInt64ToString(0)
	routinesList = append(routinesList, routines)
	return routinesList, nil
}

func RoutineListToRoutineResponse(routineList []db_models.Routines) (orderManagement.RoutineOrderList, error) {
	log.Printf("mysql query: OrderListToOrderResponse\n")
	var orderListResponse orderManagement.RoutineOrderList
	jsonString, err := json.Marshal(routineList)
	if err != nil {
		return orderListResponse, err
	}
	json.Unmarshal(jsonString, &orderListResponse)
	if err != nil {
		return orderListResponse, err
	}
	for i, routine := range routineList {
		var err error
		routinePattern, err := StringToRoutinePattern(routine.RoutinePattern)
		if err != nil {
			return orderListResponse, err
		}
		orderListResponse[i].RoutinePattern = routinePattern
		orderListResponse[i].NextDeliveryDate, err = GetNextDeliveryDate(routinePattern)
		if err != nil {
			return orderListResponse, err
		}
	}
	log.Println(orderListResponse)
	return orderListResponse, nil
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

	// -------------------- Not Implemented -----------------------------

	var patternString string
	patternString = utils.GetTimeNowString()
	return patternString, nil
}

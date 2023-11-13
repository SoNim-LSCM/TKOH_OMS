package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/SoNim-LSCM/TKOH_OMS/constants/orderStatus"
	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"
	"github.com/SoNim-LSCM/TKOH_OMS/models/orderManagement"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"
)

func AddOrders(orderRequest dto.AddDeliveryOrderDTO, userId int) (orderManagement.OrderList, error) {
	database.CheckDatabaseConnection()
	log.Printf("mysql query: AddOrder\n")
	var schedulesCount int64
	var ordersCount int64
	var orderList orderManagement.OrderList
	timeNow, err := StringToDatetime(time.Now().Format("2006-01-02T15:04:05"))
	if err != nil {
		return orderList, err
	}
	// check no of schedules
	if err := database.DB.Table("schedules").Count(&schedulesCount).Error; err != nil {
		return orderList, err
	}
	// check no of orders
	if err := database.DB.Table("orders").Count(&ordersCount).Error; err != nil {
		return orderList, err
	}
	// translate order request to orders
	orders, err := OrderRequestToOrders(orderRequest, int(schedulesCount)+1, userId, int(ordersCount)+1)
	if err != nil {
		return orderList, err
	}
	// create new schedule
	if err := database.DB.Create(db_models.Schedules{ScheduleID: int(schedulesCount + 1), ScheduleStatus: "CREATED", ScheduleCraeteTime: timeNow, OrderType: orderRequest.OrderType, NumberOfAmrRequire: orderRequest.NumberOfAmrRequire}).Error; err != nil {
		return orderList, err
	}
	// create new orders
	if err := database.DB.Create(&orders).Error; err != nil {
		return orderList, err
	}
	// translate new orders to order response
	orderList, err = OrderListToOrderResponse(orders)
	if err != nil {
		return orderList, err
	}
	return orderList, nil
}

func TriggerOrderOrderIds(orderId []int) (orderManagement.OrderList, error) {
	var orders []db_models.Orders
	var orderList orderManagement.OrderList
	updateFields := []string{"order_status", "order_start_time"}
	timeNow := utils.GetTimeNowString()
	updateMap := utils.CreateMap(updateFields, string(orderStatus.Processing), timeNow)
	// updateValues := []string{string(orderStatus.Processing), timeNow}
	err := UpdateRecords(&orders, "orders", updateMap, "order_id IN ?", orderId)
	if err != nil {
		return orderList, err
	}
	orderList, err = OrderListToOrderResponse(orders)
	if err != nil {
		return orderList, err
	}
	return orderList, nil
}

func TriggerOrderScheduleId(scheduleId int) (orderManagement.OrderList, error) {
	var orders []db_models.Orders
	var orderList orderManagement.OrderList
	updateFields := []string{"order_status", "order_start_time"}
	timeNow := utils.GetTimeNowString()
	updateMap := utils.CreateMap(updateFields, string(orderStatus.Processing), timeNow)
	// updateValues := []string{string(orderStatus.Processing), timeNow}
	err := UpdateRecords(&orders, "orders", updateMap, "schedule_id = ?", scheduleId)
	if err != nil {
		return orderList, err
	}
	orderList, err = OrderListToOrderResponse(orders)
	if err != nil {
		return orderList, err
	}

	return orderList, nil
}

func UpdateOrders(request dto.UpdateDeliveryOrderDTO) (orderManagement.OrderList, error) {
	var orderList orderManagement.OrderList

	schedules := []db_models.Schedules{}
	if FindRecords(&schedules, "schedules", "schedule_id = ?", request.ScheduleID) != nil {
		return orderList, errors.New("Failed to find schedule with schedule id")
	}

	orders := []db_models.Orders{}
	if FindRecords(&orders, "orders", "schedule_id = ? AND order_status = ?", request.ScheduleID, orderStatus.Created) != nil {
		return orderList, errors.New("Failed to find order with schedule id")
	}

	if len(orders) != schedules[0].NumberOfAmrRequire {
		return orderList, errors.New("Update fail, some orders already started")
	}
	expectedStartTime, err := StringToDatetime(request.ExpectedStartTime)
	if err != nil {
		return orderList, errors.New("Fail translate expectedStartTime")
	}
	expectedDeliveryTime, err := StringToDatetime(request.ExpectedDeliveryTime)
	if err != nil {
		return orderList, errors.New("Fail translate expectedDeliveryTime")
	}

	updateMap := utils.CreateMap([]string{"number_of_amr_require"}, request.NumberOfAmrRequire)
	err = UpdateRecords(&[]db_models.Schedules{}, "schedules", updateMap, "schedule_id = ?", request.ScheduleID)
	if err != nil {
		return orderList, errors.New("Failed to update schedule table")
	}
	for i := 0; i < utils.Max(request.NumberOfAmrRequire, schedules[0].NumberOfAmrRequire); i++ {
		// cancel orders
		if i < (schedules[0].NumberOfAmrRequire - request.NumberOfAmrRequire) {
			updatedOrderList := []db_models.Orders{}
			updateMap := utils.CreateMap([]string{"order_status"}, orderStatus.Canceled)
			if UpdateRecords(&updatedOrderList, "orders", updateMap, "order_id = ?", orders[i].OrderID) != nil {
				return orderList, errors.New("Failed to translate orders to order response")
			}
			updatedOrderResponse, err := OrderListToOrderResponse(updatedOrderList)
			if err != nil {
				return orderList, errors.New("Failed to translate orders to order response")
			}
			orderList = append(orderList, updatedOrderResponse[0])
			// change orders
		} else if i < schedules[0].NumberOfAmrRequire {
			updatedOrderList := []db_models.Orders{}
			updateMap := utils.CreateMap([]string{"schedule_id", "start_location_id", "end_location_id", "expected_start_time", "expected_delivey_time"}, request.ScheduleID, request.StartLocationID, request.EndLocationID, expectedStartTime, expectedDeliveryTime)
			UpdateRecords(&updatedOrderList, "orders", updateMap, "order_id = ?", orders[i].OrderID)
			updatedOrderResponse, err := OrderListToOrderResponse(updatedOrderList)
			if err != nil {
				return orderList, errors.New("Failed to translate orders to order response")
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
			uploadOrders, err := OrderRequestToOrders(orderRequest, request.ScheduleID, 6, 0)
			if err != nil {
				return orderList, err
			}
			uploadOrders[0].OrderType = schedules[0].OrderType
			// create new orders
			if err := database.DB.Create(&uploadOrders).Error; err != nil {
				return orderList, err
			}
			updatedOrderResponse, err := OrderListToOrderResponse(uploadOrders)
			if err != nil {
				return orderList, errors.New("Failed to translate orders to order response")
			}
			orderList = append(orderList, updatedOrderResponse[0])
		}
	}

	return orderList, err
}

func CancelOrders(scheduleId int) (orderManagement.OrderList, error) {
	var orderList orderManagement.OrderList
	var schedules []db_models.Schedules
	var orders []db_models.Orders
	var updatedOrders []db_models.Orders

	if err := FindRecords(&schedules, "schedules", "schedule_id = ?", scheduleId); err != nil {
		return orderList, err
	}

	if err := FindRecords(&orders, "orders", "schedule_id = ?", scheduleId); err != nil {
		return orderList, err
	}
	fmt.Println(orders)

	var amrs = schedules[0].NumberOfAmrRequire

	for _, order := range orders {
		if (order.OrderStatus != string(orderStatus.Created)) && (order.OrderStatus != string(orderStatus.Canceled)) {
			return orderList, errors.New("Cancel failed, order started")
		} else if order.OrderStatus == string(orderStatus.Created) {
			amrs -= 1
		}
	}

	if amrs != 0 {
		return orderList, errors.New("Cancel failed, amr number not match")
	}

	updateMap := utils.CreateMap([]string{"order_status"}, orderStatus.Canceled)
	err := UpdateRecords(&updatedOrders, "orders", updateMap, "schedule_id = ?", scheduleId)
	if err != nil {
		return orderList, err
	}
	orderList, err = OrderListToOrderResponse(updatedOrders)
	if err != nil {
		return orderList, err
	}
	return orderList, nil
}

func OrderRequestToOrders(orderRequest dto.AddDeliveryOrderDTO, scheduleNo int, userId int, orderNo int) ([]db_models.Orders, error) {
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
		orders = append(orders, order)
	}
	return orders, nil
}

func OrderListToOrderResponse(orderList []db_models.Orders) (orderManagement.OrderList, error) {
	log.Printf("mysql query: OrderListToOrderResponse\n")
	var orderListResponse orderManagement.OrderList
	roomList, err := FindAllDutyRooms()
	if err != nil {
		return orderListResponse, err
	}
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
		orderListResponse[i].StartLocationName = roomList[orderListResponse[i].StartLocationID-1].LocationName
		orderListResponse[i].EndLocationName = roomList[orderListResponse[i].EndLocationID-1].LocationName
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

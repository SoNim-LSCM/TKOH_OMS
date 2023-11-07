package service

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/SoNim-LSCM/TKOH_OMS/database"
	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	dto "github.com/SoNim-LSCM/TKOH_OMS/models/DTO"
	"github.com/SoNim-LSCM/TKOH_OMS/models/orderManagement"
)

var ORDER_TYPE = []string{"PICK_AND_DELIVERY"}

func FindOrdersWithStatus(orderStatus []string) ([]db_models.Orders, error) {
	database.CheckDatabaseConnection()
	orders := []db_models.Orders{}
	log.Printf("mysql query: FindOrders: %s\n", orderStatus)
	if err := database.DB.Table("orders").Where("order_status IN ?", orderStatus).Find(&orders).Error; err != nil {
		return orders, err
	}
	if len(orders) == 0 {
		return orders, errors.New("Order Status not found")
	}
	return orders, nil
}

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
	if err := database.DB.Create(db_models.Schedules{ScheduleID: int(schedulesCount + 1), ScheduleStatus: "CREATED", ScheduleCraeteTime: timeNow}).Error; err != nil {
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

func OrderRequestToOrders(orderRequest dto.AddDeliveryOrderDTO, scheduleNo int, userId int, orderNo int) ([]db_models.Orders, error) {
	var orders []db_models.Orders
	log.Printf("mysql query: OrderRequestToOrders\n")
	for i := 0; i < orderRequest.NumberOfAmrRequire; i++ {
		var err error
		var order db_models.Orders
		order.ScheduleID = scheduleNo
		order.OrderID = i + orderNo
		order.OrderType = orderRequest.OrderType
		order.OrderCreatedType = "ADHOC"
		order.OrderCreatedBy = userId
		order.OrderStatus = "CREATED"
		order.OrderStartTime, err = StringToDatetime("")
		if err != nil {
			return orders, err
		}
		order.OrderEndTime, err = StringToDatetime("")
		if err != nil {
			return orders, err
		}
		order.ActualArrivalTime, err = StringToDatetime("")
		if err != nil {
			return orders, err
		}
		order.StartLocationID = orderRequest.StartLocationID
		order.EndLocationID = orderRequest.EndLocationID
		order.ExpectStartTime, err = StringToDatetime(orderRequest.ExpectingStartTime)
		if err != nil {
			return orders, err
		}
		order.ExpectDeliveryTime, err = StringToDatetime(orderRequest.ExpectingDeliveryTime)
		if err != nil {
			return orders, err
		}
		order.ExpectArrivalTime, err = StringToDatetime("")
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
		orderListResponse[i].EndTime, err = StringToResponseTime(order.OrderEndTime)
		if err != nil {
			return orderListResponse, err
		}
		orderListResponse[i].ActualArrivalTime, err = StringToResponseTime(order.ActualArrivalTime)
		if err != nil {
			return orderListResponse, err
		}
		orderListResponse[i].ExpectingStartTime, err = StringToResponseTime(order.ExpectStartTime)
		if err != nil {
			return orderListResponse, err
		}
		orderListResponse[i].ExpectingDeliveryTime, err = StringToResponseTime(order.ExpectDeliveryTime)
		if err != nil {
			return orderListResponse, err
		}
		orderListResponse[i].ExpectArrivalTime, err = StringToResponseTime(order.ExpectArrivalTime)
		if err != nil {
			return orderListResponse, err
		}
		orderListResponse[i].StartLocationName = roomList[orderListResponse[i].StartLocationID-1].LocationName
		orderListResponse[i].EndLocationName = roomList[orderListResponse[i].EndLocationID-1].LocationName
	}
	log.Println(orderListResponse)
	return orderListResponse, nil
}

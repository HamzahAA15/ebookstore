package order

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/model/request"
	"ebookstore/internal/model/response"
	"ebookstore/internal/repository"
	"ebookstore/internal/service"
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"
)

type orderService struct {
	orderRepository repository.IOrderRepository
	bookRepository  repository.IBookRepository
}

func NewOrderService(orderRepository repository.IOrderRepository, bookRepository repository.IBookRepository) service.IOrderService {
	return &orderService{
		orderRepository: orderRepository,
		bookRepository:  bookRepository,
	}
}

func (o *orderService) GetUserOrders(ctx context.Context) ([]response.OrderData, error) {
	customerID := ctx.Value("id").(uint)
	mapOrderIDToOrderData := make(map[uint]response.OrderData)
	resp := []response.OrderData{}

	//get all order from customer
	orders, err := o.orderRepository.GetOrderHistoryByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order history: %s", err.Error())
	}

	//assign order data into response
	for _, order := range orders {
		mapOrderIDToOrderData[order.ID] = response.OrderData{
			OrderID:           order.ID,
			CustomerReference: order.CustomerReference,
			ReceiverName:      order.ReceiverName,
			Address:           order.Address,
			City:              order.City,
			District:          order.District,
			PostalCode:        order.PostalCode,
			Shipper:           order.Shipper,
			AirwaybillNumber:  order.AirwaybillNumber,
			OrderDate:         order.OrderDate,
			TotalPrice:        order.TotalPrice,
			TotalItem:         order.TotalItem,
		}
	}

	//assign order items into response
	for orderID, order := range mapOrderIDToOrderData {
		items, err := o.orderRepository.GetItemsByOrderID(ctx, orderID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items: %s", err.Error())
		}

		for _, item := range items {
			book, err := o.bookRepository.GetBookByID(ctx, item.BookID)
			if err != nil {
				return nil, fmt.Errorf("failed to get book: %s", err.Error())
			}

			order.Items = append(order.Items, response.Item{
				BookID:   item.BookID,
				Title:    book.Title,
				Author:   book.Author,
				Quantity: item.Quantity,
				Price:    book.Price,
			})
		}

		resp = append(resp, order)
	}

	return resp, nil
}

func (o *orderService) CreateOrder(ctx context.Context, req request.CreateOrder) (response.Order, error) {
	var totalPrice float64
	var totalQuantity int

	//insert data to order
	var order model.Order
	customerID := ctx.Value("id").(uint)

	if req.ReceiverName == "" {
		userName := ctx.Value("username").(string)
		req.ReceiverName = userName
	}

	order.ReceiverName = req.ReceiverName
	order.CustomerID = customerID
	order.OrderDate = time.Now().UTC()
	order.CustomerReference = generateCustomerReference(order.OrderDate)
	order.Address = req.Address
	order.City = req.City
	order.District = req.District
	order.PostalCode = req.PostalCode
	order.Shipper = req.Shipper
	order.AirwaybillNumber = generateAirwaybillNumber(req.Shipper)
	orderID, err := o.orderRepository.CreateOrder(ctx, order)
	if err != nil {
		return response.Order{}, err
	}

	for _, item := range req.Items {
		totalPrice += item.Price * float64(item.Quantity)
		totalQuantity += item.Quantity

		err = o.orderRepository.CreateItem(ctx, model.Item{
			BookID:    item.BookID,
			Quantity:  item.Quantity,
			OrderID:   orderID,
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			return response.Order{}, err
		}
	}

	if err != nil {
		return response.Order{}, err
	}

	var orderUpdate model.Order
	orderUpdate.ID = orderID
	orderUpdate.TotalItem = totalQuantity
	orderUpdate.TotalPrice = totalPrice
	orderUpdate.UpdatedAt = pq.NullTime{Time: time.Now().UTC(), Valid: true}

	err = o.orderRepository.UpdateOrderByOrderID(ctx, orderUpdate)
	if err != nil {
		return response.Order{}, err
	}

	data := response.CreateOrderData{
		OrderID:           orderID,
		CustomerReference: order.CustomerReference,
		AirwaybillNumber:  order.AirwaybillNumber,
		OrderDate:         order.OrderDate,
	}

	return response.Order{
		Data: data,
	}, nil
}

func generateCustomerReference(orderDate time.Time) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(orderDate.UnixNano()))
	result := make([]byte, 8)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

func generateAirwaybillNumber(shipper string) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, 10)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return shipper + "-" + string(result)
}

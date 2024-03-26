package order

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/model/request"
	"ebookstore/internal/model/response"
	"ebookstore/internal/repository"
	"ebookstore/internal/service"
	"ebookstore/utils/transactioner"
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"
)

var mapBook = make(map[uint]float64)

type orderService struct {
	orderRepository     repository.IOrderRepository
	TransactionProvider transactioner.TransactionProvider
	bookRepository      repository.IBookRepository
}

func NewOrderService(orderRepository repository.IOrderRepository, bookRepository repository.IBookRepository, tx transactioner.TransactionProvider) service.IOrderService {
	return &orderService{
		orderRepository:     orderRepository,
		bookRepository:      bookRepository,
		TransactionProvider: tx,
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

	tx, err := o.TransactionProvider.NewTransaction(ctx)
	if err != nil {
		return response.Order{}, fmt.Errorf("failed to start transaction: %s", err.Error())
	}
	defer tx.Rollback()

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

	//create order
	orderID, err := o.orderRepository.CreateOrder(ctx, tx, order)
	if err != nil {
		return response.Order{}, fmt.Errorf("failed to create order: %s", err.Error())
	}

	for _, item := range req.Items {
		//get book price from inMemory Cache first
		bookPrice, ok := mapBook[item.BookID]
		if !ok {
			book, err := o.bookRepository.GetBookByID(ctx, item.BookID)
			if err != nil {
				return response.Order{}, fmt.Errorf("failed to get book: %s", err.Error())
			}

			bookPrice = book.Price
			mapBook[item.BookID] = bookPrice
		}

		//set total price & quantity
		totalPrice += bookPrice * float64(item.Quantity)
		totalQuantity += item.Quantity

		//create item
		err = o.orderRepository.CreateItem(ctx, tx, model.Item{
			BookID:    item.BookID,
			Quantity:  item.Quantity,
			OrderID:   orderID,
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			return response.Order{}, fmt.Errorf("failed to create item: %s", err.Error())
		}
	}

	var orderUpdate model.Order
	orderUpdate.ID = orderID
	orderUpdate.TotalItem = totalQuantity
	orderUpdate.TotalPrice = totalPrice
	orderUpdate.UpdatedAt = pq.NullTime{Time: time.Now().UTC(), Valid: true}

	//update order
	err = o.orderRepository.UpdateOrderByOrderID(ctx, tx, orderUpdate)
	if err != nil {
		return response.Order{}, fmt.Errorf("failed to update order: %s", err.Error())
	}

	data := response.CreateOrderData{
		OrderID:           orderID,
		CustomerReference: order.CustomerReference,
		AirwaybillNumber:  order.AirwaybillNumber,
		OrderDate:         order.OrderDate.String(),
	}

	err = tx.Commit()
	if err != nil {
		return response.Order{}, fmt.Errorf("failed to commit transaction: %s", err.Error())
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

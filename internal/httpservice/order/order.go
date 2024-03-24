package httpservice

import (
	"ebookstore/internal/service/book"
	"ebookstore/internal/service/customer"
	"ebookstore/internal/service/order"
	"ebookstore/utils/request"
	"ebookstore/utils/response"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	customerService customer.ICustomerService
	bookService     book.IBookService
	orderService    order.IOrderService
}

func NewOrderHandler(customerService customer.ICustomerService, bookService book.IBookService, orderService order.IOrderService) *OrderHandler {
	return &OrderHandler{
		customerService: customerService,
		bookService:     bookService,
		orderService:    orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	req := request.CreateOrder{}

	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Order{
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	err = isValidCreateOrderReq(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Order{
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	fmt.Printf("%+v", req)
	resp, err := h.orderService.CreateOrder(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Order{
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	customer := c.Locals("username").(string)
	return c.Status(fiber.StatusOK).JSON(response.Order{
		StatusCode: fiber.StatusOK,
		Message:    fmt.Sprintf("successfully created order for %s", customer),
		TotalItem:  resp.TotalItem,
		TotalPrice: resp.TotalPrice,
	})
}

func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	orders, err := h.orderService.GetUserOrders(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.GetUserOrders{
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.GetUserOrders{
		StatusCode: fiber.StatusOK,
		Message:    "success",
		Data:       orders,
	})
}

func isValidCreateOrderReq(req request.CreateOrder) error {
	if len(req.Items) == 0 {
		return errors.New("items cannot be empty")
	}

	if req.ReceiverName == "" {
		return errors.New("receiver name cannot be empty")
	}

	if req.Address == "" {
		return errors.New("receiver address cannot be empty")
	}

	if req.City == "" {
		return errors.New("receiver city cannot be empty")
	}

	if req.District == "" {
		return errors.New("receiver district cannot be empty")
	}

	if req.PostalCode == "" {
		return errors.New("receiver postal code cannot be empty")
	}

	return nil
}

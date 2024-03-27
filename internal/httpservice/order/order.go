package order

import (
	"ebookstore/internal/model/request"
	"ebookstore/internal/model/response"
	"ebookstore/internal/service"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderService service.IOrderService
}

func NewOrderHandler(orderService service.IOrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
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

	ctx := c.Context()
	data, err := h.orderService.CreateOrder(ctx, req)
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
		Data:       data,
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

	for _, item := range req.Items {
		if item.BookID == 0 {
			return errors.New("book id cannot be empty")
		}

		if item.Quantity == 0 {
			return errors.New("quantity cannot be empty")
		}
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

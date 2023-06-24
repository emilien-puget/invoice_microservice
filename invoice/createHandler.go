package invoice

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/emilien-puget/invoice_microservice/money"
	"github.com/emilien-puget/invoice_microservice/user"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CreateInvoiceHandler struct {
	invoiceRepository interface {
		Create(ctx context.Context, invoice Invoice) (int64, error)
	}
	userRepository interface {
		GetById(ctx context.Context, id int64) (*user.User, error)
	}
	validator *validator.Validate
}

func NewCreateInvoiceHandler(validate *validator.Validate, repository *Repository, userRepository *user.Repository) *CreateInvoiceHandler {
	return &CreateInvoiceHandler{validator: validate, invoiceRepository: repository, userRepository: userRepository}
}

type createInvoicePayload struct {
	UserID int64   `json:"user_id" validate:"required"`
	Amount float64 `json:"amount" validate:"required"`
	Label  string  `json:"label" validate:"required"`
}

func (h CreateInvoiceHandler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse the request body into an createInvoicePayload struct
	payload := new(createInvoicePayload)
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := h.validator.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := h.userRepository.GetById(ctx, payload.UserID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "user not found")
		}
		return fmt.Errorf("userRepository.GetById: %w", err)
	}

	// Create the invoice in the repository
	invoice := Invoice{
		UserID: payload.UserID,
		Amount: money.NewMoneyFromFloat(payload.Amount),
		Label:  payload.Label,
	}

	invoiceID, err := h.invoiceRepository.Create(ctx, invoice)
	if err != nil {
		return fmt.Errorf("invoiceRepository.Create: %w", err)
	}

	return c.JSON(http.StatusCreated, map[string]int64{"invoice_id": invoiceID})
}

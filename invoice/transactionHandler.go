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

type DoTransactionHandler struct {
	invoiceRepository interface {
		GetByID(ctx context.Context, id int64) (*Invoice, error)
		MarkAsPaid(ctx context.Context, id int64) error
	}
	userRepository interface {
		ModifyBalance(ctx context.Context, userID int64, amount money.Money) error
	}
	validator *validator.Validate
}

func NewDoTransactionHandler(invoiceRepository *Repository, userRepository *user.Repository, validate *validator.Validate) *DoTransactionHandler {
	return &DoTransactionHandler{invoiceRepository: invoiceRepository, userRepository: userRepository, validator: validate}
}

type TransactionPayload struct {
	InvoiceID int64   `json:"invoice_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required"`
	Reference string  `json:"reference" validate:"required"`
}

func (d DoTransactionHandler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse and validate the JSON payload
	var payload TransactionPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := d.validator.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Fetch the invoice by ID
	invoice, err := d.invoiceRepository.GetByID(ctx, payload.InvoiceID)
	if err != nil {
		if errors.Is(err, ErrInvoiceNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "invoice not found")
		}
		return fmt.Errorf("invoiceRepository.GetByID: %w", err)
	}

	if invoice.Amount != money.NewMoneyFromFloat(payload.Amount) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid amount")
	}
	if invoice.Status == "paid" {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "invoice already paid")
	}

	err = d.userRepository.ModifyBalance(ctx, invoice.UserID, invoice.Amount)
	if err != nil {
		return fmt.Errorf("userRepository.ModifyBalance: %w", err)
	}

	err = d.invoiceRepository.MarkAsPaid(ctx, invoice.ID)
	if err != nil {
		return fmt.Errorf("invoiceRepository.MarkAsPaid: %w", err)
	}

	return c.NoContent(http.StatusNoContent)
}

package invoice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/emilien-puget/invoice_microservice/money"
)

type Invoice struct {
	ID     int64
	UserID int64
	Status string
	Label  string
	Amount money.Money
}

type Repository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, invoice Invoice) (int64, error) {
	query := `
		INSERT INTO jump.public.invoices (user_id, label, amount)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var invoiceId int64
	row := stmt.QueryRowContext(ctx, invoice.UserID, invoice.Label, invoice.Amount)
	if err := row.Scan(&invoiceId); err != nil {
		return 0, fmt.Errorf("failed to create invoice: %w", err)
	}

	return invoiceId, nil
}

var ErrInvoiceNotFound = errors.New("invoice not found")

func (r *Repository) MarkAsPaid(ctx context.Context, id int64) error {
	query := `
		UPDATE jump.public.invoices
		SET status = 'paid'
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark invoice as paid: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrInvoiceNotFound
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Invoice, error) {
	query := `
		SELECT id, user_id, status, label, amount
		FROM jump.public.invoices
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var invoice Invoice
	if err := row.Scan(&invoice.ID, &invoice.UserID, &invoice.Status, &invoice.Label, &invoice.Amount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvoiceNotFound
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	return &invoice, nil
}

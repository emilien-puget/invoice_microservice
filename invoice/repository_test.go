package invoice

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceRepository_Create(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	repo := NewInvoiceRepository(db)
	ctx := context.Background()

	// Create a test invoice
	invoice := Invoice{
		UserID: 1,
		Label:  "Test Invoice",
		Amount: 1000,
	}

	// Mock the expected query and result
	mock.ExpectPrepare("INSERT INTO jump.public.invoices (user_id, label, amount) VALUES ($1, $2, $3) RETURNING id").
		ExpectQuery().
		WithArgs(invoice.UserID, invoice.Label, invoice.Amount).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Call the Create method
	invoiceID, err := repo.Create(ctx, invoice)
	require.NoError(t, err)
	assert.Equal(t, int64(1), invoiceID)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestInvoiceRepository_GetByID(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	repo := NewInvoiceRepository(db)
	ctx := context.Background()

	// Create a test invoice
	invoice := &Invoice{
		ID:     1,
		UserID: 1,
		Status: "pending",
		Label:  "Test Invoice",
		Amount: 1000,
	}

	// Mock the expected query and result
	mock.ExpectQuery("SELECT id, user_id, status, label, amount FROM jump.public.invoices WHERE id = ?").
		WithArgs(invoice.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "status", "label", "amount"}).
			AddRow(invoice.ID, invoice.UserID, invoice.Status, invoice.Label, invoice.Amount))

	// Call the GetByID method
	result, err := repo.GetByID(ctx, invoice.ID)
	require.NoError(t, err)
	assert.Equal(t, invoice, result)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

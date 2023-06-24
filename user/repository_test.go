package user

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/emilien-puget/invoice_microservice/money"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllUsers(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}
	defer db.Close()

	// Create a new UserRepository with the mock database connection
	repo := NewUserRepository(db)

	// Define the expected query and rows for the GetAll() method
	query := "SELECT id, first_name, last_name, balance FROM jump.public.users"
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "balance"}).
		AddRow(1, "John", "Doe", 1000).
		AddRow(2, "Jane", "Smith", 2000)

	// Expect the query to be executed and return the mocked rows
	mock.ExpectQuery(query).WillReturnRows(rows)

	// Call the GetAll() method
	users, err := repo.GetAll(context.Background())
	require.NoError(t, err)

	// Assert the expected number of users
	assert.Len(t, users, 2)

	// Assert the properties of the first user
	assert.Equal(t, int64(1), users[0].ID)
	assert.Equal(t, "John", users[0].FirstName)
	assert.Equal(t, "Doe", users[0].LastName)
	assert.Equal(t, money.Money(1000), users[0].Balance)

	// Assert the properties of the second user
	assert.Equal(t, int64(2), users[1].ID)
	assert.Equal(t, "Jane", users[1].FirstName)
	assert.Equal(t, "Smith", users[1].LastName)
	assert.Equal(t, money.Money(2000), users[1].Balance)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetUserById(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}
	defer db.Close()

	// Create a new UserRepository with the mock database connection
	repo := NewUserRepository(db)

	// Define the expected query, arguments, and rows for the GetById() method
	query := "SELECT id, first_name, last_name, balance FROM jump.public.users WHERE id = $1"
	args := []driver.Value{1}
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "balance"}).AddRow(1, "John", "Doe", 1000)

	// Expect the query to be executed and return the mocked rows
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WithArgs(args...).WillReturnRows(rows)

	// Call the GetById() method
	user, err := repo.GetById(context.Background(), 1)
	require.NoError(t, err)

	// Assert the properties of the user
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, money.Money(1000), user.Balance)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateUser(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}
	defer db.Close()

	// Create a new UserRepository with the mock database connection
	repo := NewUserRepository(db)

	// Define the expected query, arguments, and rows for the Update() method
	query := "UPDATE jump.public.users SET first_name = ?, last_name = ?, balance = ? WHERE id = ?"
	args := []driver.Value{"John", "Doe", int64(2000), 1}

	// Expect the query to be executed
	mock.ExpectPrepare(query)
	mock.ExpectExec(query).WithArgs(args...).WillReturnResult(sqlmock.NewResult(0, 1))

	// Create a new user
	user := &User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Balance:   2000,
	}

	// Call the Update() method
	err = repo.Update(context.Background(), user)
	assert.NoError(t, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

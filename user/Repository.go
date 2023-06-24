package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/emilien-puget/invoice_microservice/money"
)

type Repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Balance   money.Money
}

func (r *Repository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE jump.public.users
		SET first_name = ?, last_name = ?, balance = ?
		WHERE id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.FirstName, user.LastName, user.Balance, user.ID)
	if err != nil {
		return err
	}

	return nil
}

var ErrUserNotFound = errors.New("user not found")

func (r *Repository) GetById(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, balance
		FROM jump.public.users
		WHERE id = $1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	user := &User{}
	err = stmt.QueryRow(id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*User, error) {
	query := `
		SELECT id, first_name, last_name, balance
		FROM jump.public.users
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Balance)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

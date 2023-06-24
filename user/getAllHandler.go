package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetAllHandler struct {
	userRepository interface {
		GetAll(ctx context.Context) ([]*User, error)
	}
}

func NewGetAllHandler(userRepository *Repository) *GetAllHandler {
	return &GetAllHandler{userRepository: userRepository}
}

type GetUsersHandlerResponse struct {
	UserID    int64   `json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balance   float64 `json:"balance"`
}

func (g GetAllHandler) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := g.userRepository.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("userRepository.GetAll: %w", err)
	}

	results := make([]GetUsersHandlerResponse, len(users))
	for _, user := range users {
		results = append(results, GetUsersHandlerResponse{
			UserID:    user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Balance:   user.Balance.ToFloat(),
		})
	}

	return c.JSON(http.StatusOK, results)
}

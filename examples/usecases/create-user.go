package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"

	"github.com/wesley601/symphony/slogutils"

	"github.com/google/uuid"
)

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

type CreateUserUseCase struct {
	conn *sql.DB
}

func NewCreateUserUseCase(conn *sql.DB) *CreateUserUseCase {
	return &CreateUserUseCase{conn: conn}
}

func (c *CreateUserUseCase) Handle(ctx context.Context, event []byte) ([]byte, error) {
	slog.Info("start to create a user")
	u := new(User)
	if err := json.Unmarshal(event, u); err != nil {
		slog.Error("Error unmarshaling user:", slogutils.Error(err))
		return nil, err
	}
	u.ID = uuid.New().String()
	_, err := c.conn.Exec("INSERT INTO users (id, name, email, active) VALUES ($1, $2, $3, $4)", u.ID, u.Name, u.Email, u.Active)
	if err != nil {
		slog.Error("Error creating user:", slogutils.Error(err))
		return nil, err
	}
	slog.Info("User created successfully", slog.String("name", u.Name), slog.String("email", u.Email))
	return json.Marshal(u)
}

package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"

	"github.com/wesley601/symphony/slogutils"
)

type ActivateAccountUseCase struct {
	conn *sql.DB
}

func NewActivateAccountUseCase(conn *sql.DB) *ActivateAccountUseCase {
	return &ActivateAccountUseCase{conn: conn}
}

func (c *ActivateAccountUseCase) Handle(ctx context.Context, event []byte) ([]byte, error) {
	slog.Info("start to activate an account")
	u := new(User)
	if err := json.Unmarshal(event, u); err != nil {
		slog.Error("Error unmarshaling account:", slogutils.Error(err))
		return nil, err
	}
	u.Active = true
	_, err := c.conn.Exec("UPDATE users SET active=$1 WHERE id=$2", u.Active, u.ID)
	if err != nil {
		slog.Error("Error activating account:", slogutils.Error(err))
		return nil, err
	}

	slog.Info("Account activated successfully", slog.String("id", u.ID), slog.Bool("active", u.Active))
	return nil, nil
}

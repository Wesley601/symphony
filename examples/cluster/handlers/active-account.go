package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
)

type ActivateAccountHandler struct {
	conn *sql.DB
}

func NewActivateAccountHandler(conn *sql.DB) *ActivateAccountHandler {
	return &ActivateAccountHandler{conn: conn}
}

func (c *ActivateAccountHandler) Handle(ctx context.Context, event []byte) ([]byte, error) {
	slog.Info("start to activate an account")
	u := new(User)
	if err := json.Unmarshal(event, u); err != nil {
		slog.Error("Error unmarshaling account: " + err.Error())
		return nil, err
	}
	u.Active = true
	_, err := c.conn.Exec("UPDATE users SET active=$1 WHERE id=$2", u.Active, u.ID)
	if err != nil {
		slog.Error("Error activating account: " + err.Error())
		return nil, err
	}

	slog.Info("Account activated successfully", slog.String("id", u.ID), slog.Bool("active", u.Active))
	return nil, nil
}

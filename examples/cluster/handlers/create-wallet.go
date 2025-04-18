package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
)

type CreateWalletHandler struct {
	conn *sql.DB
}

func NewCreateWalletHandler(conn *sql.DB) *CreateWalletHandler {
	return &CreateWalletHandler{conn: conn}
}

func (c *CreateWalletHandler) Handle(ctx context.Context, event []byte) ([]byte, error) {
	slog.Info("start to create a wallet")
	u := new(User)
	if err := json.Unmarshal(event, u); err != nil {
		slog.Error("Error unmarshaling user: " + err.Error())
		return nil, err
	}
	_, err := c.conn.Exec("INSERT INTO wallets (id, user_id) VALUES ($1, $2)", uuid.New().String(), u.ID)
	if err != nil {
		slog.Error("Error creating wallet: " + err.Error())
		return nil, err
	}
	slog.Info("Wallet created successfully", slog.String("user_id", u.ID))
	return json.Marshal(u)
}

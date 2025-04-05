package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"

	"github.com/wesley601/symphony/slogutils"

	"github.com/google/uuid"
)

type CreateWalletUseCase struct {
	conn *sql.DB
}

func NewCreateWalletUseCase(conn *sql.DB) *CreateWalletUseCase {
	return &CreateWalletUseCase{conn: conn}
}

func (c *CreateWalletUseCase) Handle(ctx context.Context, event []byte) ([]byte, error) {
	slog.Info("start to create a wallet")
	u := new(User)
	if err := json.Unmarshal(event, u); err != nil {
		slog.Error("Error unmarshaling user:", slogutils.Error(err))
		return nil, err
	}
	u.ID = uuid.New().String()
	_, err := c.conn.Exec("INSERT INTO wallets (id, user_id) VALUES (?, ?)", u.ID, u.ID)
	if err != nil {
		slog.Error("Error creating wallet:", slogutils.Error(err))
		return nil, err
	}
	slog.Info("Wallet created successfully", slog.String("user_id", u.ID))
	return json.Marshal(u)
}

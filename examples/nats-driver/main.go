package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/wesley601/symphony"
	"github.com/wesley601/symphony/drivers"
	"github.com/wesley601/symphony/examples/usecases"
	"github.com/wesley601/symphony/slogutils"

	"github.com/nats-io/nats.go"
)

func main() {
	driver, err := drivers.NewNatsDriver(nats.DefaultURL, "group-id")
	if err != nil {
		slog.Error("Error creating NATS driver:", slogutils.Error(err))
		return
	}
	defer driver.Close()
	symphony := symphony.New(driver)

	createUser := new(usecases.CreateUserUseCase)
	createWallet := new(usecases.CreateWalletUseCase)
	activateAccount := new(usecases.ActivateAccountUseCase)

	finishSignal := make(chan os.Signal, 1)
	signal.Notify(finishSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if err := symphony.
		After("create.user", createUser).
		After("create.wallet", createWallet).
		After("activate.account", activateAccount).
		Play(context.Background()); err != nil {
		slog.Error("Error playing symphony:", slogutils.Error(err))
	}

	slog.Info("Symphony played successfully")

	if err := driver.Publish("create.user", []byte(`{"name": "John Doe", "email": "john.doe@example.com"}`)); err != nil {
		slog.Error("Error publishing create.user event:", slogutils.Error(err))
	}

	<-finishSignal
}

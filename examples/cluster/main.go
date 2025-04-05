package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wesley601/symphony"
	"github.com/wesley601/symphony/drivers"
	"github.com/wesley601/symphony/examples/usecases"
	"github.com/wesley601/symphony/slogutils"

	_ "github.com/lib/pq"
)

func main() {
	conn, err := migrateAndGetConnection()
	if err != nil {
		slog.Error("Error migrating database:", slogutils.Error(err))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	driver, err := drivers.NewNatsDriver(os.Getenv("NATS_URL"), "group-id")
	if err != nil {
		slog.Error("Error creating NATS driver:", slogutils.Error(err))
		return
	}
	defer driver.Close()
	symphony := symphony.New(driver)

	createUser := usecases.NewCreateUserUseCase(conn)
	createWallet := usecases.NewCreateWalletUseCase(conn)
	activateAccount := usecases.NewActivateAccountUseCase(conn)

	finishSignal := make(chan os.Signal, 1)
	signal.Notify(finishSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if err := symphony.
		After("create.user", createUser).
		After("create.wallet", createWallet).
		After("activate.account", activateAccount).
		Play(ctx); err != nil {
		slog.Error("Error playing symphony:", slogutils.Error(err))
	}

	slog.Info("Symphony played successfully")

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		if err := driver.Publish("create.user", []byte(`{"name": "John Doe", "email": "john.doe@example.com"}`)); err != nil {
			slog.Error("Error publishing create.user event:", slogutils.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "User creation initiated"}`))
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		slog.Info("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error:", slogutils.Error(err))
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error:", slogutils.Error(err))
	} else {
		slog.Info("HTTP server shutdown gracefully")
	}

	cancel()

	// Give a short grace period for clean symphony shutdown
	time.Sleep(2 * time.Second)

	slog.Info("Application shutdown complete")
}

func migrateAndGetConnection() (conn *sql.DB, err error) {
	conn, err = sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	if _, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id uuid PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			active BOOLEAN NOT NULL DEFAULT FALSE
		)`); err != nil {
		return nil, err
	}

	if _, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS wallets (
			id uuid PRIMARY KEY,
			user_id uuid NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`); err != nil {
		return nil, err
	}

	return conn, nil
}

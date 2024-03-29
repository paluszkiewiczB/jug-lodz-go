// Package main is an entrypoint to fullstack exposing a list of items to be done.
// It uses no dependencies aside from SQL connector implementation, since a standard library provides just the interface (like JDBC in Java).
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"todo/internal/data"
	"todo/internal/server"
	"todo/internal/todo"
)

func main() {
	ctx := gracefulShutdown()

	storage, err := data.NewSQLiteTaskStorage("./todos.db")
	if err != nil {
		slog.Error("failed to create SQLite storage", slog.String("err", err.Error()))
		return
	}

	err = storage.Initialize()
	if err != nil {
		slog.Error("failed to initialize SQLite storage", slog.String("err", err.Error()))
		return
	}

	handler := todo.NewHandler(storage)
	s, err := server.NewHttp(nil, handler, storage)
	if err != nil {
		slog.ErrorContext(ctx, "creating the server: %w", err)
		return
	}

	err = s.Start(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to start a server", slog.String("err", err.Error()))
		return
	}
}

// listens for SIGINT and SIGTERM and cancels context if received
func gracefulShutdown() context.Context {
	ctx, cancel := context.WithCancelCause(context.Background())

	go func() {
		c := make(chan os.Signal, 3)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		sig := <-c
		cancel(fmt.Errorf("received shutdown signal: %s", sig))
	}()

	return ctx
}

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanaru-ssk/go-http-server/entity/task"
	"github.com/kanaru-ssk/go-http-server/interface/inbound/http/handler"
	"github.com/kanaru-ssk/go-http-server/interface/outbound/postgres"
	postgrestask "github.com/kanaru-ssk/go-http-server/interface/outbound/postgres/task"
	"github.com/kanaru-ssk/go-http-server/lib/id"
	"github.com/kanaru-ssk/go-http-server/lib/tx"
	"github.com/kanaru-ssk/go-http-server/usecase"
)

func main() {
	// OSシグナルに反応してHTTPサーバーをGraceful Shutdownさせるコンテキストを用意
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	idGenerator := &id.SecureGenerator{}
	pool, err := postgres.NewPool(ctx, postgres.Config{
		Host:     "db",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Database: "postgres",
		MaxConns: 10,
	})
	if err != nil {
		slog.ErrorContext(ctx, "main.main: postgres.NewPool", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	txManager := postgres.NewManager(pool)
	app := dependencyInjection(idGenerator, txManager, pool)

	addr := ":8000"
	srv := &http.Server{
		Addr:    addr,
		Handler: app.Handler,
	}

	go func() {
		slog.InfoContext(ctx, "main.main: starting http server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "main.main: http.Server.ListenAndServe", "err", err)
		}
	}()

	<-ctx.Done()
	slog.InfoContext(context.Background(), "main.main: shutdown signal received")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(context.Background(), "main.main: http.Server.Shutdown", "err", err)
	}
}

type Application struct {
	Handler http.Handler
}

func dependencyInjection(idGenerator id.Generator, txManager tx.Manager, pool *pgxpool.Pool) Application {
	// interface/outbound
	taskRepository := postgrestask.NewRepository(pool)

	// entity
	taskFactory := task.NewFactory(idGenerator)

	// usecase
	taskUseCase := usecase.NewTaskUseCase(txManager, taskFactory, taskRepository)

	// interface/inbound
	taskHandler := handler.NewTaskHandler(taskUseCase)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handler.HandleGetHealthz)

	mux.HandleFunc("GET /core/v1/task/get", taskHandler.HandleGetV1)
	mux.HandleFunc("GET /core/v1/task/list", taskHandler.HandleListV1)
	mux.HandleFunc("POST /core/v1/task/create", taskHandler.HandleCreateV1)
	mux.HandleFunc("PUT /core/v1/task/update", taskHandler.HandleUpdateV1)
	mux.HandleFunc("DELETE /core/v1/task/delete", taskHandler.HandleDeleteV1)
	mux.HandleFunc("PUT /core/v1/task/done", taskHandler.HandleDoneV1)

	return Application{Handler: mux}
}

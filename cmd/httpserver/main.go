package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kanaru-ssk/go-http-server/entity/task"
	"github.com/kanaru-ssk/go-http-server/interface/inbound/http/handler"
	"github.com/kanaru-ssk/go-http-server/interface/outbound/memory"
	memorytask "github.com/kanaru-ssk/go-http-server/interface/outbound/memory/task"
	"github.com/kanaru-ssk/go-http-server/lib/id"
	"github.com/kanaru-ssk/go-http-server/lib/tx"
	"github.com/kanaru-ssk/go-http-server/usecase"
)

func main() {
	// OSシグナルに反応してHTTPサーバーをGraceful Shutdownさせるコンテキストを用意
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	idGenerator := &id.SecureGenerator{}
	mu := &sync.RWMutex{}
	txManager := memory.NewTxManager()
	tasks := make(map[string]*task.Task)
	app := dependencyInjection(mu, idGenerator, txManager, tasks)

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

func dependencyInjection(mu *sync.RWMutex, idGenerator id.Generator, txManager tx.Manager, tasks map[string]*task.Task) Application {
	// interface/outbound
	taskRepository := memorytask.NewRepository(mu, tasks)

	// entity
	taskFactory := task.NewFactory(idGenerator)

	// usecase
	taskUseCase := usecase.NewTaskUseCase(txManager, taskFactory, taskRepository)

	// interface/inbound
	taskHandler := handler.NewTaskHandler(taskUseCase)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handler.HandleGetHealthz)

	mux.HandleFunc("POST /core/v1/task/get", taskHandler.HandleGetV1)
	mux.HandleFunc("POST /core/v1/task/list", taskHandler.HandleListV1)
	mux.HandleFunc("POST /core/v1/task/create", taskHandler.HandleCreateV1)
	mux.HandleFunc("POST /core/v1/task/update", taskHandler.HandleUpdateV1)
	mux.HandleFunc("POST /core/v1/task/delete", taskHandler.HandleDeleteV1)
	mux.HandleFunc("POST /core/v1/task/done", taskHandler.HandleDoneV1)

	return Application{Handler: mux}
}

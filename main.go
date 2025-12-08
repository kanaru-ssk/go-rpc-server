package main

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kanaru-ssk/go-rpc-server/domain/task"
	"github.com/kanaru-ssk/go-rpc-server/interface/inbound/http/handler"
	"github.com/kanaru-ssk/go-rpc-server/interface/outbound/memory"
	memorytask "github.com/kanaru-ssk/go-rpc-server/interface/outbound/memory/task"
	"github.com/kanaru-ssk/go-rpc-server/lib/id"
	"github.com/kanaru-ssk/go-rpc-server/lib/tx"
	"github.com/kanaru-ssk/go-rpc-server/usecase"
)

func main() {
	// OSシグナルに反応してHTTPサーバーをGraceful Shutdownさせるコンテキストを用意
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	idGenerator := &id.SecureGenerator{}
	mu := &sync.RWMutex{}
	txManager := memory.NewTxManager(mu)
	tasks := make(map[string]*task.Task)
	app := di(idGenerator, txManager, tasks)

	addr := ":8000"
	go func() {
		slog.InfoContext(ctx, "main.main: starting http server on: ", "addr", addr)
		if err := http.ListenAndServe(addr, app.Handler); err != nil {
			slog.WarnContext(ctx, "main.main: http.ListenAndServe: ", "err", err)
		}
	}()

	<-ctx.Done()
}

type Application struct {
	Handler http.Handler
}

func di(idGenerator id.Generator, txManager tx.Manager, tasks map[string]*task.Task) Application {
	mu := &sync.RWMutex{}

	// outbound interface
	taskRepository := memorytask.NewRepository(mu, tasks)

	// domain
	taskFactory := task.NewFactory(idGenerator)

	// usecase
	userUsecase := usecase.NewTaskUsecase(txManager, taskFactory, taskRepository)

	// inbound interface
	taskHandler := handler.NewTaskHandler(userUsecase)

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

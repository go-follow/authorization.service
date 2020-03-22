package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-follow/authorization.service/pkg/logger"
)

//Options - настройки для http сервера
type Options struct {
	//Port - порт, для запуска сервера
	Port uint16
}

//ListenServer - запуск сервера, для прослушивания
func ListenServer(r chi.Router, o *Options) {
	if err := checkOptions(o); err != nil {
		logger.Fatal("invalid options for initialization server")
	}
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	logger.Infof("rest server run in port: %d", o.Port)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", o.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	ctx := context.Background()
	go gracefullShutdown(ctx, httpServer, quit, done)

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("не удалось запустить рест сервер на порту: %d", o.Port)
	}

	<-done
	logger.Info("рест сервер дождался выполнения всех задач клиентов и завершил свою работу")
}

func checkOptions(o *Options) error {
	if o == nil {
		return fmt.Errorf("empty options")
	}
	if o.Port == 0 {
		return fmt.Errorf("port should be more 0")
	}
	return nil
}

func gracefullShutdown(ctx context.Context, server *http.Server,
	quit <-chan os.Signal, done chan<- bool) {
	<-quit
	logger.Info("Сервер начинает процедуру остановки и дожидается всех клиентов...")

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Could not gracefully shutdown the server: ", err)
	}
	close(done)
}
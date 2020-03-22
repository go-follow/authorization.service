package apiserver

import (
	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi"
	"github.com/go-follow/authorization.service/config"
	"github.com/go-follow/authorization.service/internal/service"
	"github.com/go-follow/authorization.service/internal/service/store/sqlstore"
	"github.com/go-follow/authorization.service/internal/service/transport"
	"github.com/go-follow/authorization.service/pkg/db"
	"github.com/go-follow/authorization.service/pkg/logger"
	"github.com/go-follow/authorization.service/pkg/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//Start - function start server
func Start(cfg *config.Config) {
	db, err := db.Connect(cfg.Db)
	if err != nil {
		logger.Fatalf("не удалось соединиться с БД %v", err)
	}
	defer db.Close()

	s := sqlstore.New(db)

	router := chi.NewRouter()
	metrics := chiprometheus.NewMiddleware("komiac-authorization-service")
	router.Use(metrics)
	router.Handle("/metrics", promhttp.Handler())
	svc := service.New(s, cfg)

	router.Route("/api/v1", func(rapi chi.Router) {
		transport.NewHTTP(svc, rapi)
	})
	o := &server.Options{Port: cfg.Server.Port}
	server.ListenServer(router, o)
}

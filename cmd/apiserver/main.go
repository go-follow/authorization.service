package main

import (
	"github.com/go-follow/authorization.service/config"
	"github.com/go-follow/authorization.service/internal/apiserver"
	"github.com/go-follow/authorization.service/pkg/logger"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		logger.Fatal(err)
	}
	apiserver.Start(cfg)
}

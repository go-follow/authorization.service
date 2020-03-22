package main

import (
	"fmt"

	"github.com/go-follow/authorization.service/config"
	"github.com/go-follow/authorization.service/pkg/logger"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(cfg)
}
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/klausborkowski/wordofwisdom/config"
	"github.com/klausborkowski/wordofwisdom/internal/client"
)

func main() {
	fmt.Println("start client")

	configs, err := config.ParseConfig()
	if err != nil {
		log.Panicf("failed to load configuration: %s\n", err.Error())
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configs)

	address := fmt.Sprintf("%s:%d", configs.ServerConfig.Host, configs.ServerConfig.Port)

	err = client.Run(ctx, address)
	if err != nil {
		fmt.Println("client error:", err)
	}
}

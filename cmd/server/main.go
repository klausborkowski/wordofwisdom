package main

import (
	"context"
	"fmt"
	"log"

	"github.com/klausborkowski/wordofwisdom/config"
	"github.com/klausborkowski/wordofwisdom/internal/cache"
	"github.com/klausborkowski/wordofwisdom/internal/clock"
	"github.com/klausborkowski/wordofwisdom/internal/quotes"
)

func main() {
	fmt.Println("start server")

	configs, err := config.ParseConfig()
	if err != nil {
		log.Panicf("failed to load configuration: %s\n", err.Error())
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, config.ConfigCtxKey, configs)
	ctx = context.WithValue(ctx, clock.ClockCtxKey, clock.SystemClock{})

	cache, err := cache.InitRedisCache(ctx, configs.CacheConfig.Host, configs.CacheConfig.Port)
	if err != nil {
		fmt.Println("error init cache:", err)
		return
	}
	ctx = context.WithValue(ctx, "cache", cache)

	serverAddress := fmt.Sprintf("%s:%d", configs.ServerConfig.Host, configs.ServerConfig.Port)
	err = quotes.RunServer(ctx, serverAddress)
	if err != nil {
		fmt.Println("server error:", err)
	}
}

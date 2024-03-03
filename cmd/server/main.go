package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ebdonato/go-rate-limiter/configs"
	db "github.com/ebdonato/go-rate-limiter/pkg/db/cache"
	server "github.com/ebdonato/go-rate-limiter/pkg/http/server"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	log.Println("Creating cache...")
	redis, err := db.NewRedisCache(fmt.Sprintf("%s:%s", configs.RedisHost, configs.RedisPort))
	if err != nil {
		panic(err)
	}

	r := server.NewWebServer(configs.RateLimiterMaxIPRequests, configs.RateLimiterMaxTokenRequests, configs.RateLimiterBlockTimeSeconds, redis)

	log.Println("Starting web server on port", configs.WebServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", configs.WebServerPort), r)
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ebdonato/go-rate-limiter/configs"
	server "github.com/ebdonato/go-rate-limiter/pkg/http/server"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	r := server.NewWebServer(configs.RateLimiterMaxIPRequests, configs.RateLimiterMaxTokenRequests, configs.RateLimiterBlockTimeSeconds, configs.RedisHost, configs.RedisPort)

	log.Println("Starting web server on port", configs.WebServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", configs.WebServerPort), r.Router)
}

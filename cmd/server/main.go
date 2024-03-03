package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ebdonato/go-rate-limiter/configs"
	db "github.com/ebdonato/go-rate-limiter/pkg/db/cache"
	middleware "github.com/ebdonato/go-rate-limiter/pkg/http/middleware"

	"github.com/go-chi/chi/v5"
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

	limiter := middleware.NewRateLimiter(configs.RateLimiterMaxIPRequests, configs.RateLimiterMaxTokenRequests, configs.RateLimiterBlockTimeSeconds, redis)

	r := chi.NewRouter()
	r.Use(limiter.RateLimit)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	})
	log.Println("Starting web server on port", configs.WebServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", configs.WebServerPort), r)
}

package http

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/ebdonato/go-rate-limiter/pkg/db/cache"
	middleware "github.com/ebdonato/go-rate-limiter/pkg/http/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router chi.Router
}

func NewWebServer(
	maxIpRequests, maxTokenRequests, blockTime int, redisHost, redisPort string) *Server {

	log.Println("Creating cache...")
	redis, err := db.NewRedisCache(fmt.Sprintf("%s:%s", redisHost, redisPort))
	if err != nil {
		panic(err)
	}

	log.Println("Limiter Configs: ", "maxIpRequests: ", maxIpRequests, " | maxTokenRequests: ", maxTokenRequests, " | blockTime: ", blockTime)
	limiter := middleware.NewRateLimiter(maxIpRequests, maxTokenRequests, blockTime, redis)

	r := chi.NewRouter()
	r.Use(limiter.RateLimit)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	})

	return &Server{
		Router: r,
	}
}

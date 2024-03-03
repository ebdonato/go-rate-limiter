package http

import (
	"log"
	"net/http"

	db "github.com/ebdonato/go-rate-limiter/pkg/db/cache"
	middleware "github.com/ebdonato/go-rate-limiter/pkg/http/middleware"
	"github.com/go-chi/chi/v5"
)

func NewWebServer(
	maxIpRequests, maxTokenRequests, blockTime int, cache db.Cache) *chi.Mux {

	log.Println("Limiter Configs: ", "maxIpRequests: ", maxIpRequests, " | maxTokenRequests: ", maxTokenRequests, " | blockTime: ", blockTime)
	limiter := middleware.NewRateLimiter(maxIpRequests, maxTokenRequests, blockTime, cache)

	r := chi.NewRouter()
	r.Use(limiter.RateLimit)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	})

	return r
}

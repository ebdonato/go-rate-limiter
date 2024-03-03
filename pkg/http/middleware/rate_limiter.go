package http

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "github.com/ebdonato/go-rate-limiter/pkg/db/cache"
)

type RateLimiter struct {
	maxIpRequests, maxTokenRequests, blockTime int
	cache                                      db.Cache
}

// from https://gist.github.com/miguelmota/7b765edff00dc676215d6174f3f30216
func getIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("REQUEST IP NOT FOUND")
}

func NewRateLimiter(maxIpRequests, maxTokenRequests, blockTime int, cache db.Cache) *RateLimiter {
	return &RateLimiter{
		maxIpRequests,
		maxTokenRequests,
		blockTime,
		cache,
	}
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getIP(r)
		if err != nil {
			log.Println("Failed to get remote address: ", err)
			http.Error(w, "Internal Server Error", http.StatusBadRequest)
			return
		}
		kind := "addr"
		if r.Header.Get("API_KEY") != "" {
			id = r.Header.Get("API_KEY")
			kind = "token"
		}

		key := fmt.Sprintf("%s:%s", kind, id)
		log.Println("Key: ", key)

		val, err := rl.cache.Get(key)
		if err != nil {
			log.Println("Failed to get value from cache: ", err)
			http.Error(w, "Internal Server Error", http.StatusBadGateway)
			return
		}

		if val == "" {
			val = "1"
			err := rl.cache.Set(key, val, time.Duration(rl.blockTime)*time.Second)
			if err != nil {
				log.Println("Failed to set value in cache: ", err)
				http.Error(w, "Internal Server Error", http.StatusBadGateway)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		count, err := strconv.Atoi(val)
		if err != nil {
			log.Println("Failed to convert value to int: ", err)
			http.Error(w, "Internal Server Error", http.StatusBadGateway)
			return
		}

		maxRequest := rl.maxIpRequests
		if kind == "token" {
			maxRequest = rl.maxTokenRequests
		}

		if count+1 > maxRequest {
			log.Println("Too many requests")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("You have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}

		err = rl.cache.Set(key, strconv.Itoa(count+1), time.Duration(rl.blockTime)*time.Second)
		if err != nil {
			log.Println("Filed to set value in cache: ", err)
			http.Error(w, "Internal Server Error", http.StatusBadGateway)
			return
		}

		next.ServeHTTP(w, r)
	})
}

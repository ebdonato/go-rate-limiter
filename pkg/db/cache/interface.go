package db

import "time"

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string, expiration time.Duration) error
}

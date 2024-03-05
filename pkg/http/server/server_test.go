package http

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/ebdonato/go-rate-limiter/pkg/db/cache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type webServerSuite struct {
	suite.Suite
	router           http.Handler
	maxIpRequests    int
	maxTokenRequests int
	blockTime        int
	redisHost        string
	redisPort        string
	redis            *db.RedisCache
}

func (suite *webServerSuite) SetupSuite() {
	suite.maxIpRequests = 10
	suite.maxTokenRequests = 100
	suite.blockTime = 5
	suite.redisHost = "localhost"
	suite.redisPort = "6379"

	log.Println("Creating cache...")
	redis, err := db.NewRedisCache(fmt.Sprintf("%s:%s", suite.redisHost, suite.redisPort))
	if err != nil {
		panic(err)
	}

	suite.redis = redis

	r := NewWebServer(suite.maxIpRequests, suite.maxTokenRequests, suite.blockTime, redis)

	suite.router = r
}

func (suite *webServerSuite) SetupTest() {
	suite.redis.Clear()
}

func (suite *webServerSuite) TestWebServerRunning() {
	ts := httptest.NewServer(suite.router)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *webServerSuite) TestIpRateLimits() {
	ts := httptest.NewServer(suite.router)
	defer ts.Close()

	okResponses := []int{}
	blockedResponses := []int{}

	const numberRequests = 100

	for i := 0; i < numberRequests; i++ {
		resp, err := http.Get(ts.URL)
		assert.NoError(suite.T(), err)

		switch resp.StatusCode {
		case http.StatusOK:
			okResponses = append(okResponses, resp.StatusCode)
		case http.StatusTooManyRequests:
			blockedResponses = append(blockedResponses, resp.StatusCode)
		}
	}
	assert.Equal(suite.T(), suite.maxIpRequests, len(okResponses))
	assert.Equal(suite.T(), numberRequests-suite.maxIpRequests, len(blockedResponses))
}

func (suite *webServerSuite) TestTokenRateLimits() {
	ts := httptest.NewServer(suite.router)
	defer ts.Close()

	okResponses := []int{}
	blockedResponses := []int{}

	const numberRequests = 1000

	for i := 0; i < numberRequests; i++ {
		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		assert.NoError(suite.T(), err)
		req.Header.Set("API_KEY", "some_token")
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(suite.T(), err)

		switch resp.StatusCode {
		case http.StatusOK:
			okResponses = append(okResponses, resp.StatusCode)
		case http.StatusTooManyRequests:
			blockedResponses = append(blockedResponses, resp.StatusCode)
		}
	}
	assert.Equal(suite.T(), suite.maxTokenRequests, len(okResponses))
	assert.Equal(suite.T(), numberRequests-suite.maxTokenRequests, len(blockedResponses))
}

func TestWebServerSuite(t *testing.T) {
	suite.Run(t, new(webServerSuite))
}

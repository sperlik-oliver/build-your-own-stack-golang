package network

import (
	"build-your-own/load-balancer/config"
	"build-your-own/load-balancer/pool"
	"log"
	"net/http"
	"time"
)

// Runs an infinite loop through the upstream pool. Checks health of upstream and sets the alive property accordingly.
//
// Should be ran as a go-routine on startup.
func Healthcheck() {
	i := 0
	upstreams := pool.GetUpstreams()
	for {
		time.Sleep(time.Duration(config.GetConfig().HealthcheckIntervalSeconds) * time.Second)
		upstream := pool.GetUpstream(i)
		if upstream == nil {
			continue
		}
		url := upstream.URL.String() + "/health"
		_, err := http.Get(url)
		if err == nil {
			log.Printf("Health check succeeded: " + url)
			upstream.SetAlive(true)
		} else {
			log.Printf("Health check failed: " + url)
			upstream.SetAlive(false)
		}
		upstream.SetRetries(0)
		i = (i + 1) % len(upstreams)
	}
}

// Fetches the next available upstream server and reverse proxies the HTTP request to it.
// Retries a request up to 3 times.
//
// After 3 retries upstream is marked as unhealthy. Periodic health checker will revive it if it responds successfully.
func Lb(w http.ResponseWriter, r *http.Request) {
	upstream := pool.GetNextUpstream(config.GetConfig().Mode)

	if upstream == nil {
		log.Fatal("Service unavailable. No healthy upstream servers.")
	}

	defer func() {
		upstream.SetConnections(upstream.GetConnetions() - 1)
	}()

	upstream.SetConnections(upstream.GetConnetions() + 1)

	upstream.ReverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		log.Printf("Request to upstream failed: [%+v] [%+v]\n", upstream, e.Error())

		upstream.SetRetries(upstream.GetRetries() + 1)
		if upstream.GetRetries() >= 3 {
			upstream.SetAlive(false)
			log.Printf("Request to upstream failed 3 times, upstream marked as unhealthy: [%+v]", upstream)
		}
		Lb(w, r)
	}

	log.Printf("Serving HTTP request to upstream: [%+v]", upstream)

	if upstream.GetRetries() > 0 {
		time.Sleep(time.Duration(config.GetConfig().RetryIntervalSeconds) * time.Second)
	}

	upstream.ReverseProxy.ServeHTTP(w, r)
	log.Printf("Request to upstream successful: [%+v]", upstream)
}

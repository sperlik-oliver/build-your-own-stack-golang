package pool

import (
	"log"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

var upstreamPool UpstreamPoolStruct

// Holds information about a single upstream
type Upstream struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
	mux          sync.RWMutex
	alive        bool
	retries      int
	connections  int
}

// Holds all upstreams
type UpstreamPoolStruct struct {
	upstreams []*Upstream
	current   uint64
	mux       sync.RWMutex
}

// Adds a new upstream to the pool in a thread safe manner
func AddUpstream(url *url.URL) {
	reverseProxy := httputil.NewSingleHostReverseProxy(url)
	upstreamPool.mux.RLock()
	upstreamPool.upstreams = append(upstreamPool.upstreams, &Upstream{
		URL:          url,
		alive:        true,
		ReverseProxy: reverseProxy,
	})
	upstreamPool.mux.RUnlock()
}

// Returns the upstream pool in a thread safe manner
func GetUpstreams() []*Upstream {
	upstreamPool.mux.RLock()
	upstreams := upstreamPool.upstreams
	upstreamPool.mux.RUnlock()
	return upstreams
}

// Returns a single upstream in a thread safe manner
func GetUpstream(index int) *Upstream {
	upstreamPool.mux.RLock()
	upstream := upstreamPool.upstreams[index]
	upstreamPool.mux.RUnlock()
	return upstream
}

// Sets the Alive property on an upstream server in a thread safe manner
func (upstream *Upstream) SetAlive(alive bool) {
	upstream.mux.RLock()
	upstream.alive = alive
	upstream.mux.RUnlock()
}

// Retrieves the Alive property on an upstream server in a thread safe manner
func (upstream *Upstream) IsAlive() bool {
	upstream.mux.RLock()
	alive := upstream.alive
	upstream.mux.RUnlock()
	return alive
}

// Sets the retries property om an upstream server in a thread safe manner
func (upstream *Upstream) SetRetries(retries int) {
	upstream.mux.RLock()
	upstream.retries = retries
	upstream.mux.RUnlock()
}

// Retrieves the retries property on an upstream server in a thread safe manner
func (upstream *Upstream) GetRetries() int {
	upstream.mux.RLock()
	retries := upstream.retries
	upstream.mux.RUnlock()
	return retries
}

// Sets the connections property om an upstream server in a thread safe manner
func (upstream *Upstream) SetConnections(connections int) {
	upstream.mux.RLock()
	upstream.connections = connections
	upstream.mux.RUnlock()
}

// Retrieves the connections property on an upstream server in a thread safe manner
func (upstream *Upstream) GetConnetions() int {
	upstream.mux.RLock()
	connections := upstream.connections
	upstream.mux.RUnlock()
	return connections
}

// GetNextUpstreamRoundRobin returns the next available upstream server from the pool and increases the current upstream pool index by 1 atomically
//
// Atomically means that there is no state of "currently being edited", only "not increased" and "increased".
//
// returns:
//
//	*Upstream: A pointer to the next available upstream server, or nil if none are available
func getNextUpstreamRoundRobin() *Upstream {
	nextIndex := int(atomic.AddUint64(&upstreamPool.current, 1) % uint64(len(GetUpstreams())))
	nextUpstream := GetUpstream(nextIndex)
	if nextUpstream.IsAlive() {
		return nextUpstream
	} else {
		return getNextUpstreamRoundRobin()
	}
}

// GetNextUpstreamByConnections returns the upstream server from the pool that has the lowest connections
//
// returns:
//
//	*Upstream: A pointer to the next upstream server with least connections, or nil if none are available
func getNextUpstreamByConnections() *Upstream {
	var minUpstream *Upstream = nil
	for _, upstream := range GetUpstreams() {
		if upstream.IsAlive() && (minUpstream == nil || upstream.connections < minUpstream.connections) {
			minUpstream = upstream
		}
	}
	return minUpstream
}

func GetNextUpstream(mode string) *Upstream {
	if mode == "roundrobin" {
		return getNextUpstreamRoundRobin()
	}
	if mode == "leastconnections" {
		return getNextUpstreamByConnections()
	}
	log.Panicf("Unsupported load balancing mode: %s", mode)
	panic(1)
}

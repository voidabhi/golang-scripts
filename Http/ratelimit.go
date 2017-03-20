package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/context"
)

var (
	ErrRateLimit = "Rate Limit Exceeded"	
)

const (
	bucketSize = 100
	reqPerSec  = 50
)

type limCounter struct {
	m sync.Mutex
	n int64
	t int64
}

func (lc *limCounter) try(t int64) (ok bool) {
	lc.m.Lock()
	defer lc.m.Unlock()

	d := t - lc.t
	if d < 0 {
		// Taking the abs value of d
		// is motivated by paranoia that
		// t might not always be monotonically increasing.
		d = -d
	}
	lc.n += reqPerSec * d
	if lc.n > bucketSize {
		lc.n = bucketSize
	}
	lc.t = t

	if lc.n < 1 {
		return false
	}
	lc.n--
	return true
}

func limitKeyIP(r *http.Request) string {
	if ip := r.Header.Get("Chain-Forwarded-For"); ip != "" {
		return ip
	}
	return r.Header.Get("X-Forwarded-For")
}

func limitKeyAuth(r *http.Request) string {
	id, secret := getAuthToken(r)
	return id + ":" + secret
}

var (
	limCtrsMu sync.Mutex //protects the following:
	limCtrs   = map[string]*limCounter{}
)

func init() {
	go func() {
		for range time.Tick(time.Hour) {
			limCtrsMu.Lock()
			for k, ctr := range limCtrs {
				ctr.m.Lock()
				n := ctr.n
				ctr.m.Unlock()
				if n > bucketSize {
					delete(limCtrs, k)
				}
			}
			limCtrsMu.Unlock()
		}
	}()
}

type limitHandler struct {
	h handler
	f func(r *http.Request) string
}

func (l limitHandler) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	k := l.f(r)
	limCtrsMu.Lock()
	ctr, ok := limCtrs[k]
	if !ok {
		ctr = &limCounter{n: bucketSize, t: time.Now().Unix()}
		limCtrs[k] = ctr
	}
	limCtrsMu.Unlock()

	ok = ctr.try(time.Now().Unix())
	if !ok {
		nRateLimit.Add()
		ip := r.Header.Get("X-Forwarded-For")
		user, _ := getAuthToken(r)
		log.Println("rate-limit:", ip, user)
		httpError(ctx, w, ErrRateLimit)
		return
	}
	l.h.ServeHTTPContext(ctx, w, r)
}

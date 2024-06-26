package ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/ranna-go/ranna/pkg/models"
	"github.com/zekroTJA/ratelimit"
	"github.com/zekroTJA/timedmap/v2"
)

const (
	cleanupInterval = 15 * time.Minute
	entryLifetime   = 1 * time.Hour
)

type Limiter interface {
	Allow() bool
}

type limit struct {
	Burst int
	Limit time.Duration
}

type dummyLimiter struct{}

func (dummyLimiter) Allow() bool {
	return true
}

type RateLimitManager struct {
	limits   map[models.OpCode]limit
	pool     *sync.Pool
	limiters *timedmap.TimedMap[string, Limiter]
}

func NewRateLimitManager(cfg ConfigProvider) *RateLimitManager {
	rlCfg := cfg.Config().API.WS.RateLimit

	limits := map[models.OpCode]limit{
		models.OpExec: {
			Burst: rlCfg.Burst,
			Limit: time.Duration(rlCfg.LimitSeconds) * time.Second,
		},
	}

	return &RateLimitManager{
		limits: limits,
		pool: &sync.Pool{
			New: func() interface{} {
				return ratelimit.NewLimiter(0, 0)
			},
		},
		limiters: timedmap.New[string, Limiter](cleanupInterval),
	}
}

func (rlm *RateLimitManager) GetLimiter(c *websocket.Conn, op models.OpCode) Limiter {
	limits, ok := rlm.limits[op]
	if !ok || limits.Burst == 0 && limits.Limit == 0 {
		return dummyLimiter{}
	}
	key := fmt.Sprintf("%d::%s", op, getAddr(c))
	limiter, ok := rlm.limiters.GetValue(key)
	if ok {
		return limiter
	}
	return rlm.createLimiter(key, limits.Limit, limits.Burst)
}

func (rlm *RateLimitManager) createLimiter(key string, limit time.Duration, burst int) Limiter {
	limiter := rlm.pool.Get().(*ratelimit.Limiter)
	limiter.SetLimit(limit)
	limiter.SetBurst(burst)
	limiter.Reset()
	rlm.limiters.Set(key, limiter, entryLifetime, func(v Limiter) {
		rlm.pool.Put(v)
	})
	return limiter
}

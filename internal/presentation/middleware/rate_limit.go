package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type visitor struct {
	count    int
	lastSeen time.Time
}

func RateLimitMiddleware(maxRequests int, window time.Duration) echo.MiddlewareFunc {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	go func() {
		for {
			time.Sleep(window)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > window {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()

			mu.Lock()
			v, exists := visitors[ip]
			if !exists {
				visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
				mu.Unlock()
				return next(c)
			}

			if time.Since(v.lastSeen) > window {
				v.count = 1
				v.lastSeen = time.Now()
				mu.Unlock()
				return next(c)
			}

			v.count++
			v.lastSeen = time.Now()

			if v.count > maxRequests {
				mu.Unlock()
				return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "Too many requests"})
			}

			mu.Unlock()
			return next(c)
		}
	}
}

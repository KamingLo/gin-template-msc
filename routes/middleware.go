package routes

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"template/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter       *rate.Limiter
	lastSeen      time.Time
	isLockedUntil time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func init() {
	go func() {
		for {
			time.Sleep(2 * time.Minute)
			mu.Lock()
			for ip, v := range clients {
				if time.Since(v.lastSeen) > 5*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

// AuthMiddleware ensures the request has a valid JWT.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.SendError(c, http.StatusUnauthorized, "Session expired, please log in again", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString, os.Getenv("JWT_SECRET"))
		if err != nil {
			utils.SendError(c, http.StatusUnauthorized, "Invalid or expired token", err)
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("user_email", claims["email"])
		c.Next()
	}
}

// RateLimitMiddleware prevents brute-force by tracking IP + Path.
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.Request.URL.Path
		identifier := ip + ":" + path
		now := time.Now()

		mu.Lock()
		v, found := clients[identifier]
		if !found {
			v = &client{
				limiter: rate.NewLimiter(rate.Every(30*time.Second/5), 5),
			}
			clients[identifier] = v
		}
		v.lastSeen = now

		if now.Before(v.isLockedUntil) || !v.limiter.Allow() {
			if now.After(v.isLockedUntil) {
				v.isLockedUntil = now.Add(30 * time.Second)
			}

			remaining := time.Until(v.isLockedUntil).Seconds()
			mu.Unlock()

			// Using utils.SendError for consistency
			msg := fmt.Sprintf("Rate limit exceeded (5 attempts per 30s). Please try again in %.0f seconds", remaining)
			utils.SendError(c, http.StatusTooManyRequests, msg, nil)
			c.Abort()
			return
		}

		mu.Unlock()
		c.Next()
	}
}

// CORSMiddleware manages cross-origin access based on environment settings.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		origins := strings.Split(allowedOrigins, ",")
		origin := c.GetHeader("Origin")

		isAllowed := false
		if origin == "" {
			isAllowed = true
		} else {
			for _, o := range origins {
				if origin == strings.TrimSpace(o) {
					isAllowed = true
					break
				}
			}
		}

		if !isAllowed {
			utils.SendError(c, http.StatusForbidden, "Access denied: Origin not allowed by CORS policy", nil)
			c.Abort()
			return
		}

		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

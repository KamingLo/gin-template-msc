package routes

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

var (
	mu      sync.Mutex
	clients = make(map[string]*rate.Limiter)
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 7 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token dibutuhkan"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kedaluwarsa"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.ClientIP() mendukung X-Forwarded-For (Next.js) & Direct IP (Flutter)
		ip := c.ClientIP()

		mu.Lock()
		if _, found := clients[ip]; !found {
			// Limit: 5 request per detik, Burst: 10
			clients[ip] = rate.NewLimiter(5, 10)
		}
		limiter := clients[ip]
		mu.Unlock()

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Terlalu banyak permintaan, akses dibatasi",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		origins := strings.Split(allowedOrigins, ",")
		origin := c.GetHeader("Origin")

		isAllowed := false
		// Jika origin kosong (misal request dari Postman atau Mobile App non-browser)
		// Biasanya kita izinkan, atau Anda bisa perketat jika mau.
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

		// 1. Jika TIDAK diizinkan, langsung stop dan beri pesan error
		if !isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Akses ditolak, Origin tidak diizinkan oleh kebijakan CORS",
			})
			c.Abort()
			return
		}

		// 2. Jika diizinkan, set header Allow-Origin
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 3. Tangani Pre-flight OPTIONS
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

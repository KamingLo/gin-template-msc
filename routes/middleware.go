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
				// Hapus IP yang sudah tidak aktif lebih dari 5 menit
				if time.Since(v.lastSeen) > 5*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil Header Authorization
		authHeader := c.GetHeader("Authorization")

		// 2. Validasi format: Harus dimulai dengan "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// Gunakan utils.SendError agar format JSON seragam
			utils.SendError(c, http.StatusUnauthorized, "Sesi berakhir, silakan login kembali", nil)
			c.Abort()
			return
		}

		// 3. Potong string untuk mendapatkan token murni
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 4. Validasi JWT
		claims, err := utils.ValidateToken(tokenString, os.Getenv("JWT_SECRET"))
		if err != nil {
			utils.SendError(c, http.StatusUnauthorized, "Token tidak valid atau kadaluarsa", err)
			c.Abort()
			return
		}

		// 5. Simpan data user ke context Gin
		// Pastikan claims["user_id"] dan claims["email"] sesuai dengan struct JWT kamu
		c.Set("user_id", claims["user_id"])
		c.Set("user_email", claims["email"])

		c.Next()
	}
}

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
			fmt.Printf("[RateLimit] New Registration: %s\n", identifier)
		}
		v.lastSeen = now

		if now.Before(v.isLockedUntil) || !v.limiter.Allow() {
			if now.After(v.isLockedUntil) {
				v.isLockedUntil = now.Add(30 * time.Second)
				fmt.Printf("[RateLimit] LIMIT TRIGGERED: %s | Locked for 30s\n", identifier)
			}

			remaining := time.Until(v.isLockedUntil).Seconds()
			mu.Unlock()

			fmt.Printf("[RateLimit] REJECTED: %s | Retry in %.0fs\n", identifier, remaining)

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": fmt.Sprintf("Batas 5 kali percobaan per 30 detik tercapai. Coba lagi dalam %.0f detik", remaining),
			})
			return
		}

		mu.Unlock()

		c.Next()
	}
}

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
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Akses ditolak, Origin tidak diizinkan oleh kebijakan CORS",
			})
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
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

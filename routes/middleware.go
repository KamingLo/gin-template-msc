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
		var tokenString string

		// 1. Cek Cookie (Untuk Web)
		if cookie, err := c.Cookie("auth_token"); err == nil {
			tokenString = cookie
		} else {
			// 2. Cek Header (Untuk Mobile/Flutter)
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = authHeader[7:]
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi berakhir, silakan login kembali"})
			c.Abort()
			return
		}

		// Validasi JWT (Gunakan utils.ValidateToken buatanmu)
		claims, err := utils.ValidateToken(tokenString, os.Getenv("JWT_SECRET"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		// Simpan data user ke context agar bisa dipakai di controller
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

		// GABUNGKAN: Cek apakah masih terkunci ATAU apakah jatah token habis
		if now.Before(v.isLockedUntil) || !v.limiter.Allow() {
			// Jika baru saja kena limit (belum ada isLockedUntil), pasang lock-nya
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

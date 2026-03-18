package config

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func InitOAuth() {
	// Key untuk session cookie (ganti dengan string random di .env)
	key := os.Getenv("SESSION_SECRET")
	maxAge := 86400 * 30

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false // Set true jika sudah HTTPS

	gothic.Store = store

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_CALLBACK_URL"), // Sesuaikan domain
			"email", "profile",
		),
	)
}

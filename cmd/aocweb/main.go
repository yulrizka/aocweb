package main

import (
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/yulrizka/aocweb"
)

func env(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	} else {
		return def
	}
}

func main() {
	addr := flag.String("http-addr", env("AOC_HTTP_ADDR", "localhost:3001"), "")
	templates := flag.String("templates", env("AOC_TEMPLATES", "templates"), "templates folder")

	r := gin.Default()
	r.LoadHTMLGlob(*templates + "/*")

	cfg := aocweb.MustConfig()

	var sessionStore = sessions.NewCookieStore([]byte(cfg.SessionKey))
	auth := aocweb.NewAuth(sessionStore)

	handler := aocweb.Handler{
		Config:       cfg,
		SessionStore: sessionStore,
		Auth:         auth,
	}

	// middle ware
	authorized := r.Group("/", auth.Required())
	authorized.GET("/submit", handler.Submit)

	r.GET("/", handler.Index)
	r.GET("/auth/github", handler.GithubAuth)
	r.GET("/auth/github/callback", handler.GithubAuthCallback)

	if err := r.Run(*addr); err != nil {
		log.Fatal("run ", err)
	}
}

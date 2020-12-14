package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"

	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/yulrizka/aocweb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	r.HTMLRender = loadTemplates(*templates)

	cfg := aocweb.MustConfig()

	var sessionStore = cookie.NewStore([]byte(cfg.SessionKey))
	auth := aocweb.NewAuth(sessionStore)
	r.Use(sessions.Sessions("aocweb", sessionStore))

	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	aocweb.AutoMigrate(db)

	var handler = aocweb.Handler{
		Config:       cfg,
		SessionStore: sessionStore,
		Auth:         auth,
		DB:           db,
	} // middle ware
	authorized := r.Group("/", auth.Required())
	authorized.GET("/logout", auth.Logout)
	authorized.GET("/submit", handler.Submit)
	authorized.POST("/submit", handler.DoSubmit)

	r.GET("/", handler.Index)
	r.GET("/login", handler.GithubAuth)
	r.GET("/auth/github", handler.GithubAuth)
	r.GET("/auth/github/callback", handler.GithubAuthCallback)

	if err := r.Run(*addr); err != nil {
		log.Fatal("run ", err)
	}
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.gohtml")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.gohtml")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}

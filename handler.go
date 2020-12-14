package aocweb

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

type Handler struct {
	Config       *Config
	SessionStore cookie.Store
	Auth         *Auth
	DB           *gorm.DB
}

func (h *Handler) Index(c *gin.Context) {
	user := h.Auth.User(c)
	c.HTML(http.StatusOK, "index.gohtml", gin.H{
		"title": "(unofficial) Advent of Code solutions",
		"user":  user,
	})
}

func (h *Handler) GithubAuth(c *gin.Context) {
	conf := &oauth2.Config{
		ClientID:     h.Config.GitHub.ClientID,
		ClientSecret: h.Config.GitHub.ClientSecret,
		Endpoint:     github.Endpoint,
	}

	u := conf.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func (h *Handler) GithubAuthCallback(c *gin.Context) {
	client := &http.Client{Timeout: 2 * time.Second}
	ctx := context.WithValue(c.Request.Context(), oauth2.HTTPClient, client)

	cfg := h.Config.GitHub.OAuth2
	tok, err := cfg.Exchange(ctx, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	client = cfg.Client(ctx, tok)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.Error(err)
		return
	}
	defer resp.Body.Close()

	//just use what is necessary
	info := struct {
		Name  string `json:"name"`
		Login string `json:"login"`
	}{}
	//info := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		c.Error(err)
		return
	}

	var u User
	err = h.DB.Where("github = ?", info.Login).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		u.Name = info.Name
		u.Github = info.Login
		if err := h.DB.Create(&u).Error; err != nil {
			c.Error(err)
			return
		}
	} else if err != nil {
		c.Error(err)
		return
	}
	h.Auth.WriteUser(c, u)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *Handler) Submit(c *gin.Context) {
	user := h.Auth.User(c)
	c.HTML(http.StatusOK, "submit.gohtml", gin.H{
		"title": "Submit solution",
		"user":  user,
	})
}

func (h *Handler) DoSubmit(c *gin.Context) {
	user := h.Auth.User(c)
	_ = user
}

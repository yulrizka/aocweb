package aocweb

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Handler struct {
	Config       *Config
	SessionStore *sessions.CookieStore
	Auth         *Auth
}

func (h *Handler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
		"user":  h.Auth.User(c),
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

	// just use what is necessary
	info := struct {
		Login string `json:"login"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		c.Error(err)
		return
	}
	h.Auth.WriteUser(c, info.Login)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *Handler) Submit(c *gin.Context) {
	user := h.Auth.User(c)
	_ = user
}

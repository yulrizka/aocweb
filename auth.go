package aocweb

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

const userKey = "auth"
const userSession = "user-session"

type Auth struct {
	store *sessions.CookieStore
}

func NewAuth(store *sessions.CookieStore) *Auth {
	return &Auth{
		store: store,
	}
}

func (a *Auth) User(c *gin.Context) string {
	s, err := a.store.Get(c.Request, userSession)
	if err != nil {
		return ""
	}
	v, ok := s.Values[userKey].(string)
	if !ok {
		return ""
	}

	return v
}

func (a *Auth) WriteUser(c *gin.Context, user string) {
	s, err := a.store.Get(c.Request, userSession)
	if err != nil {
		c.Error(err)
	}
	s.Values[userKey] = user
	if err := s.Save(c.Request, c.Writer); err != nil {
		c.Error(err)
	}
}

func (a *Auth) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		if a.User(c) == "" {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
	}
}

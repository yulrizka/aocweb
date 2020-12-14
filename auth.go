package aocweb

import (
	"net/http"

	"github.com/gin-contrib/sessions"

	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

const userKey = "auth"

type Auth struct {
	store cookie.Store
}

func NewAuth(store cookie.Store) *Auth {
	return &Auth{
		store: store,
	}
}

func (a *Auth) User(c *gin.Context) *User {
	s := sessions.Default(c)
	user, ok := s.Get(userKey).(User)
	if !ok {
		return nil
	}
	return &user
}

func (a *Auth) WriteUser(c *gin.Context, user User) {
	s := sessions.Default(c)

	s.Set(userKey, user)
	if err := s.Save(); err != nil {
		c.Error(err)
		return
	}
}

func (a *Auth) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		if a.User(c) == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
	}
}

func (a *Auth) Logout(c *gin.Context) {
	s := sessions.Default(c)
	s.Delete(userKey)
	if err := s.Save(); err != nil {
		c.Error(err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

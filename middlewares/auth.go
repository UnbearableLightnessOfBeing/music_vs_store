package middlewares

import (
	"log"
	db "music_vs_store/db/sqlc"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
  queries *db.Queries
}

func NewAuthMiddleware(queries *db.Queries) AuthMiddleware {
  return AuthMiddleware{
    queries,
  }
} 

func (m AuthMiddleware) RequireAuth() gin.HandlerFunc {
  return func(c *gin.Context) {
    session := sessions.Default(c)
    sessionID := session.Get("id")

    isUserPresent := true

    var user db.User
    var err error

    if sessionID == nil {
      isUserPresent = false
    } else {
      user, err = m.queries.GetUser(c, sessionID.(int32))
      if err != nil {
        panic(err)
      }
      isUserPresent = user.ID > 0
    }

    if isUserPresent {
      c.Set("user_id", uint64(user.ID))
      c.Set("user_username", user.Username)
    }

    c.Next()
  }
}

func (m AuthMiddleware) RequireAdmin() gin.HandlerFunc {
  return func(c *gin.Context) {
    session := sessions.Default(c)
    sessionID := session.Get("id")

    if sessionID != nil {
      user, err := m.queries.GetUser(c, sessionID.(int32))
      if err != nil {
        panic(err)
      }

      log.Println("valid:", user.IsAdmin.Valid)
      log.Println("bool:", user.IsAdmin.Bool)

      if user.IsAdmin.Valid && user.IsAdmin.Bool {
        c.Next()
        return
      }
    }
    c.Redirect(http.StatusMovedPermanently, "/")
  }
}

func (m *AuthMiddleware) RequireAdminInPanel() gin.HandlerFunc {
  return func(c *gin.Context) {
    session := sessions.Default(c)
    sessionID := session.Get("id")

    if sessionID != nil {
      user, err := m.queries.GetUser(c, sessionID.(int32))
      if err != nil {
        panic(err)
      }

      if user.IsAdmin.Valid && user.IsAdmin.Bool {
        c.Next()
        return
      }
    }
    c.JSON(http.StatusForbidden, gin.H{
      "message": "You are not allowed to see this page",
    })
  }
}

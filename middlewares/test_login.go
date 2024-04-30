package middlewares

import (
	"music_vs_store/helpers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginForTesting() gin.HandlerFunc {
  return func(c *gin.Context) {
    session := sessions.Default(c)
    sessionID := session.Get("id")

    if sessionID == nil || sessionID == 0 {
      helpers.SetSession(c, 1)
    }

    c.Next()
  }
}

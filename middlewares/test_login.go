package middlewares

import (
	"music_vs_store/helpers"

	"github.com/gin-gonic/gin"
)

func LoginForTesting() gin.HandlerFunc {
  return func(c *gin.Context) {
    sessionID := helpers.GetSession(c)

    if sessionID == 0 {
      helpers.SetSession(c, 1)
    }

    c.Next()
  }
}

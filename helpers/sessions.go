package helpers

import "github.com/gin-contrib/sessions"
import "github.com/gin-gonic/gin"

func SetSession(c *gin.Context, userID int32) {
  session := sessions.Default(c)
  var idInterface interface{} = &userID
  session.Set("id", idInterface)
  session.Save()
}

func GetSession(c *gin.Context) int32 {
  session := sessions.Default(c)
  id, ok :=  session.Get("id").(int32)
  if !ok {
    return 0
  }
  return id
}

func ClearSession(c *gin.Context) {
  session := sessions.Default(c)
  session.Clear()
  session.Save()
}

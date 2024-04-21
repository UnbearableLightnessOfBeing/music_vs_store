package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
  Username string `form:"username" binding:"required"`
	Email    string `form:"email" binding:"required"`
}

func Logger() gin.HandlerFunc {
  return func(c *gin.Context) {
    t := time.Now()

    c.Next()

    latency := time.Since(t)
    log.Println("latency: ", latency)

    status := c.Writer.Status()
    log.Println("status: ", status)
  }
}

func RequireAuth() gin.HandlerFunc {
  return func(c *gin.Context) {
    u := User{}

    err := c.ShouldBind(&u)
    if err != nil {
      c.AbortWithError(http.StatusBadRequest, err)
      return
    }

    if u.Username != "admin" {
      c.String(http.StatusUnauthorized, "You are not authorized to see this page")
      return
    }

    c.Next()
  }
}

func mainPageHandler(c *gin.Context) {
  c.String(200, "main page content") 
}

func protectedPageHandler(c *gin.Context) {
  c.String(200, "this is a protected route") 
}

func main() {
	r := gin.Default()
  r.Use(Logger())
	r.GET("/", mainPageHandler)
	r.GET("/protected", RequireAuth(), protectedPageHandler)
	r.Run()
}

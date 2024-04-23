package main

import (
	"fmt"
	"music_vs_store/controllers"
	"music_vs_store/driver"
	"music_vs_store/middlewares"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"

	_ "github.com/lib/pq"
)

func init() {
  gotenv.Load()
}

func main() {
	r := gin.Default()
  r.Static("/styles", "./static/styles")
  r.LoadHTMLGlob("templates/**/*")

  queries := driver.GetQueries()

  // sessions
  store := memstore.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("users", store))

  // middlewares
  authMiddleware := middlewares.NewAuthMiddleware(queries)
  r.Use(authMiddleware.RequireAuth())

  // controllers
  usersController := controllers.NewUsersController(queries)
  sessionsController := controllers.NewSessionsController(queries)
  dashboardController := controllers.NewDashboardController(queries)

  r.GET("/", usersController.ListUsers)
  r.GET("/signup", usersController.CreateUserView)
  r.GET("/login", usersController.LoginView)
  r.GET("/admin", authMiddleware.RequireAdmin(), dashboardController.Index)

  r.POST("/signup", sessionsController.Signup)
  r.POST("/login", sessionsController.Login)
  r.POST("/logout", sessionsController.Logout)

  // admin api
  r.POST("/admin/categories", authMiddleware.RequireAdmin(), dashboardController.CreateCategory)

  // HTMX test
  r.GET("/admin/htmx", authMiddleware.RequireAdmin(), dashboardController.TestHtmx)

  var port = os.Getenv("SERVER_PORT")
  fmt.Println("starting server at: " + port)
  r.Run(port)
}

package main

import (
	"fmt"
	"music_vs_store/controllers"
	"music_vs_store/driver"
	"music_vs_store/middlewares"
	"music_vs_store/web"
	"os"
	"text/template"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"

	_ "github.com/lib/pq"
)

func init() {
  gotenv.Load()
}

func Mul(param1 int32, param2 int32) int32 {
  return param1 * param2
}

func main() {
	r := gin.Default()
  r.Static("/styles", "./static/styles")
  r.Static("/js", "./static/js")
  r.Static("/assets", "./static/assets")
  r.Static("/storage", "./storage/images/")
  r.SetFuncMap(template.FuncMap{
    "mul": Mul,
  })
  r.LoadHTMLGlob("templates/**/*")

  queries, db := driver.GetQueriesWithDb()

  // sessions
  store := memstore.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("users", store))

  // for TESTING
  // r.Use(middlewares.LoginForTesting())

  // middlewares
  authMiddleware := middlewares.NewAuthMiddleware(queries)
  r.Use(authMiddleware.RequireAuth())

  // controllers
  // usersController := controllers.NewUsersController(queries)
  sessionsController := controllers.NewSessionsController(queries)
  dashboardController := controllers.NewDashboardController(queries)

  webController := web.NewWebController(queries, db)

  // web
  r.GET("/", webController.RenderMainPage)
  r.GET("/catalogue", webController.RenderCataloguePage)
  r.GET("/catalogue/:slug", webController.RenderCategoryPage)
  r.GET("/catalogue/:slug/htmx", webController.RenderProducts)
  r.GET("/catalogue/:slug/:id", webController.RenderProductPage)
  r.GET("/about", webController.RenderAboutPage)
  r.GET("/comments", webController.RenderCommentsPage)
  r.GET("/delivery", webController.RenderDeliveryPage)
  r.GET("/contacts", webController.RenderContactsPage)
  r.GET("/cart", webController.RenderCartPage)
  r.GET("/checkout", webController.RenderCheckoutPage)
  r.GET("/orders", webController.RenderOrdersPage)
  r.GET("/search", webController.RenderSearchPage)

  // auth
  r.GET("/signup", webController.RenderSignupPage)
  r.GET("/login", webController.RenderLoginPage)
  r.GET("/admin", authMiddleware.RequireAdmin(), dashboardController.Index)

  r.POST("/signup", sessionsController.Signup)
  r.POST("/login", sessionsController.Login)
  r.POST("/logout", sessionsController.Logout)

  // htmx
  r.POST("/add-to-cart", webController.AddItemToCart)
  r.POST("/decrement-quantity", webController.DecrementQuantity)
  r.POST("/increment-quantity", webController.IncrementQuantity)
  r.DELETE("/delete-cart-item/:product_id", webController.DeleteCartItem)
  r.POST("/orders", webController.CreateOrder)
  r.POST("/search", webController.SearchItems)


  // admin api
  r.POST("/admin/categories", authMiddleware.RequireAdmin(), dashboardController.CreateCategory)

  // HTMX test
  r.GET("/admin/htmx", authMiddleware.RequireAdmin(), dashboardController.TestHtmx)

  var port = os.Getenv("SERVER_PORT")
  fmt.Println("starting server at: " + port)
  r.Run(port)
}

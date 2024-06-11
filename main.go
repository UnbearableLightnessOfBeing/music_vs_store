package main

import (
	"fmt"
	"music_vs_store/api"
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

func Arr(elements ...any) []any {
  return elements
}

func main() {
	r := gin.Default()
  r.Static("/styles", "./static/styles")
  r.Static("/js", "./static/js")
  r.Static("/assets", "./static/assets")
  r.Static("/storage", "./storage/images/")
  r.SetFuncMap(template.FuncMap{
    "mul": Mul,
    "arr": Arr,
  })
  r.LoadHTMLGlob("templates/**/*")

  queries, db := driver.GetQueriesWithDb()

  // sessions
  store := memstore.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("users", store))

  r.Use(func(c *gin.Context) {
    c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
    c.Next()
  })

  // for TESTING
  // r.Use(middlewares.LoginForTesting())

  // middlewares
  authMiddleware := middlewares.NewAuthMiddleware(queries)
  r.Use(authMiddleware.RequireAuth())

  // controllers
  // usersController := controllers.NewUsersController(queries)
  sessionsController := controllers.NewSessionsController(queries)
  // dashboardController := controllers.NewDashboardController(queries)

  webController := web.NewWebController(queries, db)
  apiController := api.NewApiController(queries, db)

  // web
  r.GET("/", webController.RenderMainPage)
  r.GET("/catalogue", webController.RenderCataloguePage)
  r.GET("/catalogue/:slug", webController.RenderCategoryPage)
  r.GET("/catalogue/:slug/htmx", webController.RenderProducts)
  r.GET("/products/:id", webController.RenderProductPage)
  r.GET("/about", webController.RenderAboutPage)
  r.GET("/comments", webController.RenderCommentsPage)
  r.GET("/delivery", webController.RenderDeliveryPage)
  r.GET("/contacts", webController.RenderContactsPage)
  r.GET("/cart", webController.RenderCartPage)
  r.GET("/checkout", webController.RenderCheckoutPage)
  r.GET("/orders", webController.RenderOrdersPage)
  r.GET("/orders/:id", webController.RenderOrderPage)
  r.GET("/search", webController.RenderSearchPage)

  // auth
  r.GET("/signup", webController.RenderSignupPage)
  r.GET("/login", webController.RenderLoginPage)
  // r.GET("/admin", authMiddleware.RequireAdmin(), dashboardController.Index)

  r.POST("/signup", sessionsController.Signup)
  r.POST("/login", sessionsController.Login)
  r.POST("/logout", sessionsController.Logout)

  // htmx
  r.POST("/add-to-cart", webController.AddItemToCart)
  r.POST("/buy-product", webController.BuyProduct)
  r.POST("/decrement-quantity", webController.DecrementQuantity)
  r.POST("/increment-quantity", webController.IncrementQuantity)
  r.DELETE("/delete-cart-item/:product_id", webController.DeleteCartItem)
  r.POST("/orders", webController.CreateOrder)
  r.POST("/search", webController.SearchItems)
  r.POST("/comments", webController.CreateComment)

  api := r.Group("/api/admin")
  api.Use(authMiddleware.RequireAdminInPanel())
  // admin api
  // ---AUTH---
  api.GET("/check-auth", authMiddleware.RequireAdminInPanel(), apiController.CheckAuth)
  // ---USERS---
  api.GET("/users", authMiddleware.RequireAdminInPanel(), apiController.Users)
  api.PUT("/users/:id", apiController.UserToggleIsAdmin)

  // ---PRODUCTS---
  api.GET("/products", apiController.Products)
  api.GET("/products/:id", apiController.Product)
  api.POST("/products", apiController.CreateProduct)
  api.PUT("/products/:id", apiController.UpdateProduct)
  api.DELETE("/products/:id", apiController.DeleteProduct)
  api.POST("/products/:id/images_add", apiController.AddImageToProdut)
  api.POST("/products/:id/images_remove", apiController.RemoveImageFromProduct)

  // ---CATEGORIES---
  api.GET("/categories", apiController.Categories)
  api.POST("/categories", apiController.CreateCategory)
  api.GET("/categories/:id", apiController.Category)
  api.PUT("/categories/:id", apiController.UpdateCategory)
  api.DELETE("/categories/:id", apiController.DeleteCategory)
  api.POST("/categories/:id/image", apiController.SetCategoryImage)

  // ---LABELS---
  api.GET("/labels", apiController.Labels)
  api.GET("/labels/:id", apiController.Label)
  api.POST("/labels", apiController.CreateLabel)
  api.PUT("/labels/:id", apiController.UpdateLabel)
  api.DELETE("/labels/:id", apiController.DeleteLabel)

  // ---ORDERS---
  api.GET("/orders", apiController.Orders)
  api.GET("/orders/:id", apiController.Order)
  // ---DASHBOARD---
  api.GET("/dashboard", apiController.Dashboard)

  // HTMX test
  // r.GET("/admin/htmx", authMiddleware.RequireAdmin(), dashboardController.TestHtmx)

  var port = os.Getenv("SERVER_PORT")
  fmt.Println("starting server at: " + port)
  r.Run(port)
}

package controllers

import (
	db "music_vs_store/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	queries *db.Queries
}

func NewDashboardController(queries *db.Queries) DashboardController {
    return DashboardController{
    queries,
  }
}

func (d DashboardController) Index(c *gin.Context) {

  categories, err := d.queries.ListCategories(c, db.ListCategoriesParams{
    Limit: 10,
    Offset: 0,
  })

  if err != nil {
    panic(err)
  }

	c.HTML(http.StatusOK, "dashboard/index.html", gin.H{
		"title": "Welcome to music_vs_store",
    "categories": categories,
	})
}

type CreateCategoryParams struct {
  Name string `form:"name" binding:"required"`
}

func (d DashboardController) CreateCategory(c *gin.Context) {
  var params CreateCategoryParams 
  c.ShouldBind(&params)

  _, err := d.queries.CreateCategory(c, params.Name)
  if err != nil {
    panic(err)
  }

  c.Redirect(http.StatusMovedPermanently, "/admin")

  // c.HTML(http.StatusOK, "dashboard/index.html", gin.H{})
} 

func (d DashboardController) TestHtmx(c *gin.Context) {
  c.HTML(http.StatusOK, "dashboard/htmx.html", gin.H{})
}

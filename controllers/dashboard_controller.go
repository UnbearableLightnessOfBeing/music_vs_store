package controllers

import (
	"database/sql"
	db "music_vs_store/db/sqlc"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

  created, err := d.queries.CreateCategory(c, params.Name)
  if err != nil {
    panic(err)
  }

  file, err := c.FormFile("image")
  if err != nil {
    panic(err)
  }
  ext := filepath.Ext(file.Filename)
  uuidStr := uuid.New().String()

  dst := "./storage/images/" + uuidStr + ext
  if err = c.SaveUploadedFile(file, dst); err != nil {
    panic(err)
  }

  imgUrl := "/storage/" + uuidStr + ext
  updateParams := db.UpdateCategoryImageUrlParams{
    ID: created.ID,
    ImgUrl: sql.NullString{
      String: imgUrl,
      Valid: true,
    },
  }
  _, err = d.queries.UpdateCategoryImageUrl(c, updateParams)
  if err != nil {
    panic(err)
  }

  c.Redirect(http.StatusMovedPermanently, "/admin")
} 

func (d DashboardController) TestHtmx(c *gin.Context) {
  c.HTML(http.StatusOK, "dashboard/htmx.html", gin.H{})
}

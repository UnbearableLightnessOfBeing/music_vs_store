package web

import (
	db "music_vs_store/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebController struct {
	queries *db.Queries
}

func NewWebController(queries *db.Queries) WebController {
    return WebController{
    queries,
  }
}


func (w WebController) RenderMainPage(c *gin.Context) {
  categories, err := w.queries.ListCategories(c, db.ListCategoriesParams{
    Limit: 10,
    Offset: 0,
  })
  if err != nil {
    panic(err)
  }

	c.HTML(http.StatusOK, "home/index.html", gin.H{
		"categories": categories,
    "isLoggedIn": c.GetUint64("user_id") > 0,
	})
}

func (w WebController) RenderCataloguePage(c *gin.Context) {
  categories, err := w.queries.ListCategories(c, db.ListCategoriesParams{
    Limit: 999,
    Offset: 0,
  })
  if err != nil {
    panic(err)
  }

	c.HTML(http.StatusOK, "home/catalogue.html", gin.H{
		"categories": categories,
    "isLoggedIn": c.GetUint64("user_id") > 0,
	})
}

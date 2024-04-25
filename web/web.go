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

type PagesInfo struct {
  // name - value map
	Pages       map[string]string
	CurrentPage string
}

var pages = map[string]string{
  "Главная": "home",
  "Каталог": "catalogue",
  "Категории": "categories",
  "О нас": "about",
  "Контакты": "contacts",
}

func (w WebController) RenderMainPage(c *gin.Context) {
	categories, err := w.queries.ListCategories(c, db.ListCategoriesParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "home/index.html", gin.H{
		"pages":      PagesInfo{
      Pages: pages,
      CurrentPage: "home",
    },
		"title":      "aboba",
		"categories": categories,
		"isLoggedIn": c.GetUint64("user_id") > 0,
	})
}

func (w WebController) RenderCataloguePage(c *gin.Context) {
	categories, err := w.queries.ListCategories(c, db.ListCategoriesParams{
		Limit:  999,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "home/catalogue.html", gin.H{
		"pages":      PagesInfo{
      Pages: pages,
      CurrentPage: "catalogue",
    },
		"categories":  categories,
		"isLoggedIn":  c.GetUint64("user_id") > 0,
	})
}

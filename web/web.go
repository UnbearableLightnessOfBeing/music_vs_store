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

type PageLink struct {
	Name  string
	Value string
}

type PagesInfo struct {
	// name - value map
	Pages       []PageLink
	CurrentPage string
}

var pages = []PageLink{
	{Name: "Главная", Value: ""},
	{Name: "Каталог", Value: "catalogue"},
	{Name: "О компании", Value: "about"},
	{Name: "Отзывы", Value: "comments"},
	{Name: "Доставка и оплата", Value: "delivery"},
	{Name: "Контакты", Value: "contacts"},
}

func (w WebController) RenderMainPage(c *gin.Context) {
	categories, err := w.queries.ListCategories(c, db.ListCategoriesParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/index.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "",
		},
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

	c.HTML(http.StatusOK, "web/catalogue.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "catalogue",
		},
		"categories": categories,
		"isLoggedIn": c.GetUint64("user_id") > 0,
	})
}

func (w WebController) RenderAboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "web/about.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "about",
		},
		"message": "About page",
	})
}

func (w WebController) RenderCommentsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "web/comments.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "comments",
		},
		"message": "Comments page",
	})
}

func (w WebController) RenderDeliveryPage(c *gin.Context) {
	c.HTML(http.StatusOK, "web/delivery.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "delivery",
		},
		"message": "Delivery page",
	})
}

func (w WebController) RenderContactsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "web/contacts.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "contacts",
		},
		"message": "Contacts page",
	})
}

type CategoryPage struct {
	Uri string `uri:"slug" binding:"required"`
}

type CategoryFilters struct {
	MinPrice int32 `form:"min_price"`
	MaxPrice int32 `form:"max_price"`
	LabelID  int32 `form:"label_id"`
}

func (w WebController) RenderCategoryPage(c *gin.Context) {
	var page CategoryPage
	if err := c.ShouldBindUri(&page); err != nil {
		panic(err)
	}

	category, err := w.queries.GetCategoryBySlug(c, page.Uri)
	if err != nil {
		panic(err)
	}

	var query CategoryFilters

	if err := c.ShouldBindQuery(&query); err != nil {
		panic(err)
	}

	products, err := w.queries.GetProductsByCategory(c, db.GetProductsByCategoryParams{
		CategoryID: category.ID,
		MinPrice:   query.MinPrice,
		MaxPrice:   query.MaxPrice,
		LabelID:    query.LabelID,
	})
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/category.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "categories",
		},
		"isLoggedIn":   c.GetUint64("user_id") > 0,
		"categoryName": category.Name,
		"products":     products,
	})
}

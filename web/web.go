package web

import (
	"database/sql"
	db "music_vs_store/db/sqlc"
	"music_vs_store/helpers"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type WebController struct {
	queries *db.Queries
  db *sql.DB
}

func NewWebController(queries *db.Queries, db *sql.DB) WebController {
	return WebController{
		queries,
    db,
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

func getCartProductsCount(c *gin.Context, q *db.Queries) (int32, error) {
  userID := helpers.GetSession(c)

  var cartProductsCount int64 = 0

  session, err := q.GetShoppingSessionByUserId(c, userID)
  if err != nil && err != sql.ErrNoRows {
    return 0, err
  }

  cartProductsCount, err = q.GetCartProductsCount(c, session.ID)
  if err != nil && err != sql.ErrNoRows {
    return 0, err
  }

  return int32(cartProductsCount), nil
}

func (w WebController) RenderMainPage(c *gin.Context) {
	categories, err := w.queries.ListCategories(c, db.ListCategoriesParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}

  cartProductsCount, err := getCartProductsCount(c, w.queries)
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
    "cartProductsCount": cartProductsCount,
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

  cartProductsCount, err := getCartProductsCount(c, w.queries)
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
    "cartProductsCount": cartProductsCount,
	})
}

func (w WebController) RenderAboutPage(c *gin.Context) {

  cartProductsCount, err := getCartProductsCount(c, w.queries)
  if err != nil {
    panic(err)
  }

	c.HTML(http.StatusOK, "web/about.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "about",
		},
		"message": "About page",
		"isLoggedIn": c.GetUint64("user_id") > 0,
    "cartProductsCount": cartProductsCount,
	})
}

func (w WebController) RenderCommentsPage(c *gin.Context) {
  cartProductsCount, err := getCartProductsCount(c, w.queries)
  if err != nil {
    panic(err)
  }

	c.HTML(http.StatusOK, "web/comments.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "comments",
		},
		"message": "Comments page",
		"isLoggedIn": c.GetUint64("user_id") > 0,
    "cartProductsCount": cartProductsCount,
	})
}

func (w WebController) RenderDeliveryPage(c *gin.Context) {
  cartProductsCount, err := getCartProductsCount(c, w.queries)
  if err != nil {
    panic(err)
  }

	c.HTML(http.StatusOK, "web/delivery.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "delivery",
		},
		"message": "Delivery page",
		"isLoggedIn": c.GetUint64("user_id") > 0,
    "cartProductsCount": cartProductsCount,
	})
}

func (w WebController) RenderContactsPage(c *gin.Context) {
  cartProductsCount, err := getCartProductsCount(c, w.queries)
  if err != nil {
    panic(err)
  }

	c.HTML(http.StatusOK, "web/contacts.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "contacts",
		},
		"message": "Contacts page",
		"isLoggedIn": c.GetUint64("user_id") > 0,
    "cartProductsCount": cartProductsCount,
	})
}

type CategoryPage struct {
	Uri string `uri:"slug" binding:"required"`
}

type CategoryFilters struct {
	MinPrice     int32  `form:"min_price"`
	MaxPrice     int32  `form:"max_price"`
	LabelID      int32  `form:"label_id"`
	PriceSorting string `form:"price_sorting"`
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
		CategoryID:   category.ID,
		MinPrice:     query.MinPrice,
		MaxPrice:     query.MaxPrice,
		LabelID:      query.LabelID,
		PriceSorting: query.PriceSorting,
	})
	if err != nil {
		panic(err)
	}

	labels, err := w.queries.ListLabels(c, db.ListLabelsParams{
		Limit:  99999,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}

  cartProductsCount, err := getCartProductsCount(c, w.queries)
  if err != nil {
    panic(err)
  }

	c.HTML(http.StatusOK, "web/category.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "catalogue",
		},
		"isLoggedIn":   c.GetUint64("user_id") > 0,
    "cartProductsCount": cartProductsCount,
		"categoryName": category.Name,
		"products":     products,
		"labels":       labels,
		"slug":         page.Uri,
	})
}

func (w WebController) RenderProducts(c *gin.Context) {
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
		CategoryID:   category.ID,
		MinPrice:     query.MinPrice,
		MaxPrice:     query.MaxPrice,
		LabelID:      query.LabelID,
		PriceSorting: query.PriceSorting,
	})
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "components/products.html", gin.H{
		"products":     products,
		"priceSorting": query.PriceSorting,
		"slug":         page.Uri,
	})
}

type ProductPage struct {
	Category  string `uri:"slug" binding:"required"`
	ProductID int32  `uri:"id" binding:"required"`
}

func (w WebController) RenderProductPage(c *gin.Context) {
	var productPage ProductPage

	if err := c.ShouldBindUri(&productPage); err != nil {
		panic(err)
	}

	product, err := w.queries.GetProduct(c, productPage.ProductID)
	if err != nil {
		panic(err)
	}

	isProductInCart := false

  cartProductsCount, err := getCartProductsCount(c, w.queries)
  if err != nil {
    panic(err)
  }

	respondWithHTML := func() {
		c.HTML(http.StatusOK, "web/product.html", gin.H{
			"pages": PagesInfo{
				Pages:       pages,
				CurrentPage: "catalogue",
			},
			"isLoggedIn":      c.GetUint64("user_id") > 0,
      "cartProductsCount": cartProductsCount,
			"product":         product,
			"isProductInCart": isProductInCart,
		})
	}

	userId := helpers.GetSession(c)
	if userId != 0 {
		session, err := w.queries.GetShoppingSessionByUserId(c, userId)
		if err == sql.ErrNoRows {
			respondWithHTML()
			return
		} else if err != nil {
			panic(err)
		}

		cartProducts, err := w.queries.GetProdutsInCart(c, session.ID)
		if err == sql.ErrNoRows {
			respondWithHTML()
			return
		} else if err != nil {
			panic(err)
		}

		if slices.ContainsFunc(cartProducts, func(item db.GetProdutsInCartRow) bool {
			return item.ID == product.ID
		}) {
			isProductInCart = true
		}
	}

	respondWithHTML()
}

type CartItemParams struct {
	ProductID int32 `form:"id" binding:"required"`
	Quantity  int32 `form:"quantity" binding:"required"`
}

func (w WebController) AddItemToCart(c *gin.Context) {
	var cartItemParams CartItemParams

	if err := c.ShouldBind(&cartItemParams); err != nil {
		panic(err)
	}

	userId := c.GetUint64("user_id")
	if userId == 0 {
		panic("unauthorized")
	}

  tx, err := w.db.Begin()
  if err != nil {
    panic(err)
  }
  defer tx.Rollback()
  qtx := w.queries.WithTx(tx)

	var session db.ShoppingSession
	session, err = qtx.GetShoppingSessionByUserId(c, int32(userId))
	if err != nil {
		if err == sql.ErrNoRows {
			session, err = qtx.CreateShoppingSession(c, int32(userId))
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	_, err = qtx.CreateCartItem(c, db.CreateCartItemParams{
		SessionID: session.ID,
		ProductID: cartItemParams.ProductID,
		Quantity:  cartItemParams.Quantity,
	})

	if err != nil {
		panic(err)
	}

  // update cart total
  RecalculateCartTotal(c, qtx, int32(userId), session.ID)

  cartProductsCount, err := getCartProductsCount(c, qtx)
  if err != nil {
    panic(err)
  }

  if err = tx.Commit(); err != nil {
    panic(err)
  }

	c.HTML(http.StatusOK, "components/createdCartItem.html", gin.H{
    "CartProductsCount": cartProductsCount,
  })
}

func (w WebController) RenderCartPage(c *gin.Context) {
	var cartProducts []db.GetProdutsInCartRow
  var session db.ShoppingSession


  cartProductsCount, err := getCartProductsCount(c, w.queries)
  if err != nil {
    panic(err)
  }

	respondWithHTML := func() {
		c.HTML(http.StatusOK, "web/cart.html", gin.H{
			"isLoggedIn": c.GetUint64("user_id") > 0,
      "cartProductsCount": cartProductsCount,
			"products":   cartProducts,
      "session": session,
		})
	}

	userID := helpers.GetSession(c)
	if userID == 0 {
		respondWithHTML()
		return
	}

	session, err = w.queries.GetShoppingSessionByUserId(c, userID)
	if err == sql.ErrNoRows {
		respondWithHTML()
		return
	} else if err != nil {
		panic(err)
	}

	cartProducts, err = w.queries.GetProdutsInCart(c, session.ID)
	if err == sql.ErrNoRows {
		respondWithHTML()
		return
	} else if err != nil {
    panic(err)
  }

  respondWithHTML()
}

type CartItemManipulation struct {
  ProductID int32 `form:"product_id" binding:"required"`
}

func RecalculateCartTotal(c *gin.Context, q *db.Queries, userID, sessionID int32) db.ShoppingSession {
  products, err := q.GetProdutsInCart(c, sessionID)
  if err != nil {
    panic(err)
  }

  var cartTotal int32 = 0
  for _, item := range products {
    cartTotal += item.PriceInt * item.Quantity
  }

  updatedSession, err := q.UpdateSessionTotal(c, db.UpdateSessionTotalParams{
    UserID: userID,
    TotalInt: sql.NullInt32{
      Valid: true,
      Int32: cartTotal,
    },
  })
  if err != nil {
    panic(err)
  }

  return updatedSession
}

func (w WebController) ManupulateQuantity(c *gin.Context, operation string) {
  var params CartItemManipulation
  if err := c.ShouldBind(&params); err != nil {
    panic(err)
  }

  userID := helpers.GetSession(c)
  if userID == 0 {
    panic("user not authorized")
  }

  tx, err := w.db.Begin()
  if err != nil {
    panic(err)
  }
  defer tx.Rollback()
  qtx := w.queries.WithTx(tx)

  session, err := qtx.GetShoppingSessionByUserId(c, userID)
	if err != nil {
		panic(err)
	}

  cartItem, err := qtx.GetCartItem(c, db.GetCartItemParams{
    SessionID: session.ID,
    ProductID: params.ProductID,
  })

  targetQuantity := cartItem.Quantity
  if operation == "inc" {
    targetQuantity++
  } else if targetQuantity > 1 {
    targetQuantity--
  }

  updatedItem, err := qtx.UpdateCartItemQuantity(c, db.UpdateCartItemQuantityParams{
    ID: cartItem.ID,
    Quantity: targetQuantity,
  })
  if err != nil {
    panic(err)
  }

  product ,err := qtx.GetProduct(c, params.ProductID)
  if err != nil {
    panic(err)
  }

  updatedSession := RecalculateCartTotal(c, qtx, userID, session.ID)

  if err = tx.Commit(); err != nil {
    panic(err)
  }

  c.HTML(http.StatusOK, "components/quantity.html", gin.H{
    "ID": params.ProductID,
    "Quantity": updatedItem.Quantity,
    "Total": product.PriceInt * targetQuantity,
    "CartTotal": updatedSession.TotalInt.Int32,
  })
}

func (w WebController) DecrementQuantity(c *gin.Context) {
  w.ManupulateQuantity(c, "dec")
}

func (w WebController) IncrementQuantity(c *gin.Context) {
  w.ManupulateQuantity(c, "inc")
}

type DeleteCartItem struct {
  ProductID int32 `uri:"product_id" binding:"required"`
}

func (w WebController) DeleteCartItem(c *gin.Context) {
  var params DeleteCartItem
  c.ShouldBindUri(&params)

  tx, err := w.db.Begin()
  if err != nil {
    panic(err)
  }
  defer tx.Rollback()
  qtx := w.queries.WithTx(tx)

  userID := helpers.GetSession(c)
  if userID == 0 {
    panic("user not authorized")
  }

  session, err := qtx.GetShoppingSessionByUserId(c, userID)
	if err != nil {
		panic(err)
	}

  if _, err := qtx.DeleteCartItem(c, db.DeleteCartItemParams{
    SessionID: session.ID,
    ProductID: params.ProductID,
  }); err != nil {
    panic(err)
  }

  // update cart total
  updatedSession := RecalculateCartTotal(c, qtx, userID, session.ID)

  cartProductsCount, err := getCartProductsCount(c, qtx)
  if err != nil {
    panic(err)
  }

  if err = tx.Commit(); err != nil {
    panic(err)
  }

  if updatedSession.TotalInt.Int32 == 0 {
    c.Header("HX-Refresh", "true")
  }

  c.HTML(http.StatusOK, "components/delete_cart_item.html", gin.H{
    "CartTotal": updatedSession.TotalInt.Int32,
    "CartProductsCount": cartProductsCount,
  })
}

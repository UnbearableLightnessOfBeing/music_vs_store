package web

import (
	"database/sql"
	db "music_vs_store/db/sqlc"
	"music_vs_store/helpers"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebController struct {
	queries *db.Queries
	db      *sql.DB
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

func (w WebController) RenderMainPage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
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
		"categories":        categories,
		"isLoggedIn":        c.GetUint64("user_id") > 0,
		"cartProductsCount": cartProductsCount,
	})
}

func (w WebController) RenderCataloguePage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
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
		"categories":        categories,
		"isLoggedIn":        c.GetUint64("user_id") > 0,
		"cartProductsCount": cartProductsCount,
	})
}

func (w WebController) RenderAboutPage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	cartProductsCount, err := getCartProductsCount(c, w.queries)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/about.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "about",
		},
		"message":           "About page",
		"isLoggedIn":        c.GetUint64("user_id") > 0,
		"cartProductsCount": cartProductsCount,
		"categories":        categories,
	})
}

func (w WebController) RenderCommentsPage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	cartProductsCount, err := getCartProductsCount(c, w.queries)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/comments.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "comments",
		},
		"message":           "Comments page",
		"isLoggedIn":        c.GetUint64("user_id") > 0,
		"cartProductsCount": cartProductsCount,
		"categories":        categories,
	})
}

func (w WebController) RenderDeliveryPage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	cartProductsCount, err := getCartProductsCount(c, w.queries)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/delivery.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "delivery",
		},
		"message":           "Delivery page",
		"isLoggedIn":        c.GetUint64("user_id") > 0,
		"cartProductsCount": cartProductsCount,
		"categories":        categories,
	})
}

func (w WebController) RenderContactsPage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	cartProductsCount, err := getCartProductsCount(c, w.queries)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/contacts.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "contacts",
		},
		"message":           "Contacts page",
		"isLoggedIn":        c.GetUint64("user_id") > 0,
		"cartProductsCount": cartProductsCount,
		"categories":        categories,
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
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

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
		"isLoggedIn":        c.GetUint64("user_id") > 0,
		"cartProductsCount": cartProductsCount,
		"categoryName":      category.Name,
		"categories":        categories,
		"products":          products,
		"labels":            labels,
		"slug":              page.Uri,
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
	ProductID int32  `uri:"id" binding:"required"`
}

func (w WebController) RenderProductPage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	var productPage ProductPage

	if err := c.ShouldBindUri(&productPage); err != nil {
		panic(err)
	}

	product, err := w.queries.GetProduct(c, productPage.ProductID)
	if err != nil {
		panic(err)
	}

	labelName := ""
	if product.LabelID.Valid {
		productLabel, err := w.queries.GetLabel(c, product.LabelID.Int32)
		if err == sql.ErrNoRows {
			labelName = ""
		} else if err != nil {
			panic(err)
		} else {
			labelName = productLabel.Name
		}
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
				CurrentPage: "product",
			},
			"isLoggedIn":        c.GetUint64("user_id") > 0,
			"cartProductsCount": cartProductsCount,
			"categories":        categories,
			"product":           product,
			"isProductInCart":   isProductInCart,
			"labelName":         labelName,
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
	recalculateCartTotal(c, qtx, int32(userId), session.ID)

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
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	var cartProducts []db.GetProdutsInCartRow
	var session db.ShoppingSession

	cartProductsCount, err := getCartProductsCount(c, w.queries)
	if err != nil {
		panic(err)
	}

	respondWithHTML := func() {
		c.HTML(http.StatusOK, "web/cart.html", gin.H{
			"pages": PagesInfo{
				Pages:       pages,
				CurrentPage: "cart",
			},
			"isLoggedIn":        c.GetUint64("user_id") > 0,
			"cartProductsCount": cartProductsCount,
			"categories":        categories,
			"products":          cartProducts,
			"session":           session,
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
		ID:       cartItem.ID,
		Quantity: targetQuantity,
	})
	if err != nil {
		panic(err)
	}

	product, err := qtx.GetProduct(c, params.ProductID)
	if err != nil {
		panic(err)
	}

	updatedSession := recalculateCartTotal(c, qtx, userID, session.ID)

	if err = tx.Commit(); err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "components/quantity.html", gin.H{
		"ID":        params.ProductID,
		"Quantity":  updatedItem.Quantity,
		"Total":     product.PriceInt * targetQuantity,
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
	updatedSession := recalculateCartTotal(c, qtx, userID, session.ID)

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
		"CartTotal":         updatedSession.TotalInt.Int32,
		"CartProductsCount": cartProductsCount,
	})
}

type FormField struct {
	Name   string
	IDName string
	Type   string
}

var regionFields []FormField = []FormField{
	{Name: "Регион", IDName: "district", Type: "text"},
	{Name: "Город", IDName: "city", Type: "text"},
	{Name: "Почтовый индекс", IDName: "postal_code", Type: "number"},
}

var personalInfoFields []FormField = []FormField{
	{Name: "Фамилия", IDName: "customer_lastname", Type: "text"},
	{Name: "Имя", IDName: "customer_firstname", Type: "text"},
	{Name: "Отчество", IDName: "customer_middlename", Type: "text"},
	{Name: "Телефон", IDName: "customer_phone_number", Type: "tel"},
	{Name: "Email", IDName: "customer_email", Type: "email"},
	{Name: "Адрес", IDName: "customer_address", Type: "text"},
}

func (w WebController) RenderCheckoutPage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	cartProductsCount, err := getCartProductsCount(c, w.queries)
	if err != nil {
		panic(err)
	}

	if cartProductsCount == 0 {
		c.Redirect(http.StatusMovedPermanently, "/cart")
		return
	}

	userID := helpers.GetSession(c)

	session, err := w.queries.GetShoppingSessionByUserId(c, userID)
	if err == sql.ErrNoRows {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	} else if err != nil {
		panic(err)
	}

	totalPrice := session.TotalInt.Int32
	var deliveryPrice int32 = 500
	if totalPrice >= 10000 {
		deliveryPrice = 0
	}
	finalTotal := totalPrice + deliveryPrice

	countries, err := w.queries.ListCountries(c, db.ListCountriesParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}

	deliveryMethods, err := w.queries.ListDeliveryMethods(c, db.ListDeliveryMethodsParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}
	paymentMethods, err := w.queries.ListPaymentMethods(c, db.ListPaymentMethodsParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/checkout.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "",
		},
		"isLoggedIn":         c.GetUint64("user_id") > 0,
		"cartProductsCount":  cartProductsCount,
		"categories":         categories,
		"totalPrice":         totalPrice,
		"deliveryPrice":      deliveryPrice,
		"finalTotal":         finalTotal,
		"userID":             userID,
		"countries":          countries,
		"regionFields":       regionFields,
		"deliveryMethods":    deliveryMethods,
		"paymentMethods":     paymentMethods,
		"personalInfoFields": personalInfoFields,
	})
}

type OrderParams struct {
	UserID              int32  `form:"user_id"`
	ProductCount        int32  `form:"product_count"`
	PriceInt            int32  `form:"price_int"`
	DeliveryPriceInt    int32  `form:"delivery_price_int"`
	TotalInt            int32  `form:"total_int"`
	CountryID           int32  `form:"country_id"`
	District            string `form:"district"`
	City                string `form:"city"`
	PostalCode          int32  `form:"postal_code"`
	DeliveryMethodID    int32  `form:"delivery_method_id"`
	PaymentMethodID     int32  `form:"payment_method_id"`
	CustomerFirstname   string `form:"customer_firstname"`
	CusotmerMiddlename  string `form:"customer_middlename"`
	CustomerLastname    string `form:"customer_lastname"`
	CustomerPhoneNumber string `form:"customer_phone_number"`
	CustomerEmail       string `form:"customer_email"`
	CustomerAddress     string `form:"customer_address"`
	CustomerComment     sql.NullString
}

type CommentParam struct {
	CustomerComment *string `form:"customer_comment"`
}

func (w WebController) CreateOrder(c *gin.Context) {
	var params OrderParams
	if err := c.ShouldBind(&params); err != nil {
		panic(err)
	}

	if params.District == "" ||
		params.City == "" ||
		params.PostalCode == 0 ||
		params.CustomerFirstname == "" ||
		params.CustomerLastname == "" ||
		params.CusotmerMiddlename == "" ||
		params.CustomerPhoneNumber == "" ||
		params.CustomerEmail == "" ||
		params.CustomerAddress == "" {
		c.Header("HX-Retarget", "#submit-error")
		c.HTML(http.StatusBadRequest, "components/create_order_error.html", gin.H{
			"message": "Проверте все ли поля заполнены",
		})
		return
	}

	var commentParam CommentParam
	if err := c.ShouldBind(&commentParam); err != nil {
		panic(err)
	}

	if *commentParam.CustomerComment == "" || commentParam.CustomerComment == nil {
		params.CustomerComment = sql.NullString{
			Valid:  false,
			String: "",
		}
	} else {
		params.CustomerComment = sql.NullString{
			Valid:  true,
			String: *commentParam.CustomerComment,
		}
	}

	session, err := w.queries.GetShoppingSessionByUserId(c, params.UserID)
	if err != nil {
		panic(err)
	}

	cartProducts, err := w.queries.GetProdutsInCart(c, session.ID)
	if err != nil {
		panic(err)
	}
	if len(cartProducts) == 0 {
		panic("no products in cart")
	}

	tx, err := w.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()
	qtx := w.queries.WithTx(tx)

	// create new order
	newOrderID, err := qtx.CreateOrder(c, db.CreateOrderParams(params))
	if err != nil {
		panic(err)
	}

	// add products_order relationships
	for _, cartProduct := range cartProducts {
		_, err := qtx.AddProductToOrder(c, db.AddProductToOrderParams{
			OrderID:   newOrderID,
			ProductID: cartProduct.ID,
			Count:     cartProduct.Quantity,
		})
		if err != nil {
			panic(err)
		}
	}

	// delete session
	_, err = qtx.DeleteSessionByUserId(c, params.UserID)
	if err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "components/order_created.html", gin.H{
		"orderID": newOrderID,
	})
}

func createdTimeFormatted(time string) string {
	split := strings.Split(time, " ")
	if len(split) < 2 {
		return ""
	}

	date := split[0]
	splitDate := strings.Split(date, "-")
	if len(splitDate) < 3 {
		date = ""
	} else {
		slices.Reverse(splitDate)
		date = strings.Join(splitDate, ".")
	}

	dateTime := split[1]
	splitDateTime := strings.Split(dateTime, ".")
	if len(splitDateTime) < 1 {
		dateTime = ""
	} else {
		dateTime = splitDateTime[0]
	}

	return date + " " + dateTime
}

func (w WebController) RenderOrdersPage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	userID := helpers.GetSession(c)
	if userID == 0 {
		c.Redirect(http.StatusMovedPermanently, "/")
	}

	cartProductsCount, err := getCartProductsCount(c, w.queries)
	if err != nil {
		panic(err)
	}

	orders, err := w.queries.GetOrdersByUserId(c, userID)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/orders.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "orders",
		},
		"isLoggedIn":        c.GetUint64("user_id") > 0,
		"cartProductsCount": cartProductsCount,
		"categories":        categories,
		"orders":            orders,
		"ordersCount":       len(orders),
	})
}

type OrderPage struct {
	OrderID int32 `uri:"id" binding:"required"`
}

func (w WebController) RenderOrderPage(c *gin.Context) {
	var orderPage OrderPage
	if err := c.ShouldBindUri(&orderPage); err != nil {
		panic(err)
	}

	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	userID := helpers.GetSession(c)
	if userID == 0 {
		c.Redirect(http.StatusMovedPermanently, "/")
	}

	cartProductsCount, err := getCartProductsCount(c, w.queries)
	if err != nil {
		panic(err)
	}

	order, err := w.queries.GetOrder(c, orderPage.OrderID)
	if err != nil {
		panic(err)
	}

	orderProducts, err := w.queries.GetOrderProducts(c, order.ID)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/order.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "orders",
		},
		"isLoggedIn":        c.GetUint64("user_id") > 0,
		"cartProductsCount": cartProductsCount,
		"categories":        categories,
		"order":             order,
		"orderProducts":     orderProducts,
	})
}

func (w WebController) RenderLoginPage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/login.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "",
		},
		"isLoggedIn": c.GetUint64("user_id") > 0,
		"categories": categories,
	})
}

func (w WebController) RenderSignupPage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/signup.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "signup",
		},
		"isLoggedIn": c.GetUint64("user_id") > 0,
		"categories": categories,
	})
}

func (w WebController) RenderSearchPage(c *gin.Context) {
	categories, err := getCategories(c, w.queries)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "web/search.html", gin.H{
		"pages": PagesInfo{
			Pages:       pages,
			CurrentPage: "search",
		},
		"isLoggedIn": c.GetUint64("user_id") > 0,
		"categories": categories,
	})
}

type SearchParams struct {
	Query string `form:"search"`
}

func (w WebController) SearchItems(c *gin.Context) {
	var params SearchParams
	if err := c.ShouldBind(&params); err != nil {
		panic(err)
	}

	if params.Query == "" {
		c.HTML(http.StatusOK, "components/search_products.html", gin.H{
			"isInputEmpty": true,
		})
		return
	}

	results, err := w.queries.SearchProducts(c, "%"+strings.ToLower(params.Query)+"%")
	if err == sql.ErrNoRows {
		c.HTML(http.StatusOK, "components/search_products.html", gin.H{})
		return
	} else if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "components/search_products.html", gin.H{
		"products": results,
	})
}

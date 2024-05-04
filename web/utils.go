package web

import (
	"database/sql"
	db "music_vs_store/db/sqlc"
	"music_vs_store/helpers"

	"github.com/gin-gonic/gin"
)

func getCategories(c *gin.Context, q *db.Queries) ([]db.Category, error) {
	categories, err := q.ListCategories(c, db.ListCategoriesParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}
	return categories, nil
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

func recalculateCartTotal(c *gin.Context, q *db.Queries, userID, sessionID int32) db.ShoppingSession {
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

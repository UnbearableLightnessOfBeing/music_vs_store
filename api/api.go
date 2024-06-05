package api

import (
	"database/sql"
	db "music_vs_store/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)


type ApiController struct {
	queries *db.Queries
	db      *sql.DB
}

func NewApiController(queries *db.Queries, db *sql.DB) *ApiController {
	return &ApiController{
		queries,
		db,
	}
}

func (w *ApiController) Users(c *gin.Context) {
  users, err := w.queries.ListUsers(c, db.ListUsersParams{
    Limit: 999,
    Offset: 0,
  })
  if err != nil {
    panic(err)
  }

  c.JSON(http.StatusOK, gin.H{
    "users": users,
  }) 
}

type ToggleUser struct {
	UserID int32 `uri:"id" binding:"required"`
}

func (w *ApiController) UserToggleIsAdmin(c *gin.Context) {
	var userToAdmin ToggleUser
	if err := c.ShouldBindUri(&userToAdmin); err != nil {
		panic(err)
	}

  user, err := w.queries.GetUser(c, userToAdmin.UserID)
  if err != nil {
    panic(err)
  }

  isAdmin := user.IsAdmin.Bool

  user, err = w.queries.UpdateUserIsAdmin(c, db.UpdateUserIsAdminParams{
    ID: userToAdmin.UserID,
    IsAdmin: sql.NullBool{
      Bool: !isAdmin,
      Valid: true,
    },
  })

  c.JSON(http.StatusOK, gin.H{
    "result": user,
  })
}

func (w *ApiController) Products(c *gin.Context) {
  products, err := w.queries.ListProducts(c, db.ListProductsParams{
    Limit: 999,
    Offset: 0,
  })      
  if err != nil {
    panic(err)
  }

  c.JSON(http.StatusOK, gin.H{
    "products": products,
  })
}

type Product struct {
	ProductID int32 `uri:"id" binding:"required"`
}

func (w *ApiController) Product(c *gin.Context) {
  var productPage Product
  err := c.ShouldBindUri(&productPage)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Bad url",
    })
  }

  product, err := w.queries.GetProduct(c, productPage.ProductID)
  if err == sql.ErrNoRows {
    c.JSON(http.StatusNotFound, gin.H{
      "message": "Product with such id doesn't exist",
    })
    return
  } else if err != nil {
    panic(err)
  }

  c.JSON(http.StatusOK, gin.H{
    "product": product,
  })
}

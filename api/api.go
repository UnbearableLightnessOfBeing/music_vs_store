package api

import (
	"database/sql"
	"fmt"
	db "music_vs_store/db/sqlc"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		Limit:  999,
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
			Bool:  !isAdmin,
			Valid: true,
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"result": user,
	})
}

func (w *ApiController) Products(c *gin.Context) {
	products, err := w.queries.ListProducts(c, db.ListProductsParams{
		Limit:  999,
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

// ----PRODUCTS----
type CreateProductReq struct {
	Name            string         `json:"name"`
	PriceInt        int32          `json:"price_int"`
	LabelID         sql.NullInt32  `json:"label_id"`
	Description     sql.NullString `json:"description"`
	Characteristics sql.NullString `json:"characteristics"`
	InStock         bool           `json:"in_stock"`
}

func parseRequest(c *gin.Context) (*CreateProductReq, error) {
	var productData CreateProductReq
	if err := c.ShouldBindJSON(&productData); err != nil {
		return nil, fmt.Errorf("Couldn't parse request")
	}
  return &productData, nil
}

func (w *ApiController) getProductByName(c *gin.Context, name string ) (db.Product, bool) {
  product, err := w.queries.GetProductByName(c, name)
  if err == sql.ErrNoRows {
    return db.Product{}, false
  }
  return product, true
} 

func (w *ApiController) CreateProduct(c *gin.Context) {
  productData, err := parseRequest(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Поля введены некорректно",
    })
    return
  }

  if _, exist := w.getProductByName(c, productData.Name); exist {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Товар с таким названием уже существует",
    })
    return
  }

	product, err := w.queries.CreateProduct(c, db.CreateProductParams(*productData))
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": err.Error(),
    })
    return
  }

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

type UpdateProductReq struct {
  ProductID int32 `uri:"id"`
}

func (w *ApiController) parseProductID(c *gin.Context) (int32, error) {
  var updateProductReq UpdateProductReq
  if err := c.ShouldBindUri(&updateProductReq); err != nil {
    return 0, err
  } else {
    return updateProductReq.ProductID, nil
  }
}

func (w *ApiController) UpdateProduct(c *gin.Context) {
  productData, err := parseRequest(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Поля введены некорректно",
    })
    return
  }

  productID, err := w.parseProductID(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Неверный url",
    })
    return
  }

  prod, exist := w.getProductByName(c, productData.Name)
  if exist && prod.ID != productID {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Товар с таким названием уже существует",
    })
    return
  }

	product, err := w.queries.UpdateProduct(
    c, 
    db.UpdateProductParams(db.UpdateProductParams{
      ID: productID,
      Name: productData.Name,
      PriceInt: productData.PriceInt,
      LabelID: productData.LabelID,
      Description: productData.Description,
      Characteristics: productData.Characteristics,
      InStock: productData.InStock,
    }))

  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": err.Error(),
    })
    return
  }

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (w *ApiController) AddImageToProdut(c *gin.Context) {
  productID, err := w.parseProductID(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Неверный url",
    })
    return
  }

  file, err := c.FormFile("image")
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": err.Error(),
    })
    return
  }
  ext := filepath.Ext(file.Filename)
  uuidStr := uuid.New().String()

  dst := "./storage/images/" + uuidStr + ext
  if err = c.SaveUploadedFile(file, dst); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
    return
  }

  imgUrl := "/storage/" + uuidStr + ext

  err = w.queries.AddImageToProduct(c, db.AddImageToProductParams{
    ID: productID,
    ArrayAppend: imgUrl,
  })  
}

type RemoveImageFromProductReq struct {
  ImageName string `json:"image_name"`
}

func (w *ApiController) parseImagePath(c *gin.Context) (string, error) {
  var req RemoveImageFromProductReq
  if err := c.ShouldBindJSON(&req); err != nil {
    return "", err
  } else {
    return req.ImageName, nil
  }
}

func (w *ApiController) RemoveImageFromProduct(c *gin.Context) {
  productID, err := w.parseProductID(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Неверный url",
    })
    return
  }

  filePath, err := w.parseImagePath(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": err.Error(),
    })
    return
  }

  splitFilePath := strings.Split(filePath, "/")
  fileName := splitFilePath[len(splitFilePath) - 1]
  fileUrl := "./storage/images/" + fileName

  if err := os.Remove(fileUrl); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
    return
  }

  err = w.queries.RemoveImageFromProduct(c, db.RemoveImageFromProductParams{
    ID: productID,
    ArrayRemove: filePath,
  })
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": err.Error(),
    })
    return
  }

  c.Status(http.StatusOK)
}

func (w *ApiController) DeleteProduct(c *gin.Context) {
  productID, err := w.parseProductID(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": err.Error(),
    })
    return
  }
  deleted, err := w.queries.DeleteProduct(c, productID)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "deleted": deleted,
  })
}

// -----CATEGORIES------
func (w *ApiController) Categories(c *gin.Context) {
  categories, err := w.queries.ListCategories(c, db.ListCategoriesParams{
    Limit: 999,
    Offset: 0,
  })   
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "categories": categories,
  })
}


type DeleteCategoryReq struct {
  ProductID int32 `uri:"id"`
}

func (w *ApiController) parseCategoryID(c *gin.Context) (int32, error) {
  var uriParams DeleteCategoryReq
  if err := c.ShouldBindUri(&uriParams); err != nil {
    return 0, err
  }
  return uriParams.ProductID, nil
}

func (w *ApiController) Category(c *gin.Context) {
  ctgrID, err := w.parseCategoryID(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": err.Error(),
    })
    return
  }

  ctgr, err := w.queries.GetCategory(c, ctgrID)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "category": ctgr,
  })
}

func (w *ApiController) parseCategoryReq(c *gin.Context) (db.CreateCategoryParams, error) {
  var createReq db.CreateCategoryParams
  if err := c.ShouldBindJSON(&createReq); err != nil {
    return db.CreateCategoryParams{}, err
  }
  return createReq, nil
}

func (w *ApiController) getCategoryByName(c *gin.Context, name string) (db.Category, bool) {
  ctgr, err := w.queries.GetCategoryByName(c, name)
  if err == sql.ErrNoRows {
    return db.Category{}, false
  } else if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
  }
  return ctgr, true
}

func (w *ApiController) getCategoryBySlug(c *gin.Context, slug string) (db.Category, bool) {
  ctgr, err := w.queries.GetCategoryBySlug(c, slug)
  if err == sql.ErrNoRows {
    return db.Category{}, false
  } else if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
  }
  return ctgr, true
}

func (w *ApiController) CreateCategory(c *gin.Context) {
  req, err := w.parseCategoryReq(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Поля заполнены некорректно",
    })
    return
  }

  _, exists := w.getCategoryByName(c, req.Name)
  if exists {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Категория с таким именем уже существует",
    })
    return
  }

  _, exists = w.getCategoryBySlug(c, req.Slug)
  if exists {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Категория с таким slug уже существует",
    })
    return
  }

  category, err := w.queries.CreateCategory(c, req)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "created": category,
  })
}

func (w *ApiController) DeleteCategory(c *gin.Context) {
  ctgrID, err := w.parseCategoryID(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": err.Error(),
    })
    return
  }

  deleted, err := w.queries.DeleteCategory(c, ctgrID)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "deleted": deleted,
  })
}

func (w *ApiController) UpdateCategory(c *gin.Context) {
  ctgrID, err := w.parseCategoryID(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": err.Error(),
    })
    return
  }

  req, err := w.parseCategoryReq(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Поля заполнены некорректно",
    })
    return
  }

  ctgr, exists := w.getCategoryByName(c, req.Name)
  if exists && ctgr.ID != ctgrID {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Категория с таким именем уже существует",
    })
    return
  }

  ctgr, exists = w.getCategoryBySlug(c, req.Slug)
  if exists && ctgr.ID != ctgrID {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Категория с таким slug уже существует",
    })
    return
  }

  category, err := w.queries.UpdateCategory(c, db.UpdateCategoryParams{
    ID: ctgrID,
    Name: req.Name,
    Slug: req.Slug,
  })
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "created": category,
  })
}

func (w *ApiController) SetCategoryImage(c *gin.Context) {
  ctgrID, err := w.parseCategoryID(c)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": "Неверный url",
    })
    return
  }

  file, err := c.FormFile("image")
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": err.Error(),
    })
    return
  }
  ext := filepath.Ext(file.Filename)
  uuidStr := uuid.New().String()

  dst := "./storage/images/" + uuidStr + ext
  if err = c.SaveUploadedFile(file, dst); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
    return
  }

  imgUrl := "/storage/" + uuidStr + ext

  category, err := w.queries.SetCategoryImage(c, db.SetCategoryImageParams{
    ID: ctgrID,
    ImgUrl: sql.NullString{
      String: imgUrl,
      Valid: true,
    },
  })  
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "updated": category,
  })
}

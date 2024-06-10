package api

import (
	"database/sql"
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

func (w *ApiController) parseAndStoreOSImage(c *gin.Context) (string, error, int) {
	file, err := c.FormFile("image")
	if err != nil {
		return "", err, http.StatusBadRequest
	}
	ext := filepath.Ext(file.Filename)
	uuidStr := uuid.New().String()

	dst := "./storage/images/" + uuidStr + ext
	if err = c.SaveUploadedFile(file, dst); err != nil {
		return "", err, http.StatusInternalServerError
	}

	imgUrl := "/storage/" + uuidStr + ext
	return imgUrl, nil, 0
}

func (w *ApiController) removeOSImage(filePath string) error {
	splitFilePath := strings.Split(filePath, "/")
	fileName := splitFilePath[len(splitFilePath)-1]
	fileUrl := "./storage/images/" + fileName
	if err := os.Remove(fileUrl); err != nil {
		return err
	}
	return nil
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
	products, err := w.queries.GetProductsWithCategory(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
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
		return
	}

	product, err := w.queries.GetProductWithCategory(c, productPage.ProductID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product with such id doesn't exist",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// ----PRODUCTS----
type CreateProductReq struct {
	Name            string `json:"name"`
	PriceInt        int32  `json:"price_int"`
	LabelID         int32  `json:"label_id"`
	Description     string `json:"description"`
	Characteristics string `json:"characteristics"`
	InStock         bool   `json:"in_stock"`
	CategoryID      int32  `json:"category_id"`
}

func parseRequest(c *gin.Context) (*CreateProductReq, error) {
	var productData CreateProductReq
	if err := c.ShouldBindJSON(&productData); err != nil {
		return nil, err
	}
	return &productData, nil
}

func (w *ApiController) getProductByName(c *gin.Context, name string) (db.Product, bool) {
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
			"message": "Поля введены некорректно: " + err.Error(),
		})
		return
	}

	if productData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": `Поле "Название" обязательно к заполнению`,
		})
		return
	}

	if _, exist := w.getProductByName(c, productData.Name); exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Товар с таким названием уже существует",
		})
		return
	}

	isLabelValid := productData.LabelID != 0

	createParams := db.CreateProductParams{
		Name:     productData.Name,
		PriceInt: productData.PriceInt,
		InStock:  productData.InStock,
		Description: sql.NullString{
			String: productData.Description,
			Valid:  true,
		},
		Characteristics: sql.NullString{
			String: productData.Characteristics,
			Valid:  true,
		},
		LabelID: sql.NullInt32{
			Int32: productData.LabelID,
			Valid: isLabelValid,
		},
	}

	tx, err := w.db.Begin()
	if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": err.Error(),
    })
    return
	}
	defer tx.Rollback()
	qtx := w.queries.WithTx(tx)

	product, err := qtx.CreateProduct(c, createParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

  if productData.CategoryID != 0 {
    // delete previous relations
    // deleted, err := w.queries.DeleteProductCategoryRelations(c, product.ID)
    // create a new relation
    _, err := qtx.AddProductCategoryRelation(c, db.AddProductCategoryRelationParams{
      ProductID: product.ID,
      CategoryID: productData.CategoryID,
    })
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "message": err.Error(),
      })
      return
    }
  }

  if err := tx.Commit(); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
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

	if productData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": `Поле "Название" обязательно к заполнению`,
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

	isLabelValid := productData.LabelID != 0


  tx, err := w.db.Begin()
  if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
  }
  qtx:= w.queries.WithTx(tx)
  defer tx.Rollback()

	product, err := qtx.UpdateProduct(
		c,
		db.UpdateProductParams(db.UpdateProductParams{
			ID:       productID,
			Name:     productData.Name,
			PriceInt: productData.PriceInt,
			LabelID: sql.NullInt32{
				Int32: productData.LabelID,
				Valid: isLabelValid,
			},
			Description: sql.NullString{
				String: productData.Description,
				Valid:  true,
			},
			Characteristics: sql.NullString{
				String: productData.Characteristics,
				Valid:  true,
			},
			InStock: productData.InStock,
		}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

  if productData.CategoryID != 0 {
    // delete previous relations
    _, err := qtx.DeleteProductCategoryRelations(c, product.ID)
    if err != sql.ErrNoRows && err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "message": err.Error(),
      })
      return
    }
    // create a new relation
    _, err = qtx.AddProductCategoryRelation(c, db.AddProductCategoryRelationParams{
      ProductID: product.ID,
      CategoryID: productData.CategoryID,
    })
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "message": err.Error(),
      })
      return
    }
  }
  tx.Commit()

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

	imgUrl, err, code := w.parseAndStoreOSImage(c)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
	}

	err = w.queries.AddImageToProduct(c, db.AddImageToProductParams{
		ID:          productID,
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

	w.removeOSImage(filePath)

	err = w.queries.RemoveImageFromProduct(c, db.RemoveImageFromProductParams{
		ID:          productID,
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
  tx, err := w.db.Begin()
  if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
  }

  qtx := w.queries.WithTx(tx)
  defer tx.Rollback()

	productID, err := w.parseProductID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

  _, err = qtx.DeleteProductCategoryRelations(c, productID)
  if err != sql.ErrNoRows  && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
  }

	deleted, err := qtx.DeleteProduct(c, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

  tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"deleted": deleted,
	})
}

// -----CATEGORIES------
func (w *ApiController) Categories(c *gin.Context) {
	categories, err := w.queries.ListCategories(c, db.ListCategoriesParams{
		Limit:  999,
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

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": `Поле "Название" обязательно к заполнению`,
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

  tx, err := w.db.Begin()
  if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
  }

  qtx := w.queries.WithTx(tx)
  defer tx.Rollback()

  if err := qtx.DeleteCategoryProductRelations(c, ctgrID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
  }

	deleted, err := qtx.DeleteCategory(c, ctgrID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
  tx.Commit()

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

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": `Поле "Название" обязательно к заполнению`,
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
		ID:   ctgrID,
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
		"updated": category,
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

	category, err := w.queries.GetCategory(c, ctgrID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	if category.ImgUrl.Valid {
		if err := w.removeOSImage(category.ImgUrl.String); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
	}

	imgUrl, err, code := w.parseAndStoreOSImage(c)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
	}

	category, err = w.queries.SetCategoryImage(c, db.SetCategoryImageParams{
		ID: ctgrID,
		ImgUrl: sql.NullString{
			String: imgUrl,
			Valid:  true,
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

type CreateLabelReq struct {
	Name string `json:"name"`
}

// ----LABELS----
func (w *ApiController) Labels(c *gin.Context) {
	lbls, err := w.queries.ListLabels(c, db.ListLabelsParams{
		Limit:  999,
		Offset: 0,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"labels": lbls,
	})
}

func (w *ApiController) Label(c *gin.Context) {
	var req DeleteLabelReq
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
  label, err := w.queries.GetLabel(c, req.LabelID) 
  if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
  }

  c.JSON(http.StatusOK, gin.H{
    "label": label,
  })
}

func (w *ApiController) UpdateLabel(c *gin.Context) {
	var req DeleteLabelReq
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var creaqteReq CreateLabelReq
	if err := c.ShouldBindJSON(&creaqteReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if creaqteReq.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": `Поле "Название" обязательно к заполнению`,
		})
		return
	}

  updated, err := w.queries.UpdateLabel(c, db.UpdateLabelParams{
    ID: req.LabelID,
    Name: creaqteReq.Name,
  })
  if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
  }

  c.JSON(http.StatusOK, gin.H{
    "updated": updated,
  })
}

func (w *ApiController) CreateLabel(c *gin.Context) {
	var labelReq CreateLabelReq
	if err := c.ShouldBindJSON(&labelReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if labelReq.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": `Поле "Название" обязательно к заполнению`,
		})
		return
	}

	lbl, err := w.queries.CreateLabel(c, labelReq.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"created": lbl,
	})
}

type DeleteLabelReq struct {
	LabelID int32 `uri:"id"`
}

func (w *ApiController) DeleteLabel(c *gin.Context) {
	var deleteReq DeleteLabelReq
	if err := c.ShouldBindUri(&deleteReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

  tx, err := w.db.Begin()
  if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
  }

  qtx := w.queries.WithTx(tx)
  defer tx.Rollback()

  if err = qtx.RemoveLabelProductRelations(c, sql.NullInt32{
    Int32: deleteReq.LabelID,
    Valid: true,
  }); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
  }

	lbl, err := qtx.DeleteLabel(c, deleteReq.LabelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

  if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
  }

	c.JSON(http.StatusOK, gin.H{
		"deleted": lbl,
	})
}

// ---ORDRES---
func (w *ApiController) Orders(c *gin.Context) {
	orders, err := w.queries.GetOrders(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

package controllers

import (
	db "music_vs_store/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	queries *db.Queries
}

func NewUsersController(queries *db.Queries) UsersController {
	return UsersController{queries}
}

func (q UsersController) ListUsers(c *gin.Context) {
	users, err := q.queries.ListUsers(c, db.ListUsersParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}

	categories, err := q.queries.ListCategories(c, db.ListCategoriesParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "home/index.html", gin.H{
		"title":      "Hello GinGonic",
		"users":      users,
		"categories": categories,
		"isLoggedIn": c.GetUint64("user_id") > 0,
	})
}

// func (q UsersController) CreateUserView(c *gin.Context) {
// 	c.HTML(http.StatusOK, "web/signup.html", gin.H{
// 		"title": "Create new account",
// 	})
// }

// func (q UsersController) LoginView(c *gin.Context) {
// 	categories, err := web.GetCategories(c, q.queries)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c.HTML(http.StatusOK, "home/login.html", gin.H{
// 		"isLoggedIn": c.GetUint64("user_id") > 0,
// 		"categories": categories,
// 	})
// }

func respondWithError(c *gin.Context, endpoint string, status int, message string) {
	c.HTML(status, "web/"+endpoint+".html", gin.H{
		"title":     "Create new account",
		"errorText": message,
	})
}

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
	users, err := q.queries.ListUsrs(c, db.ListUsrsParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "home/index.html", gin.H{
		"title": "Hello GinGonik",
		"users": users,
    "isLoggedIn": c.GetUint64("user_id") > 0,
	})
}

func (q UsersController) CreateUserView(c *gin.Context) {
	c.HTML(http.StatusOK, "home/signup.html", gin.H{
		"title": "Create new account",
	})
}

func (q UsersController) LoginView(c *gin.Context) {
	c.HTML(http.StatusOK, "home/login.html", gin.H{})
}

func respondWithError(c *gin.Context, endpoint string, status int, message string) {
	c.HTML(status, "home/" + endpoint + ".html", gin.H{
		"title":     "Create new account",
		"errorText": message,
	})
}

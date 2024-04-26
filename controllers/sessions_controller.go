package controllers

import (
	"database/sql"
	db "music_vs_store/db/sqlc"
	"music_vs_store/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionsController struct {
	queries *db.Queries
}

func NewSessionsController(queries *db.Queries) SessionsController {
  return SessionsController{
    queries,
  }
}

type SignupFileds struct {
	db.CreateUserParams
	RepeatPassword string `form:"repeat_password"`
}

func (q SessionsController) Signup(c *gin.Context) {
	var params SignupFileds
	c.ShouldBind(&params)

	if len(params.Username) == 0 ||
		len(params.Email) == 0 ||
		len(params.Password) == 0 {
		respondWithError(c, "signup", http.StatusBadRequest, "Fill out the form")
		return
	}

	userToCheckName, err := q.queries.GetUserByName(c, params.Username)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	if userToCheckName.Username == params.Username {
		respondWithError(c, "signup", http.StatusIMUsed, "this username is already in use")
		return
	}

	userToCheckEmail, err := q.queries.GetUserByEmail(c, params.Email)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	if userToCheckEmail.Email == params.Email {
		respondWithError(c, "signup", http.StatusIMUsed, "this email is already in use")
		return
	}

	if params.Password != params.RepeatPassword {
		respondWithError(c, "signup", http.StatusNotAcceptable, "passwords are not the same")
		return
	}

  hash, err := helpers.HashPassword(params.Password)
  if err != nil {
    panic(err)
  }

	u := db.CreateUserParams{
		Username: params.Username,
		Email:    params.Email,
		Password: hash,
	}

  created, err := q.queries.CreateUser(c, u)
	if err != nil {
		panic(err)
	}

  helpers.SetSession(c, created.ID)

	c.Redirect(http.StatusMovedPermanently, "/")
}

type LoginFields struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (q SessionsController) Login(c *gin.Context) {
	var params LoginFields
	c.ShouldBind(&params)

	user, err := q.queries.GetUserByEmail(c, params.Email)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	if len(params.Email) == 0 ||
		len(params.Password) == 0 {
		respondWithError(c, "login", http.StatusBadRequest, "Fill out the form")
		return
	}

	if !helpers.IsPasswordValid(user.Password, params.Password) {
		respondWithError(c, "login", http.StatusUnauthorized, "Wrong credentials")
		return
	}

  helpers.SetSession(c, user.ID)

  c.Redirect(http.StatusMovedPermanently, "/")
}

func (q SessionsController) Logout(c *gin.Context) {
  helpers.ClearSession(c)
  c.Header("HX-Refresh", "true") 
}

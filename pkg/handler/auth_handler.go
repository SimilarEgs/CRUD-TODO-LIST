package handler

import (
	"database/sql"
	"net/http"

	"github.com/SimilarEgs/CRUD-TODO-LIST/internal/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Hanlder) signUp(c *gin.Context) {
	// var for storing user input data
	var input entity.User

	// validate request body
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "[Error] invalid request, try again")
		return
	}

	// affter parsing and data validation, send data to the service layer via «CreateUser» method
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "[Error] operation failed, try again")
		return
	}

	// if operation was successfully done, send code 201 to the client and json with id of created user
	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})

}

// struct for parsing request body
type singInUserInput struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *Hanlder) signIn(c *gin.Context) {

	var input singInUserInput

	// validate request body
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "[Error] invalid login credentials")
		return
	}

	// generate JWT and error handling
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		// error handling for 400 -> unknown users
		if err == sql.ErrNoRows {
			newErrorResponse(c, http.StatusBadRequest, "[Error] accout with given username not found")
			return
		}
		// error handilng for 401 -> incorrect pwd
		if err == bcrypt.ErrMismatchedHashAndPassword {
			newErrorResponse(c, http.StatusUnauthorized, "[Error] invalid login credentials")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "[Error] connection error, try again")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

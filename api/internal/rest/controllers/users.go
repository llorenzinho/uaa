package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/llorenzinho/goauth/api/internal"
	"github.com/llorenzinho/goauth/api/internal/rest/dtos"
	"github.com/llorenzinho/goauth/api/internal/services"
)

type UserController struct {
	s services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{s: s}
}

func (uc *UserController) CreateUser(c *gin.Context) {

	var params dtos.CreateUserParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := uc.s.CreateUser(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := uc.s.GetUserByID(id)
	if err != nil {
		switch err {
		case internal.ErrInvalidUUID:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, user)
}

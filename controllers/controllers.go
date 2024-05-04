package controllers

import (
	"github.com/draco121/authorizationservice/core"
	"github.com/draco121/common/constants"
	"github.com/draco121/common/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controllers struct {
	service core.IAuthorizationService
}

func NewControllers(service core.IAuthorizationService) Controllers {
	c := Controllers{
		service: service,
	}
	return c
}

func (s *Controllers) Authorize(c *gin.Context) {
	var input models.AuthorizationInput
	if c.ShouldBind(&input) != nil {
		c.Status(http.StatusBadRequest)
	} else {
		output := s.service.Authorize(c.Request.Context(), &input)
		if output.Grant == constants.Allowed {
			c.JSON(http.StatusOK, output)
		} else {
			c.JSON(http.StatusOK, output)
		}
	}
}

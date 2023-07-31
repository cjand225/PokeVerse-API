// Package controller implements the HTTP handlers for the web service.
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pokeverse/web-service/service"
	"pokeverse/web-service/validation"
)

const (
	GetPokemonByIDEndpoint = "/pokemon/:lang/:id"
)

// Controller handles the HTTP routes for the web service.
type Controller struct {
	service *service.Service
}

// NewController creates a new instance of Controller with the provided service.
func NewController(service *service.Service) *Controller {
	return &Controller{service: service}
}

// SetupRoutes configures the HTTP routes for the controller using the provided router.
func (con *Controller) SetupRoutes(router *gin.Engine) {
	router.GET(GetPokemonByIDEndpoint, con.GetPokemonByIDHandler)
}

// GetPokemonByIDHandler handles the "/pokemon/:lang/:id" GET endpoint.
func (con *Controller) GetPokemonByIDHandler(c *gin.Context) {
	pid := c.Param("id")
	pLang := c.Param("lang")

	if pid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	id, err := strconv.Atoi(pid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid ID is required"})
		return
	}

	if pLang == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language code is required"})
		return
	}

	if !validation.IsValidLanguageCode(pLang) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid Language code is required"})
		return
	}

	data, err := con.service.GetPokemonByID(id, pLang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data."})
		return
	}

	c.JSON(http.StatusOK, data)
}

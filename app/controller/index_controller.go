package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// indexController index controller
type indexController struct{}

// Index action for list all swagger-JSON file list
func (s *indexController) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "list.html", map[string]string{"a": "b"})
}

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tvb-sz/serve-swagger-ui/app/service"
	"net/http"
)

// indexController index controller
type indexController struct{}

// Index action for list all swagger-JSON file list
func (s *indexController) Index(ctx *gin.Context) {
	data, err := service.ParseService.ParseWithCache()
	if err != nil {
		return
	}
	ctx.HTML(http.StatusOK, "list.html", data.Items)
}

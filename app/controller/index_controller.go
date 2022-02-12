package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/tvb-sz/serve-swagger-ui/app/service"
	"net/http"
	"os"
	"strings"
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

// Detail action for swagger-JSON file detail
func (s *indexController) Detail(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "detail.html", ctx.Query("json"))
}

// Json serve JSON file
func (s *indexController) Json(ctx *gin.Context) {
	hash := strings.TrimRight(ctx.Param("path"), ".json")
	data, err := service.ParseService.ParseWithCache()
	if err != nil {
		return
	}
	path, exist := data.Table[hash]
	if !exist {
		return
	}
	stream, err := os.ReadFile(path)
	if err != nil {
		return
	}

	ctx.DataFromReader(http.StatusOK, int64(len(stream)), "application/json", bytes.NewReader(stream), nil)
}

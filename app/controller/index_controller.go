package controller

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tvb-sz/serve-swagger-ui/app/service"
	"github.com/tvb-sz/serve-swagger-ui/render"
	"net/http"
	"os"
	"strings"
)

// indexController index controller
type indexController struct{}

var (
	SwaggerPathNotFoundFile = errors.New("the swagger JSON file was not found at the specified path")
	SwaggerFileNotExist     = errors.New("the swagger file does not exist")
)

// Index action for list all swagger-JSON file list
func (s *indexController) Index(ctx *gin.Context) {
	data, err := service.ParseService.ParseWithCache()
	if err != nil {
		render.HtmlFail(ctx, SwaggerPathNotFoundFile)
		return
	}
	ctx.HTML(http.StatusOK, "list.html", data.Items)
}

// Detail action for swagger-JSON file detail
func (s *indexController) Detail(ctx *gin.Context) {
	hash := strings.TrimRight(ctx.Param("path"), ".html")
	data, err := service.ParseService.ParseWithCache()
	if err != nil {
		render.HtmlFail(ctx, SwaggerPathNotFoundFile)
		return
	}
	if _, exist := data.Table[hash]; !exist {
		render.HtmlFail(ctx, SwaggerFileNotExist)
		return
	}
	ctx.HTML(http.StatusOK, "detail.html", "/json/"+hash+".json")
}

// Json serve JSON file
func (s *indexController) Json(ctx *gin.Context) {
	hash := strings.TrimRight(ctx.Param("path"), ".json")
	data, err := service.ParseService.ParseWithCache()
	if err != nil {
		render.F(ctx, SwaggerPathNotFoundFile)
		return
	}
	path, exist := data.Table[hash]
	if !exist {
		render.F(ctx, SwaggerFileNotExist)
		return
	}
	stream, err := os.ReadFile(path)
	if err != nil {
		render.F(ctx, err)
		return
	}

	ctx.DataFromReader(http.StatusOK, int64(len(stream)), "application/json", bytes.NewReader(stream), nil)
}

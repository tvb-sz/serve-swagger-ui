package route

import (
	"github.com/gin-gonic/gin"
	"github.com/tvb-sz/serve-swagger-ui/app/controller"
	"github.com/tvb-sz/serve-swagger-ui/stubs"
	"html/template"
	"net/http"
)

func routeSetting() {
	// serve static file route, do not need to auth
	// http://domain/static/dist/xxx.css ---> ./stubs/dist/xxx.css
	router.StaticFS("/static", http.FS(stubs.Static))

	// serve favicon.ico
	router.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "image/x-icon")
		ctx.String(200, string(stubs.Favicon))
	})

	// register index page, use embed html file property
	router.SetHTMLTemplate(template.Must(template.ParseFS(stubs.Template, "./*.html")))
	router.GET("/", controller.IndexController.Index)
	router.GET("/index", controller.IndexController.Index)
	router.GET("/index.html", controller.IndexController.Index)
	router.GET("/index.htm", controller.IndexController.Index)
}

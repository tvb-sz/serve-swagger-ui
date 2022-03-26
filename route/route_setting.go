package route

import (
	"github.com/gin-gonic/gin"
	"github.com/tvb-sz/serve-swagger-ui/app/controller"
	"github.com/tvb-sz/serve-swagger-ui/stubs"
	"html/template"
	"net/http"
)

func routeSetting() {
	// serve favicon.ico, public accessible
	router.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "image/x-icon")
		ctx.String(200, string(stubs.Favicon))
	})

	// register index page, use embed html file property
	router.Use(tryAuthenticate)
	{
		router.SetHTMLTemplate(template.Must(template.ParseFS(stubs.Template, "./*.html")))
		router.GET("/", controller.IndexController.Index)
		router.GET("/index", controller.IndexController.Index)
		router.GET("/index.html", controller.IndexController.Index)
		router.GET("/index.htm", controller.IndexController.Index)
	}

	// authenticate, nothing when not need login, or should log in not authenticate redirect to index /
	router.Use(authenticate)
	{
		// serve static file route, do not need to auth
		// http://domain/static/dist/xxx.css ---> ./stubs/dist/xxx.css
		router.StaticFS("/static", http.FS(stubs.Static))

		// serve json file
		router.GET("/json/:path", controller.IndexController.Json)

		// register detail page, use embed html file property
		router.GET("/doc/:path", controller.IndexController.Detail)
	}

	// redirect when authenticated or not need login to index /
	router.Use(redirectIfAuthenticated)
	{
		router.GET("/oauth/google", controller.AuthController.LoginUsingGoogle)
		router.GET("/callback/google", controller.AuthController.CallbackUsingGoogle)
	}
}

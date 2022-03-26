package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tvb-sz/serve-swagger-ui/app/service"
	"github.com/tvb-sz/serve-swagger-ui/render"
	"net/http"
)

// authController auth controller
type authController struct{}

// LoginUsingGoogle to google oAuth login
func (s *authController) LoginUsingGoogle(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, service.OauthService.GoogleRedirectURL(ctx))
}

// CallbackUsingGoogle google oAuth login callback
func (s *authController) CallbackUsingGoogle(ctx *gin.Context) {
	if err := service.OauthService.GoogleCallback(ctx); err != nil {
		render.HtmlFail(ctx, err)
	}
	ctx.Redirect(http.StatusFound, "/")
}

// Logout exit logout
func (s *authController) Logout(ctx *gin.Context) {
	service.OauthService.DeleteCookie(ctx)
	ctx.Redirect(http.StatusFound, "/")
}

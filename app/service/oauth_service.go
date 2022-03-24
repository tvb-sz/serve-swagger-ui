package service

import (
	"github.com/gin-gonic/gin"
)

// JwtToken jwt token struct
type JwtToken struct {
	Sub string `json:"sub"` // JWT subject, just use oauth OpenID, such as google email address
	Exp int64  `json:"exp"` // JWT expired unix timestamp
	Iat int64  `json:"iat"` // JWT issue unix timestamp
}

// oauthService oAuth service
type oauthService struct{}

func (o *oauthService) CheckAuthorization(ctx *gin.Context) bool {
	//cookie, err := ctx.Cookie(define.AuthCookieName)
	//if err != nil {
	//	return false
	//}
	return false
}

func (o *oauthService) authenticate() {}

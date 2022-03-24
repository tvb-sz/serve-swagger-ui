package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tvb-sz/serve-swagger-ui/conf"
	"github.com/tvb-sz/serve-swagger-ui/define"
	"gopkg.in/square/go-jose.v2"
	"strings"
	"time"
)

var signer jose.Signer
var AuthorizationHasExpiredOrInvalid = errors.New("authorization has expired or invalid")

// JwtToken jwt token struct
type JwtToken struct {
	Sub string `json:"sub"` // JWT subject, just use oauth email, such as google email address
	Exp int64  `json:"exp"` // JWT expired unix timestamp
	Iat int64  `json:"iat"` // JWT issue unix timestamp
}

// oauthService oAuth service
type oauthService struct{}

func (o *oauthService) GoogleRedirectURL() string {
	return ""
}

// CheckAuthorization check cookie state
func (o *oauthService) CheckAuthorization(ctx *gin.Context) (token JwtToken, valid bool) {
	cookie, err := ctx.Cookie(define.AuthCookieName)
	if err != nil {
		return token, false
	}

	token, err = o.verifyJwt(cookie)
	if err != nil {
		return token, false
	}

	// check account can attempt
	if !o.canAttempt(token.Sub) {
		return JwtToken{}, false
	}

	return token, true
}

// canAttempt check email can attempt auth
func (o *oauthService) canAttempt(email string) bool {
	// ① check specify email
	for _, item := range conf.Config.Account.Email {
		if item == email {
			return true
		}
	}

	// ② check email suffix domain
	emailItems := strings.Split(email, "@")
	if len(emailItems) == 2 {
		for _, item := range conf.Config.Account.Domain {
			if item == emailItems[1] {
				return true
			}
		}
	}
	return false
}

// generateJwt use auth email generate jwt
// email ex: google email address
func (o *oauthService) generateJwt(email string) (jwt string, err error) {
	// init signer
	if signer == nil {
		signer, err = jose.NewSigner(jose.SigningKey{
			Algorithm: "HS256",
			Key:       conf.Config.Server.JwtKey,
		}, nil)
		if err != nil {
			return
		}
	}

	// construct token struct
	now := time.Now()
	token := JwtToken{
		Sub: email,
		Exp: now.Add(time.Duration(conf.Config.Server.JwtExpiredTime) * time.Second).Unix(),
		Iat: now.Unix(),
	}

	obj, err := json.Marshal(token)
	if err != nil {
		return
	}

	jwsObj, err := signer.Sign(obj)
	if err != nil {
		return
	}

	return jwsObj.CompactSerialize()
}

// verifyJwt verify JWT string
func (o *oauthService) verifyJwt(jwt string) (token JwtToken, err error) {
	var sig *jose.JSONWebSignature
	if sig, err = jose.ParseSigned(jwt); err != nil {
		return JwtToken{}, err
	}

	if err = json.Unmarshal(sig.UnsafePayloadWithoutVerification(), &token); err != nil {
		return JwtToken{}, err
	}

	// validate expired time
	if token.Exp <= time.Now().Unix() {
		return JwtToken{}, AuthorizationHasExpiredOrInvalid
	}

	// check key
	if _, err = sig.Verify([]byte(conf.Config.Server.JwtKey)); err != nil {
		return JwtToken{}, AuthorizationHasExpiredOrInvalid
	}

	return token, nil
}

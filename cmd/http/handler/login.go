package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
)

// Login is a handler for login.
type Login struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
	Config     config.Config
}

// LoginRequest is a request for Login.
type LoginRequest struct {
	LoginID  string `json:"login_id" validate:"required,max=255" ja:"ログインID" en:"LoginID"`
	Password string `json:"password" validate:"required,max=255" ja:"パスワード" en:"Password"`
}

func (h *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	var loginReq LoginRequest
	if err = json.NewDecoder(r.Body).Decode(&loginReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			loginReq = LoginRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &loginReq, lang.GetLocale(r.Context()))
	} else {
		err = errhandle.NewJSONFormatError()
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	var jwt entity.AuthJwt
	if jwt, err = h.Service.Login(
		ctx,
		loginReq.LoginID,
		loginReq.Password,
	); err == nil {
		cookie := new(http.Cookie)
		cookie.Name = auth.AccessTokenCookieKey
		cookie.Value = jwt.AccessToken
		cookie.Expires = time.Now().Add(h.Config.AuthRefreshTokenExpiresIn)
		cookie.SameSite = http.SameSiteLaxMode
		cookie.HttpOnly = true
		cookie.Secure = !h.Config.AppDebug
		http.SetCookie(w, cookie)
		err = response.JSONResponseWriter(ctx, w, response.Success, jwt, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

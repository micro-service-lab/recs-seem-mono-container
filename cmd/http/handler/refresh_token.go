package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// RefreshToken is a handler for refresh token.
type RefreshToken struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// RefreshTokenRequest is a request for RefreshToken.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" ja:"リフレッシュトークン" en:"RefreshToken"`
}

func (h *RefreshToken) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	var refreshTokenReq RefreshTokenRequest
	if err = json.NewDecoder(r.Body).Decode(&refreshTokenReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			refreshTokenReq = RefreshTokenRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &refreshTokenReq, lang.GetLocale(r.Context()))
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
	if jwt, err = h.Service.RefreshToken(
		ctx,
		refreshTokenReq.RefreshToken,
	); err == nil {
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

package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// UpdateMember is a handler for creating member.
type UpdateMember struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// UpdateMemberRequest is a request for UpdateMember.
type UpdateMemberRequest struct {
	Email          string    `json:"email" validate:"required,email,max=255" ja:"メールアドレス" en:"Email"`
	Name           string    `json:"name" validate:"required,max=255" ja:"名前" en:"Name"`
	FirstName      string    `json:"first_name" validate:"max=255" ja:"名" en:"FirstName"`
	LastName       string    `json:"last_name" validate:"max=255" ja:"姓" en:"LastName"`
	ProfileImageID uuid.UUID `json:"profile_image_id" validate:"" ja:"プロフィール画像ID" en:"ProfileImageID"`
}

func (h *UpdateMember) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "member_id"))
	var err error
	var memberReq UpdateMemberRequest
	if err = json.NewDecoder(r.Body).Decode(&memberReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			memberReq = UpdateMemberRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &memberReq, lang.GetLocale(r.Context()))
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
	var member entity.Member
	if member, err = h.Service.UpdateMember(
		ctx,
		id,
		memberReq.Email,
		memberReq.Name,
		entity.String{
			Valid:  memberReq.FirstName != "",
			String: memberReq.FirstName,
		},
		entity.String{
			Valid:  memberReq.LastName != "",
			String: memberReq.LastName,
		},
		entity.UUID{
			Valid: memberReq.ProfileImageID != uuid.Nil,
			Bytes: memberReq.ProfileImageID,
		},
	); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.MemberTargetProfileImages:
				profileImageStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ProfileImageID", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ProfileImageID",
						Other: "ProfileImageID",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ModelNotExists", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ModelNotExists",
						Other: "{{.ID}} not found",
					},
					TemplateData: map[string]any{
						"ID": profileImageStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("profile_image_id", msgStr)
				err = ve
			}
		}
	} else {
		err = response.JSONResponseWriter(ctx, w, response.Success, member, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

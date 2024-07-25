package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// UpdateAuth is a handler for creating member.
type UpdateAuth struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

func (h *UpdateAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
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
		authUser.MemberID,
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

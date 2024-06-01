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

// DeleteImage is a handler for creating organization.
type DeleteImage struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// DeleteImageRequest is a request for DeleteImage.
type DeleteImageRequest struct {
	ImageIDS []uuid.UUID `json:"image_ids" validate:"required,unique,min=1" ja:"画像ID" en:"ImageIDs"`
}

func (h *DeleteImage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	var err error
	var req DeleteImageRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			req = DeleteImageRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &req, lang.GetLocale(r.Context()))
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
	if _, err = h.Service.PluralDeleteImages(
		ctx,
		req.ImageIDS,
		entity.UUID{Bytes: authUser.MemberID, Valid: true},
	); err != nil {
		var ce errhandle.CommonError
		if errors.As(err, &ce) {
			if ce.Code.Code == response.CannotDeleteSystemFile.Code {
				imageStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ImageIDs", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ImageIDs",
						Other: "ImageIDs",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "ContainSystemFile", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "ContainSystemFile",
							Other: "{{.Target}} contains system files",
						},
						TemplateData: map[string]any{
							"Target": imageStr,
						},
					})
				ve := errhandle.NewValidationError(nil)
				ve.Add("image_ids", msgStr)
				err = ve
			} else if ce.Code.Code == response.NotFileOwner.Code {
				imageStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ImageIDs", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ImageIDs",
						Other: "ImageIDs",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "ContainNotOwnerFile", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "ContainNotOwnerFile",
							Other: "{{.Target}} contains files that are not owned by the user",
						},
						TemplateData: map[string]any{
							"Target": imageStr,
						},
					})
				ve := errhandle.NewValidationError(nil)
				ve.Add("image_ids", msgStr)
				err = ve
			}
		}
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.ImageTargetImages:
				imageStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ImageIDs", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ImageIDs",
						Other: "ImageIDs",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "PluralModelNotExists", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "PluralModelNotExists",
							Other: "{{.ID}} not found",
						},
						TemplateData: map[string]any{
							"ID":        imageStr,
							"ValueType": "ID",
						},
					})
				ve := errhandle.NewValidationError(nil)
				ve.Add("image_ids", msgStr)
				err = ve
			}
		}
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, nil, nil)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

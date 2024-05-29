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

// UpdateStudentGroup is a handler for updating student group.
type UpdateStudentGroup struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// UpdateStudentGroupRequest is a request for UpdateStudentGroup.
type UpdateStudentGroupRequest struct {
	GroupID uuid.UUID `json:"group_id" validate:"required" ja:"班ID" en:"GroupID"`
}

func (h *UpdateStudentGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "student_id"))
	var err error
	var studentReq UpdateStudentGroupRequest
	if err = json.NewDecoder(r.Body).Decode(&studentReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			studentReq = UpdateStudentGroupRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &studentReq, lang.GetLocale(r.Context()))
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
	var student entity.Student
	if student, err = h.Service.UpdateStudentGroup(
		ctx,
		id,
		studentReq.GroupID,
	); err != nil {
		var ce errhandle.CommonError
		if errors.As(err, &ce) && ce.Code.Code == response.OnlyProfessorAction.Code {
			switch ce.Target {
			case service.MemberTargetGroups:
				groupStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "GroupID", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "GroupID",
						Other: "GroupID",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "OnlyProfessorModel", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "OnlyProfessorModel",
						Other: "{{.ID}} professor only",
					},
					TemplateData: map[string]any{
						"ID": groupStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("group_id", msgStr)
				err = ve
			}
		}
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.MemberTargetGroups:
				groupStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "GroupID", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "GroupID",
						Other: "GroupID",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ModelNotExists", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ModelNotExists",
						Other: "{{.ID}} not found",
					},
					TemplateData: map[string]any{
						"ID": groupStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("group_id", msgStr)
				err = ve
			}
		}
	} else {
		err = response.JSONResponseWriter(ctx, w, response.Success, student, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

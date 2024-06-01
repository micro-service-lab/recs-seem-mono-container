package handler

import (
	"log"
	"net/http"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// CreateFile is a handler for creating file.
type CreateFile struct {
	Service    service.ManagerInterface
	Translator i18n.Translation
}

func (h *CreateFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const (
		propName  = "files"
		maxMemory = 32 << 20
	)
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	var err error

	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		err = response.JSONResponseWriter(ctx, w, response.MultiPartFormParseError, nil, nil)
		if err != nil {
			log.Printf("failed to write response: %+v", err)
		}
		return
	}
	formdata := r.MultipartForm
	pName := propName + "[]"
	v, ok := formdata.Value[pName]

	if ok && len(v) > 0 {
		fileStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "Files", i18n.Options{
			DefaultMessage: &i18n.Message{
				ID:    "Files",
				Other: "Files",
			},
		})
		msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "PluralIsUploadFile", i18n.Options{
			DefaultMessage: &i18n.Message{
				ID:    "PluralIsUploadFile",
				Other: "{{.Target}} must be an upload file",
			},
			TemplateData: map[string]any{
				"Target": fileStr,
			},
		})
		ve := errhandle.NewValidationError(nil)
		ve.Add(pName, msgStr)
		err = ve
	}
	filHeaders, ok := formdata.File[pName]
	if !ok || len(filHeaders) == 0 {
		fileStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "Files", i18n.Options{
			DefaultMessage: &i18n.Message{
				ID:    "Files",
				Other: "Files",
			},
		})
		msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "Required", i18n.Options{
			DefaultMessage: &i18n.Message{
				ID:    "Required",
				Other: "{{.Target}} is a required field",
			},
			TemplateData: map[string]any{
				"Target": fileStr,
			},
		})
		ve := errhandle.NewValidationError(nil)
		ve.Add(pName, msgStr)
		err = ve
	}

	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}

	p := make([]parameter.CreateFileServiceParam, 0, len(filHeaders))
	for _, header := range filHeaders {
		file, err := header.Open()
		if err != nil {
			return
		}
		defer file.Close()

		if file == nil {
			continue
		}

		p = append(p, parameter.CreateFileServiceParam{
			Origin: file,
			Alias:  header.Filename,
		})
	}

	var files []entity.FileWithAttachableItem
	if files, err = h.Service.CreateFiles(
		ctx,
		entity.UUID{
			Valid: true,
			Bytes: authUser.MemberID,
		},
		p,
	); err == nil {
		err = response.JSONResponseWriter(ctx, w, response.Success, files, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// StudentHandler is a handler for roles.
func StudentHandler(svc service.ManagerInterface, vd validation.Validator, t i18n.Translation) http.Handler {
	createHandler := handler.CreateStudent{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	deleteHandler := handler.DeleteStudent{
		Service: svc,
	}
	updateGroupHandler := handler.UpdateStudentGroup{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	updateGradeHandler := handler.UpdateStudentGrade{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}

	r := chi.NewRouter()
	r.Post("/", createHandler.ServeHTTP)
	r.Delete("/{student_id}", deleteHandler.ServeHTTP)
	r.Put("/{student_id}/group", updateGroupHandler.ServeHTTP)
	r.Put("/{student_id}/grade", updateGradeHandler.ServeHTTP)

	return r
}

package handler

import (
	"log"
	"net/http"
	"reflect"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// GetGrades is a handler for getting grades.
type GetGrades struct {
	Service service.ManagerInterface
}

// GetGradesParam is a parameter for GetGrades.
type GetGradesParam struct {
	SearchName string                            `queryParam:"search_name"`
	Order      parameter.OrganizationOrderMethod `queryParam:"order"`
	Limit      parameter.Limit                   `queryParam:"limit"`
	Offset     parameter.Offset                  `queryParam:"offset"`
	Cursor     parameter.Cursor                  `queryParam:"cursor"`
	Pagination parameter.Pagination              `queryParam:"pagination"`
	WithCount  parameter.WithCount               `queryParam:"with_count"`
	With       parameter.OrganizationWithParams  `queryParam:"with[]"`
}

var getGradesParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.OrganizationOrderMethodName): parameter.ParseOrganizationOrderMethod,
	reflect.TypeOf(parameter.Limit(0)):                    parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):                   parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):                  parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):              parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):            parameter.ParseWithCountParam,
	reflect.TypeOf(parameter.OrganizationWith{}):          parameter.ParseOrganizationWithParam,
}

func (h *GetGrades) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetGradesParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getGradesParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	var organizations any
	switch param.With.Case() {
	case parameter.OrganizationWithCaseDefault:
		organizations, err = h.Service.GetOrganizations(
			ctx,
			param.SearchName,
			parameter.WhereOrganizationTypeGrade,
			uuid.UUID{},
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	case parameter.OrganizationWithCaseChatRoom:
		organizations, err = h.Service.GetOrganizationsWithChatRoom(
			ctx,
			param.SearchName,
			parameter.WhereOrganizationTypeGrade,
			uuid.UUID{},
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	case parameter.OrganizationWithCaseDetail:
		organizations, err = h.Service.GetOrganizationsWithDetail(
			ctx,
			param.SearchName,
			parameter.WhereOrganizationTypeGrade,
			uuid.UUID{},
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	case parameter.OrganizationWithCaseChatRoomAndDetail:
		organizations, err = h.Service.GetOrganizationsWithChatRoomAndDetail(
			ctx,
			param.SearchName,
			parameter.WhereOrganizationTypeGrade,
			uuid.UUID{},
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	}

	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, organizations, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

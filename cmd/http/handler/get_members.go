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

// GetMembers is a handler for getting organizations.
type GetMembers struct {
	Service service.ManagerInterface
}

// GetMembersParam is a parameter for GetMembers.
type GetMembersParam struct {
	SearchName                       string                      `queryParam:"search_name"`
	SearchHasPolicies                []parameter.EntityID        `queryParam:"search_has_policies[]"`
	SearchAttendStatuses             []parameter.EntityID        `queryParam:"search_attend_statuses[]"`
	SearchGrades                     []parameter.EntityID        `queryParam:"search_grades[]"`
	SearchGroups                     []parameter.EntityID        `queryParam:"search_groups[]"`
	SearchBelongingOrganizationID    parameter.EntityID          `queryParam:"search_belonging_organization_id"`
	SearchNotBelongingOrganizationID parameter.EntityID          `queryParam:"search_not_belonging_organization_id"`
	SearchBelongingChatRoomID        parameter.EntityID          `queryParam:"search_belonging_chat_room_id"`
	SearchNotBelongingChatRoomID     parameter.EntityID          `queryParam:"search_not_belonging_chat_room_id"`
	Order                            parameter.MemberOrderMethod `queryParam:"order"`
	Limit                            parameter.Limit             `queryParam:"limit"`
	Offset                           parameter.Offset            `queryParam:"offset"`
	Cursor                           parameter.Cursor            `queryParam:"cursor"`
	Pagination                       parameter.Pagination        `queryParam:"pagination"`
	WithCount                        parameter.WithCount         `queryParam:"with_count"`
	With                             parameter.MemberWithParams  `queryParam:"with[]"`
}

var getMembersParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.MemberOrderMethodName): parameter.ParseMemberOrderMethod,
	reflect.TypeOf(parameter.EntityID(uuid.UUID{})): parameter.ParseEntityIDParam,
	reflect.TypeOf(parameter.Limit(0)):              parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):             parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):            parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):        parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):      parameter.ParseWithCountParam,
	reflect.TypeOf(parameter.MemberWith{}):          parameter.ParseMemberWithParam,
}

func (h *GetMembers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetMembersParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getMembersParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	inHasPolicies := make([]uuid.UUID, 0, len(param.SearchHasPolicies))
	for _, v := range param.SearchHasPolicies {
		if uuid.UUID(v) != uuid.Nil {
			inHasPolicies = append(inHasPolicies, uuid.UUID(v))
		}
	}
	inAttendStatuses := make([]uuid.UUID, 0, len(param.SearchAttendStatuses))
	for _, v := range param.SearchAttendStatuses {
		if uuid.UUID(v) != uuid.Nil {
			inAttendStatuses = append(inAttendStatuses, uuid.UUID(v))
		}
	}
	inGrades := make([]uuid.UUID, 0, len(param.SearchGrades))
	for _, v := range param.SearchGrades {
		if uuid.UUID(v) != uuid.Nil {
			inGrades = append(inGrades, uuid.UUID(v))
		}
	}
	inGroups := make([]uuid.UUID, 0, len(param.SearchGroups))
	for _, v := range param.SearchGroups {
		if uuid.UUID(v) != uuid.Nil {
			inGroups = append(inGroups, uuid.UUID(v))
		}
	}
	var members any
	switch param.With.Case() {
	case parameter.MemberWithCaseDefault:
		members, err = h.Service.GetMembers(
			ctx,
			param.SearchName,
			inHasPolicies,
			inAttendStatuses,
			inGrades,
			inGroups,
			uuid.UUID(param.SearchBelongingOrganizationID),
			uuid.UUID(param.SearchNotBelongingOrganizationID),
			uuid.UUID(param.SearchBelongingChatRoomID),
			uuid.UUID(param.SearchNotBelongingChatRoomID),
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	case parameter.MemberWithCaseDetail:
		members, err = h.Service.GetMembersWithDetail(
			ctx,
			param.SearchName,
			inHasPolicies,
			inAttendStatuses,
			inGrades,
			inGroups,
			uuid.UUID(param.SearchBelongingOrganizationID),
			uuid.UUID(param.SearchNotBelongingOrganizationID),
			uuid.UUID(param.SearchBelongingChatRoomID),
			uuid.UUID(param.SearchNotBelongingChatRoomID),
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	case parameter.MemberWithCaseProfileImage:
		members, err = h.Service.GetMembersWithProfileImage(
			ctx,
			param.SearchName,
			inHasPolicies,
			inAttendStatuses,
			inGrades,
			inGroups,
			uuid.UUID(param.SearchBelongingOrganizationID),
			uuid.UUID(param.SearchNotBelongingOrganizationID),
			uuid.UUID(param.SearchBelongingChatRoomID),
			uuid.UUID(param.SearchNotBelongingChatRoomID),
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	case parameter.MemberWithCasePersonalOrganization:
		members, err = h.Service.GetMembersWithPersonalOrganization(
			ctx,
			param.SearchName,
			inHasPolicies,
			inAttendStatuses,
			inGrades,
			inGroups,
			uuid.UUID(param.SearchBelongingOrganizationID),
			uuid.UUID(param.SearchNotBelongingOrganizationID),
			uuid.UUID(param.SearchBelongingChatRoomID),
			uuid.UUID(param.SearchNotBelongingChatRoomID),
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	case parameter.MemberWithCaseCrew:
		members, err = h.Service.GetMembersWithCrew(
			ctx,
			param.SearchName,
			inHasPolicies,
			inAttendStatuses,
			inGrades,
			inGroups,
			uuid.UUID(param.SearchBelongingOrganizationID),
			uuid.UUID(param.SearchNotBelongingOrganizationID),
			uuid.UUID(param.SearchBelongingChatRoomID),
			uuid.UUID(param.SearchNotBelongingChatRoomID),
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	case parameter.MemberWithCaseAttendStatus:
		members, err = h.Service.GetMembersWithAttendStatus(
			ctx,
			param.SearchName,
			inHasPolicies,
			inAttendStatuses,
			inGrades,
			inGroups,
			uuid.UUID(param.SearchBelongingOrganizationID),
			uuid.UUID(param.SearchNotBelongingOrganizationID),
			uuid.UUID(param.SearchBelongingChatRoomID),
			uuid.UUID(param.SearchNotBelongingChatRoomID),
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	case parameter.MemberWithCaseRole:
		members, err = h.Service.GetMembersWithRole(
			ctx,
			param.SearchName,
			inHasPolicies,
			inAttendStatuses,
			inGrades,
			inGroups,
			uuid.UUID(param.SearchBelongingOrganizationID),
			uuid.UUID(param.SearchNotBelongingOrganizationID),
			uuid.UUID(param.SearchBelongingChatRoomID),
			uuid.UUID(param.SearchNotBelongingChatRoomID),
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	case parameter.MemberWithCaseCrew | parameter.MemberWithCaseAttendStatus | parameter.MemberWithCaseProfileImage:
		members, err = h.Service.GetMembersWithCrewAndProfileImageAndAttendStatus(
			ctx,
			param.SearchName,
			inHasPolicies,
			inAttendStatuses,
			inGrades,
			inGroups,
			uuid.UUID(param.SearchBelongingOrganizationID),
			uuid.UUID(param.SearchNotBelongingOrganizationID),
			uuid.UUID(param.SearchBelongingChatRoomID),
			uuid.UUID(param.SearchNotBelongingChatRoomID),
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	default:
		members, err = h.Service.GetMembers(
			ctx,
			param.SearchName,
			inHasPolicies,
			inAttendStatuses,
			inGrades,
			inGroups,
			uuid.UUID(param.SearchBelongingOrganizationID),
			uuid.UUID(param.SearchNotBelongingOrganizationID),
			uuid.UUID(param.SearchBelongingChatRoomID),
			uuid.UUID(param.SearchNotBelongingChatRoomID),
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
	err = response.JSONResponseWriter(ctx, w, response.Success, members, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

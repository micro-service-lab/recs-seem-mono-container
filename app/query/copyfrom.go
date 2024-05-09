// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: copyfrom.go

package query

import (
	"context"

	"github.com/google/uuid"
)

// iteratorForCreateAbsences implements pgx.CopyFromSource.
type iteratorForCreateAbsences struct {
	rows                 []uuid.UUID
	skippedFirstNextCall bool
}

func (r *iteratorForCreateAbsences) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateAbsences) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0],
	}, nil
}

func (r iteratorForCreateAbsences) Err() error {
	return nil
}

func (q *Queries) CreateAbsences(ctx context.Context, attendanceID []uuid.UUID) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_absences"}, []string{"attendance_id"}, &iteratorForCreateAbsences{rows: attendanceID})
}

// iteratorForCreateAttachableItems implements pgx.CopyFromSource.
type iteratorForCreateAttachableItems struct {
	rows                 []CreateAttachableItemsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateAttachableItems) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateAttachableItems) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Url,
		r.rows[0].Size,
		r.rows[0].OwnerID,
		r.rows[0].MimeTypeID,
	}, nil
}

func (r iteratorForCreateAttachableItems) Err() error {
	return nil
}

func (q *Queries) CreateAttachableItems(ctx context.Context, arg []CreateAttachableItemsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_attachable_items"}, []string{"url", "size", "owner_id", "mime_type_id"}, &iteratorForCreateAttachableItems{rows: arg})
}

// iteratorForCreateAttachedMessages implements pgx.CopyFromSource.
type iteratorForCreateAttachedMessages struct {
	rows                 []CreateAttachedMessagesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateAttachedMessages) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateAttachedMessages) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].MessageID,
		r.rows[0].FileUrl,
	}, nil
}

func (r iteratorForCreateAttachedMessages) Err() error {
	return nil
}

func (q *Queries) CreateAttachedMessages(ctx context.Context, arg []CreateAttachedMessagesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_attached_messages"}, []string{"message_id", "file_url"}, &iteratorForCreateAttachedMessages{rows: arg})
}

// iteratorForCreateAttendStatuses implements pgx.CopyFromSource.
type iteratorForCreateAttendStatuses struct {
	rows                 []CreateAttendStatusesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateAttendStatuses) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateAttendStatuses) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Key,
	}, nil
}

func (r iteratorForCreateAttendStatuses) Err() error {
	return nil
}

func (q *Queries) CreateAttendStatuses(ctx context.Context, arg []CreateAttendStatusesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_attend_statuses"}, []string{"name", "key"}, &iteratorForCreateAttendStatuses{rows: arg})
}

// iteratorForCreateAttendanceTypes implements pgx.CopyFromSource.
type iteratorForCreateAttendanceTypes struct {
	rows                 []CreateAttendanceTypesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateAttendanceTypes) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateAttendanceTypes) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Key,
		r.rows[0].Color,
	}, nil
}

func (r iteratorForCreateAttendanceTypes) Err() error {
	return nil
}

func (q *Queries) CreateAttendanceTypes(ctx context.Context, arg []CreateAttendanceTypesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_attendance_types"}, []string{"name", "key", "color"}, &iteratorForCreateAttendanceTypes{rows: arg})
}

// iteratorForCreateAttendances implements pgx.CopyFromSource.
type iteratorForCreateAttendances struct {
	rows                 []CreateAttendancesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateAttendances) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateAttendances) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].AttendanceTypeID,
		r.rows[0].MemberID,
		r.rows[0].Description,
		r.rows[0].Date,
		r.rows[0].MailSendFlag,
		r.rows[0].SendOrganizationID,
		r.rows[0].PostedAt,
		r.rows[0].LastEditedAt,
	}, nil
}

func (r iteratorForCreateAttendances) Err() error {
	return nil
}

func (q *Queries) CreateAttendances(ctx context.Context, arg []CreateAttendancesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_attendances"}, []string{"attendance_type_id", "member_id", "description", "date", "mail_send_flag", "send_organization_id", "posted_at", "last_edited_at"}, &iteratorForCreateAttendances{rows: arg})
}

// iteratorForCreateChatRoomBelongings implements pgx.CopyFromSource.
type iteratorForCreateChatRoomBelongings struct {
	rows                 []CreateChatRoomBelongingsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateChatRoomBelongings) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateChatRoomBelongings) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].MemberID,
		r.rows[0].ChatRoomID,
		r.rows[0].AddedAt,
	}, nil
}

func (r iteratorForCreateChatRoomBelongings) Err() error {
	return nil
}

func (q *Queries) CreateChatRoomBelongings(ctx context.Context, arg []CreateChatRoomBelongingsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_chat_room_belongings"}, []string{"member_id", "chat_room_id", "added_at"}, &iteratorForCreateChatRoomBelongings{rows: arg})
}

// iteratorForCreateChatRooms implements pgx.CopyFromSource.
type iteratorForCreateChatRooms struct {
	rows                 []CreateChatRoomsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateChatRooms) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateChatRooms) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].IsPrivate,
		r.rows[0].CoverImageUrl,
		r.rows[0].OwnerID,
		r.rows[0].FromOrganization,
		r.rows[0].CreatedAt,
		r.rows[0].UpdatedAt,
	}, nil
}

func (r iteratorForCreateChatRooms) Err() error {
	return nil
}

func (q *Queries) CreateChatRooms(ctx context.Context, arg []CreateChatRoomsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_chat_rooms"}, []string{"name", "is_private", "cover_image_url", "owner_id", "from_organization", "created_at", "updated_at"}, &iteratorForCreateChatRooms{rows: arg})
}

// iteratorForCreateEarlyLeavings implements pgx.CopyFromSource.
type iteratorForCreateEarlyLeavings struct {
	rows                 []CreateEarlyLeavingsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateEarlyLeavings) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateEarlyLeavings) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].AttendanceID,
		r.rows[0].LeaveTime,
	}, nil
}

func (r iteratorForCreateEarlyLeavings) Err() error {
	return nil
}

func (q *Queries) CreateEarlyLeavings(ctx context.Context, arg []CreateEarlyLeavingsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_early_leavings"}, []string{"attendance_id", "leave_time"}, &iteratorForCreateEarlyLeavings{rows: arg})
}

// iteratorForCreateEventTypes implements pgx.CopyFromSource.
type iteratorForCreateEventTypes struct {
	rows                 []CreateEventTypesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateEventTypes) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateEventTypes) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Key,
		r.rows[0].Color,
	}, nil
}

func (r iteratorForCreateEventTypes) Err() error {
	return nil
}

func (q *Queries) CreateEventTypes(ctx context.Context, arg []CreateEventTypesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_event_types"}, []string{"name", "key", "color"}, &iteratorForCreateEventTypes{rows: arg})
}

// iteratorForCreateEvents implements pgx.CopyFromSource.
type iteratorForCreateEvents struct {
	rows                 []CreateEventsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateEvents) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateEvents) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].EventTypeID,
		r.rows[0].Title,
		r.rows[0].Description,
		r.rows[0].OrganizationID,
		r.rows[0].StartTime,
		r.rows[0].EndTime,
		r.rows[0].MailSendFlag,
		r.rows[0].SendOrganizationID,
		r.rows[0].PostedBy,
		r.rows[0].LastEditedBy,
		r.rows[0].PostedAt,
		r.rows[0].LastEditedAt,
	}, nil
}

func (r iteratorForCreateEvents) Err() error {
	return nil
}

func (q *Queries) CreateEvents(ctx context.Context, arg []CreateEventsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_events"}, []string{"event_type_id", "title", "description", "organization_id", "start_time", "end_time", "mail_send_flag", "send_organization_id", "posted_by", "last_edited_by", "posted_at", "last_edited_at"}, &iteratorForCreateEvents{rows: arg})
}

// iteratorForCreateFiles implements pgx.CopyFromSource.
type iteratorForCreateFiles struct {
	rows                 []uuid.UUID
	skippedFirstNextCall bool
}

func (r *iteratorForCreateFiles) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateFiles) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0],
	}, nil
}

func (r iteratorForCreateFiles) Err() error {
	return nil
}

func (q *Queries) CreateFiles(ctx context.Context, attachableItemID []uuid.UUID) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_files"}, []string{"attachable_item_id"}, &iteratorForCreateFiles{rows: attachableItemID})
}

// iteratorForCreateGrades implements pgx.CopyFromSource.
type iteratorForCreateGrades struct {
	rows                 []CreateGradesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateGrades) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateGrades) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Key,
		r.rows[0].OrganizationID,
	}, nil
}

func (r iteratorForCreateGrades) Err() error {
	return nil
}

func (q *Queries) CreateGrades(ctx context.Context, arg []CreateGradesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_grades"}, []string{"key", "organization_id"}, &iteratorForCreateGrades{rows: arg})
}

// iteratorForCreateGroups implements pgx.CopyFromSource.
type iteratorForCreateGroups struct {
	rows                 []CreateGroupsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateGroups) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateGroups) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Key,
		r.rows[0].OrganizationID,
	}, nil
}

func (r iteratorForCreateGroups) Err() error {
	return nil
}

func (q *Queries) CreateGroups(ctx context.Context, arg []CreateGroupsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_groups"}, []string{"key", "organization_id"}, &iteratorForCreateGroups{rows: arg})
}

// iteratorForCreateImages implements pgx.CopyFromSource.
type iteratorForCreateImages struct {
	rows                 []CreateImagesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateImages) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateImages) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Height,
		r.rows[0].Width,
		r.rows[0].AttachableItemID,
	}, nil
}

func (r iteratorForCreateImages) Err() error {
	return nil
}

func (q *Queries) CreateImages(ctx context.Context, arg []CreateImagesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_images"}, []string{"height", "width", "attachable_item_id"}, &iteratorForCreateImages{rows: arg})
}

// iteratorForCreateLabIOHistories implements pgx.CopyFromSource.
type iteratorForCreateLabIOHistories struct {
	rows                 []CreateLabIOHistoriesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateLabIOHistories) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateLabIOHistories) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].MemberID,
		r.rows[0].EnteredAt,
		r.rows[0].ExitedAt,
	}, nil
}

func (r iteratorForCreateLabIOHistories) Err() error {
	return nil
}

func (q *Queries) CreateLabIOHistories(ctx context.Context, arg []CreateLabIOHistoriesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_lab_io_histories"}, []string{"member_id", "entered_at", "exited_at"}, &iteratorForCreateLabIOHistories{rows: arg})
}

// iteratorForCreateLateArrivals implements pgx.CopyFromSource.
type iteratorForCreateLateArrivals struct {
	rows                 []CreateLateArrivalsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateLateArrivals) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateLateArrivals) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].AttendanceID,
		r.rows[0].ArriveTime,
	}, nil
}

func (r iteratorForCreateLateArrivals) Err() error {
	return nil
}

func (q *Queries) CreateLateArrivals(ctx context.Context, arg []CreateLateArrivalsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_late_arrivals"}, []string{"attendance_id", "arrive_time"}, &iteratorForCreateLateArrivals{rows: arg})
}

// iteratorForCreateMembers implements pgx.CopyFromSource.
type iteratorForCreateMembers struct {
	rows                 []CreateMembersParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateMembers) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateMembers) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].LoginID,
		r.rows[0].Password,
		r.rows[0].Email,
		r.rows[0].Name,
		r.rows[0].AttendStatusID,
		r.rows[0].GradeID,
		r.rows[0].GroupID,
		r.rows[0].ProfileImageUrl,
		r.rows[0].RoleID,
		r.rows[0].PersonalOrganizationID,
		r.rows[0].CreatedAt,
		r.rows[0].UpdatedAt,
	}, nil
}

func (r iteratorForCreateMembers) Err() error {
	return nil
}

func (q *Queries) CreateMembers(ctx context.Context, arg []CreateMembersParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_members"}, []string{"login_id", "password", "email", "name", "attend_status_id", "grade_id", "group_id", "profile_image_url", "role_id", "personal_organization_id", "created_at", "updated_at"}, &iteratorForCreateMembers{rows: arg})
}

// iteratorForCreateMessages implements pgx.CopyFromSource.
type iteratorForCreateMessages struct {
	rows                 []CreateMessagesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateMessages) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateMessages) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ChatRoomID,
		r.rows[0].SenderID,
		r.rows[0].Body,
		r.rows[0].PostedAt,
		r.rows[0].LastEditedAt,
	}, nil
}

func (r iteratorForCreateMessages) Err() error {
	return nil
}

func (q *Queries) CreateMessages(ctx context.Context, arg []CreateMessagesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_messages"}, []string{"chat_room_id", "sender_id", "body", "posted_at", "last_edited_at"}, &iteratorForCreateMessages{rows: arg})
}

// iteratorForCreateMimeTypes implements pgx.CopyFromSource.
type iteratorForCreateMimeTypes struct {
	rows                 []CreateMimeTypesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateMimeTypes) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateMimeTypes) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Key,
		r.rows[0].Kind,
	}, nil
}

func (r iteratorForCreateMimeTypes) Err() error {
	return nil
}

func (q *Queries) CreateMimeTypes(ctx context.Context, arg []CreateMimeTypesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_mime_types"}, []string{"name", "key", "kind"}, &iteratorForCreateMimeTypes{rows: arg})
}

// iteratorForCreateOrganizations implements pgx.CopyFromSource.
type iteratorForCreateOrganizations struct {
	rows                 []CreateOrganizationsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateOrganizations) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateOrganizations) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Description,
		r.rows[0].Color,
		r.rows[0].IsPersonal,
		r.rows[0].IsWhole,
		r.rows[0].ChatRoomID,
		r.rows[0].CreatedAt,
		r.rows[0].UpdatedAt,
	}, nil
}

func (r iteratorForCreateOrganizations) Err() error {
	return nil
}

func (q *Queries) CreateOrganizations(ctx context.Context, arg []CreateOrganizationsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_organizations"}, []string{"name", "description", "color", "is_personal", "is_whole", "chat_room_id", "created_at", "updated_at"}, &iteratorForCreateOrganizations{rows: arg})
}

// iteratorForCreatePermissionAssociations implements pgx.CopyFromSource.
type iteratorForCreatePermissionAssociations struct {
	rows                 []CreatePermissionAssociationsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreatePermissionAssociations) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreatePermissionAssociations) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].PermissionID,
		r.rows[0].WorkPositionID,
	}, nil
}

func (r iteratorForCreatePermissionAssociations) Err() error {
	return nil
}

func (q *Queries) CreatePermissionAssociations(ctx context.Context, arg []CreatePermissionAssociationsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_permission_associations"}, []string{"permission_id", "work_position_id"}, &iteratorForCreatePermissionAssociations{rows: arg})
}

// iteratorForCreatePermissionCategories implements pgx.CopyFromSource.
type iteratorForCreatePermissionCategories struct {
	rows                 []CreatePermissionCategoriesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreatePermissionCategories) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreatePermissionCategories) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Description,
		r.rows[0].Key,
	}, nil
}

func (r iteratorForCreatePermissionCategories) Err() error {
	return nil
}

func (q *Queries) CreatePermissionCategories(ctx context.Context, arg []CreatePermissionCategoriesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_permission_categories"}, []string{"name", "description", "key"}, &iteratorForCreatePermissionCategories{rows: arg})
}

// iteratorForCreatePermissions implements pgx.CopyFromSource.
type iteratorForCreatePermissions struct {
	rows                 []CreatePermissionsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreatePermissions) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreatePermissions) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Description,
		r.rows[0].Key,
		r.rows[0].PermissionCategoryID,
	}, nil
}

func (r iteratorForCreatePermissions) Err() error {
	return nil
}

func (q *Queries) CreatePermissions(ctx context.Context, arg []CreatePermissionsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_permissions"}, []string{"name", "description", "key", "permission_category_id"}, &iteratorForCreatePermissions{rows: arg})
}

// iteratorForCreatePolicies implements pgx.CopyFromSource.
type iteratorForCreatePolicies struct {
	rows                 []CreatePoliciesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreatePolicies) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreatePolicies) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Description,
		r.rows[0].Key,
		r.rows[0].PolicyCategoryID,
	}, nil
}

func (r iteratorForCreatePolicies) Err() error {
	return nil
}

func (q *Queries) CreatePolicies(ctx context.Context, arg []CreatePoliciesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_policies"}, []string{"name", "description", "key", "policy_category_id"}, &iteratorForCreatePolicies{rows: arg})
}

// iteratorForCreatePolicyCategories implements pgx.CopyFromSource.
type iteratorForCreatePolicyCategories struct {
	rows                 []CreatePolicyCategoriesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreatePolicyCategories) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreatePolicyCategories) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Description,
		r.rows[0].Key,
	}, nil
}

func (r iteratorForCreatePolicyCategories) Err() error {
	return nil
}

func (q *Queries) CreatePolicyCategories(ctx context.Context, arg []CreatePolicyCategoriesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_policy_categories"}, []string{"name", "description", "key"}, &iteratorForCreatePolicyCategories{rows: arg})
}

// iteratorForCreatePositionHistories implements pgx.CopyFromSource.
type iteratorForCreatePositionHistories struct {
	rows                 []CreatePositionHistoriesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreatePositionHistories) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreatePositionHistories) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].MemberID,
		r.rows[0].XPos,
		r.rows[0].YPos,
		r.rows[0].SentAt,
	}, nil
}

func (r iteratorForCreatePositionHistories) Err() error {
	return nil
}

func (q *Queries) CreatePositionHistories(ctx context.Context, arg []CreatePositionHistoriesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_position_histories"}, []string{"member_id", "x_pos", "y_pos", "sent_at"}, &iteratorForCreatePositionHistories{rows: arg})
}

// iteratorForCreateProfessors implements pgx.CopyFromSource.
type iteratorForCreateProfessors struct {
	rows                 []uuid.UUID
	skippedFirstNextCall bool
}

func (r *iteratorForCreateProfessors) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateProfessors) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0],
	}, nil
}

func (r iteratorForCreateProfessors) Err() error {
	return nil
}

func (q *Queries) CreateProfessors(ctx context.Context, memberID []uuid.UUID) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_professors"}, []string{"member_id"}, &iteratorForCreateProfessors{rows: memberID})
}

// iteratorForCreateRecordTypes implements pgx.CopyFromSource.
type iteratorForCreateRecordTypes struct {
	rows                 []CreateRecordTypesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateRecordTypes) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateRecordTypes) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Key,
	}, nil
}

func (r iteratorForCreateRecordTypes) Err() error {
	return nil
}

func (q *Queries) CreateRecordTypes(ctx context.Context, arg []CreateRecordTypesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_record_types"}, []string{"name", "key"}, &iteratorForCreateRecordTypes{rows: arg})
}

// iteratorForCreateRecords implements pgx.CopyFromSource.
type iteratorForCreateRecords struct {
	rows                 []CreateRecordsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateRecords) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateRecords) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RecordTypeID,
		r.rows[0].Title,
		r.rows[0].Body,
		r.rows[0].OrganizationID,
		r.rows[0].PostedBy,
		r.rows[0].LastEditedBy,
		r.rows[0].PostedAt,
		r.rows[0].LastEditedAt,
	}, nil
}

func (r iteratorForCreateRecords) Err() error {
	return nil
}

func (q *Queries) CreateRecords(ctx context.Context, arg []CreateRecordsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"t_records"}, []string{"record_type_id", "title", "body", "organization_id", "posted_by", "last_edited_by", "posted_at", "last_edited_at"}, &iteratorForCreateRecords{rows: arg})
}

// iteratorForCreateRoleAssociations implements pgx.CopyFromSource.
type iteratorForCreateRoleAssociations struct {
	rows                 []CreateRoleAssociationsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateRoleAssociations) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateRoleAssociations) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RoleID,
		r.rows[0].PolicyID,
	}, nil
}

func (r iteratorForCreateRoleAssociations) Err() error {
	return nil
}

func (q *Queries) CreateRoleAssociations(ctx context.Context, arg []CreateRoleAssociationsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_role_associations"}, []string{"role_id", "policy_id"}, &iteratorForCreateRoleAssociations{rows: arg})
}

// iteratorForCreateRoles implements pgx.CopyFromSource.
type iteratorForCreateRoles struct {
	rows                 []CreateRolesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateRoles) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateRoles) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Description,
		r.rows[0].CreatedAt,
		r.rows[0].UpdatedAt,
	}, nil
}

func (r iteratorForCreateRoles) Err() error {
	return nil
}

func (q *Queries) CreateRoles(ctx context.Context, arg []CreateRolesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_roles"}, []string{"name", "description", "created_at", "updated_at"}, &iteratorForCreateRoles{rows: arg})
}

// iteratorForCreateStudents implements pgx.CopyFromSource.
type iteratorForCreateStudents struct {
	rows                 []uuid.UUID
	skippedFirstNextCall bool
}

func (r *iteratorForCreateStudents) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateStudents) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0],
	}, nil
}

func (r iteratorForCreateStudents) Err() error {
	return nil
}

func (q *Queries) CreateStudents(ctx context.Context, memberID []uuid.UUID) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_students"}, []string{"member_id"}, &iteratorForCreateStudents{rows: memberID})
}

// iteratorForCreateWorkPositions implements pgx.CopyFromSource.
type iteratorForCreateWorkPositions struct {
	rows                 []CreateWorkPositionsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateWorkPositions) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateWorkPositions) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].OrganizationID,
		r.rows[0].Description,
		r.rows[0].CreatedAt,
		r.rows[0].UpdatedAt,
	}, nil
}

func (r iteratorForCreateWorkPositions) Err() error {
	return nil
}

func (q *Queries) CreateWorkPositions(ctx context.Context, arg []CreateWorkPositionsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_work_positions"}, []string{"name", "organization_id", "description", "created_at", "updated_at"}, &iteratorForCreateWorkPositions{rows: arg})
}

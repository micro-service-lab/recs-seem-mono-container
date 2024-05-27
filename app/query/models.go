// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Absence struct {
	TAbsencesPkey pgtype.Int8 `json:"t_absences_pkey"`
	AbsenceID     uuid.UUID   `json:"absence_id"`
	AttendanceID  uuid.UUID   `json:"attendance_id"`
}

type AttachableItem struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	Alias                string        `json:"alias"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
}

type AttachedMessage struct {
	TAttachedMessagesPkey pgtype.Int8 `json:"t_attached_messages_pkey"`
	AttachedMessageID     uuid.UUID   `json:"attached_message_id"`
	MessageID             uuid.UUID   `json:"message_id"`
	AttachableItemID      pgtype.UUID `json:"attachable_item_id"`
}

type AttendStatus struct {
	MAttendStatusesPkey pgtype.Int8 `json:"m_attend_statuses_pkey"`
	AttendStatusID      uuid.UUID   `json:"attend_status_id"`
	Name                string      `json:"name"`
	Key                 string      `json:"key"`
}

type Attendance struct {
	TAttendancesPkey   pgtype.Int8 `json:"t_attendances_pkey"`
	AttendanceID       uuid.UUID   `json:"attendance_id"`
	AttendanceTypeID   uuid.UUID   `json:"attendance_type_id"`
	MemberID           uuid.UUID   `json:"member_id"`
	Description        string      `json:"description"`
	Date               pgtype.Date `json:"date"`
	MailSendFlag       bool        `json:"mail_send_flag"`
	SendOrganizationID pgtype.UUID `json:"send_organization_id"`
	PostedAt           time.Time   `json:"posted_at"`
	LastEditedAt       time.Time   `json:"last_edited_at"`
}

type AttendanceType struct {
	MAttendanceTypesPkey pgtype.Int8 `json:"m_attendance_types_pkey"`
	AttendanceTypeID     uuid.UUID   `json:"attendance_type_id"`
	Name                 string      `json:"name"`
	Key                  string      `json:"key"`
	Color                string      `json:"color"`
}

type ChatRoom struct {
	MChatRoomsPkey   pgtype.Int8 `json:"m_chat_rooms_pkey"`
	ChatRoomID       uuid.UUID   `json:"chat_room_id"`
	Name             string      `json:"name"`
	IsPrivate        bool        `json:"is_private"`
	CoverImageID     pgtype.UUID `json:"cover_image_id"`
	OwnerID          pgtype.UUID `json:"owner_id"`
	FromOrganization bool        `json:"from_organization"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

type ChatRoomBelonging struct {
	MChatRoomBelongingsPkey pgtype.Int8 `json:"m_chat_room_belongings_pkey"`
	MemberID                uuid.UUID   `json:"member_id"`
	ChatRoomID              uuid.UUID   `json:"chat_room_id"`
	AddedAt                 time.Time   `json:"added_at"`
}

type EarlyLeaving struct {
	TEarlyLeavingsPkey pgtype.Int8 `json:"t_early_leavings_pkey"`
	EarlyLeavingID     uuid.UUID   `json:"early_leaving_id"`
	AttendanceID       uuid.UUID   `json:"attendance_id"`
	LeaveTime          time.Time   `json:"leave_time"`
}

type Event struct {
	TEventsPkey        pgtype.Int8 `json:"t_events_pkey"`
	EventID            uuid.UUID   `json:"event_id"`
	EventTypeID        uuid.UUID   `json:"event_type_id"`
	Title              string      `json:"title"`
	Description        pgtype.Text `json:"description"`
	OrganizationID     pgtype.UUID `json:"organization_id"`
	StartTime          time.Time   `json:"start_time"`
	EndTime            time.Time   `json:"end_time"`
	MailSendFlag       bool        `json:"mail_send_flag"`
	SendOrganizationID pgtype.UUID `json:"send_organization_id"`
	PostedBy           pgtype.UUID `json:"posted_by"`
	LastEditedBy       pgtype.UUID `json:"last_edited_by"`
	PostedAt           time.Time   `json:"posted_at"`
	LastEditedAt       time.Time   `json:"last_edited_at"`
}

type EventType struct {
	MEventTypesPkey pgtype.Int8 `json:"m_event_types_pkey"`
	EventTypeID     uuid.UUID   `json:"event_type_id"`
	Name            string      `json:"name"`
	Key             string      `json:"key"`
	Color           string      `json:"color"`
}

type File struct {
	TFilesPkey       pgtype.Int8 `json:"t_files_pkey"`
	FileID           uuid.UUID   `json:"file_id"`
	AttachableItemID uuid.UUID   `json:"attachable_item_id"`
}

type Grade struct {
	MGradesPkey    pgtype.Int8 `json:"m_grades_pkey"`
	GradeID        uuid.UUID   `json:"grade_id"`
	Key            string      `json:"key"`
	OrganizationID uuid.UUID   `json:"organization_id"`
}

type Group struct {
	MGroupsPkey    pgtype.Int8 `json:"m_groups_pkey"`
	GroupID        uuid.UUID   `json:"group_id"`
	Key            string      `json:"key"`
	OrganizationID uuid.UUID   `json:"organization_id"`
}

type Image struct {
	TImagesPkey      pgtype.Int8   `json:"t_images_pkey"`
	ImageID          uuid.UUID     `json:"image_id"`
	Height           pgtype.Float8 `json:"height"`
	Width            pgtype.Float8 `json:"width"`
	AttachableItemID uuid.UUID     `json:"attachable_item_id"`
}

type LabIOHistory struct {
	TLabIoHistoriesPkey pgtype.Int8        `json:"t_lab_io_histories_pkey"`
	LabIoHistoryID      uuid.UUID          `json:"lab_io_history_id"`
	MemberID            uuid.UUID          `json:"member_id"`
	EnteredAt           time.Time          `json:"entered_at"`
	ExitedAt            pgtype.Timestamptz `json:"exited_at"`
}

type LateArrival struct {
	TLateArrivalsPkey pgtype.Int8 `json:"t_late_arrivals_pkey"`
	LateArrivalID     uuid.UUID   `json:"late_arrival_id"`
	AttendanceID      uuid.UUID   `json:"attendance_id"`
	ArriveTime        time.Time   `json:"arrive_time"`
}

type Member struct {
	MMembersPkey           pgtype.Int8 `json:"m_members_pkey"`
	MemberID               uuid.UUID   `json:"member_id"`
	LoginID                string      `json:"login_id"`
	Password               string      `json:"password"`
	Email                  string      `json:"email"`
	Name                   string      `json:"name"`
	FirstName              string      `json:"first_name"`
	LastName               string      `json:"last_name"`
	AttendStatusID         uuid.UUID   `json:"attend_status_id"`
	ProfileImageID         pgtype.UUID `json:"profile_image_id"`
	GradeID                uuid.UUID   `json:"grade_id"`
	GroupID                uuid.UUID   `json:"group_id"`
	PersonalOrganizationID uuid.UUID   `json:"personal_organization_id"`
	RoleID                 pgtype.UUID `json:"role_id"`
	CreatedAt              time.Time   `json:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at"`
}

type Membership struct {
	MMembershipsPkey pgtype.Int8 `json:"m_memberships_pkey"`
	MemberID         uuid.UUID   `json:"member_id"`
	OrganizationID   uuid.UUID   `json:"organization_id"`
	WorkPositionID   pgtype.UUID `json:"work_position_id"`
	AddedAt          time.Time   `json:"added_at"`
}

type Message struct {
	TMessagesPkey pgtype.Int8 `json:"t_messages_pkey"`
	MessageID     uuid.UUID   `json:"message_id"`
	ChatRoomID    uuid.UUID   `json:"chat_room_id"`
	SenderID      pgtype.UUID `json:"sender_id"`
	Body          string      `json:"body"`
	PostedAt      time.Time   `json:"posted_at"`
	LastEditedAt  time.Time   `json:"last_edited_at"`
}

type MimeType struct {
	MMimeTypesPkey pgtype.Int8 `json:"m_mime_types_pkey"`
	MimeTypeID     uuid.UUID   `json:"mime_type_id"`
	Name           string      `json:"name"`
	Kind           string      `json:"kind"`
	Key            string      `json:"key"`
}

type Organization struct {
	MOrganizationsPkey pgtype.Int8 `json:"m_organizations_pkey"`
	OrganizationID     uuid.UUID   `json:"organization_id"`
	Name               string      `json:"name"`
	Description        pgtype.Text `json:"description"`
	Color              pgtype.Text `json:"color"`
	IsPersonal         bool        `json:"is_personal"`
	IsWhole            bool        `json:"is_whole"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	ChatRoomID         pgtype.UUID `json:"chat_room_id"`
}

type Permission struct {
	MPermissionsPkey     pgtype.Int8 `json:"m_permissions_pkey"`
	PermissionID         uuid.UUID   `json:"permission_id"`
	Name                 string      `json:"name"`
	Description          string      `json:"description"`
	Key                  string      `json:"key"`
	PermissionCategoryID uuid.UUID   `json:"permission_category_id"`
}

type PermissionAssociation struct {
	MPermissionAssociationsPkey pgtype.Int8 `json:"m_permission_associations_pkey"`
	PermissionID                uuid.UUID   `json:"permission_id"`
	WorkPositionID              uuid.UUID   `json:"work_position_id"`
}

type PermissionCategory struct {
	MPermissionCategoriesPkey pgtype.Int8 `json:"m_permission_categories_pkey"`
	PermissionCategoryID      uuid.UUID   `json:"permission_category_id"`
	Name                      string      `json:"name"`
	Description               string      `json:"description"`
	Key                       string      `json:"key"`
}

type Policy struct {
	MPoliciesPkey    pgtype.Int8 `json:"m_policies_pkey"`
	PolicyID         uuid.UUID   `json:"policy_id"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	Key              string      `json:"key"`
	PolicyCategoryID uuid.UUID   `json:"policy_category_id"`
}

type PolicyCategory struct {
	MPolicyCategoriesPkey pgtype.Int8 `json:"m_policy_categories_pkey"`
	PolicyCategoryID      uuid.UUID   `json:"policy_category_id"`
	Name                  string      `json:"name"`
	Description           string      `json:"description"`
	Key                   string      `json:"key"`
}

type PositionHistory struct {
	TPositionHistoriesPkey pgtype.Int8 `json:"t_position_histories_pkey"`
	PositionHistoryID      uuid.UUID   `json:"position_history_id"`
	MemberID               uuid.UUID   `json:"member_id"`
	XPos                   float64     `json:"x_pos"`
	YPos                   float64     `json:"y_pos"`
	SentAt                 time.Time   `json:"sent_at"`
}

type Professor struct {
	MProfessorsPkey pgtype.Int8 `json:"m_professors_pkey"`
	ProfessorID     uuid.UUID   `json:"professor_id"`
	MemberID        uuid.UUID   `json:"member_id"`
}

type ReadReceipt struct {
	TReadReceiptsPkey pgtype.Int8        `json:"t_read_receipts_pkey"`
	MemberID          uuid.UUID          `json:"member_id"`
	MessageID         uuid.UUID          `json:"message_id"`
	ReadAt            pgtype.Timestamptz `json:"read_at"`
}

type Record struct {
	TRecordsPkey   pgtype.Int8 `json:"t_records_pkey"`
	RecordID       uuid.UUID   `json:"record_id"`
	RecordTypeID   uuid.UUID   `json:"record_type_id"`
	Title          string      `json:"title"`
	Body           pgtype.Text `json:"body"`
	OrganizationID pgtype.UUID `json:"organization_id"`
	PostedBy       pgtype.UUID `json:"posted_by"`
	LastEditedBy   pgtype.UUID `json:"last_edited_by"`
	PostedAt       time.Time   `json:"posted_at"`
	LastEditedAt   time.Time   `json:"last_edited_at"`
}

type RecordType struct {
	MRecordTypesPkey pgtype.Int8 `json:"m_record_types_pkey"`
	RecordTypeID     uuid.UUID   `json:"record_type_id"`
	Name             string      `json:"name"`
	Key              string      `json:"key"`
}

type Role struct {
	MRolesPkey  pgtype.Int8 `json:"m_roles_pkey"`
	RoleID      uuid.UUID   `json:"role_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type RoleAssociation struct {
	MRoleAssociationsPkey pgtype.Int8 `json:"m_role_associations_pkey"`
	RoleID                uuid.UUID   `json:"role_id"`
	PolicyID              uuid.UUID   `json:"policy_id"`
}

type Student struct {
	MStudentsPkey pgtype.Int8 `json:"m_students_pkey"`
	StudentID     uuid.UUID   `json:"student_id"`
	MemberID      uuid.UUID   `json:"member_id"`
}

type WorkPosition struct {
	MWorkPositionsPkey pgtype.Int8 `json:"m_work_positions_pkey"`
	WorkPositionID     uuid.UUID   `json:"work_position_id"`
	OrganizationID     uuid.UUID   `json:"organization_id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

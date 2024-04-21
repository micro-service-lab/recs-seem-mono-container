// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: copyfrom.go

package query

import (
	"context"
)

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
		r.rows[0].IsPersonal,
		r.rows[0].IsWhole,
		r.rows[0].CreatedAt,
		r.rows[0].UpdatedAt,
	}, nil
}

func (r iteratorForCreateOrganizations) Err() error {
	return nil
}

func (q *Queries) CreateOrganizations(ctx context.Context, arg []CreateOrganizationsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_organizations"}, []string{"name", "description", "is_personal", "is_whole", "created_at", "updated_at"}, &iteratorForCreateOrganizations{rows: arg})
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

// CREATE TABLE m_permissions (
//
//		m_permissions_pkey BIGSERIAL,
//	    permission_id UUID NOT NULL DEFAULT uuid_generate_v4(),
//	    name VARCHAR(255) NOT NULL,
//	    description TEXT NOT NULL,
//		key VARCHAR(255) NOT NULL,
//		permission_category_id UUID NOT NULL
//
// );
// ALTER TABLE m_permissions ADD CONSTRAINT m_permissions_pkey PRIMARY KEY (m_permissions_pkey);
// ALTER TABLE m_permissions ADD CONSTRAINT fk_m_permissions_permission_category_id FOREIGN KEY (permission_category_id) REFERENCES m_permission_categories(permission_category_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
// CREATE UNIQUE INDEX idx_m_permissions_id ON m_permissions(permission_id);
// CREATE UNIQUE INDEX idx_m_permissions_key ON m_permissions(key);
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
		r.rows[0].Description,
		r.rows[0].CreatedAt,
		r.rows[0].UpdatedAt,
	}, nil
}

func (r iteratorForCreateWorkPositions) Err() error {
	return nil
}

func (q *Queries) CreateWorkPositions(ctx context.Context, arg []CreateWorkPositionsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"m_work_positions"}, []string{"name", "description", "created_at", "updated_at"}, &iteratorForCreateWorkPositions{rows: arg})
}

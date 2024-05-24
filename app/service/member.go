package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
)

// ManageMember メンバー管理サービス。
type ManageMember struct {
	DB         store.Store
	Translator i18n.Translation
}

// CreateMember メンバーを作成する。
func (m *ManageMember) CreateMember(
	ctx context.Context,
	loginID,
	password,
	email,
	name,
	firstName,
	lastName string,
	attendStatusID,
	gradeID,
	groupID uuid.UUID,
	profileImageID,
	roleID entity.UUID,
	personalOrganizationID uuid.UUID,
) (e entity.Member, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	orgDscStr := m.Translator.TranslateWithOpts(
		lang.GetLocaleForTranslation(ctx), "OrganizationDescriptionForPersonal", i18n.Options{
			DefaultMessage: &i18n.Message{
				ID:    "OrganizationDescriptionForPersonal",
				Other: "Personal organization",
			},
			TemplateData: map[string]any{
				"UserName": name,
			},
		})
	orgDsc := entity.String{
		Valid:  true,
		String: orgDscStr,
	}
	col := entity.String{
		Valid:  true,
		String: RandomColor(),
	}
	_, err = m.DB.CreateOrganizationWithSd(ctx, sd, parameter.CreateOrganizationParam{
		Name:        name,
		Description: orgDsc,
		Color:       col,
		IsPersonal:  true,
		IsWhole:     false,
		ChatRoomID:  entity.UUID{},
	})
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to create organization: %w", err)
	}
	p := parameter.CreateMemberParam{
		// Name:        name,
		// Description: description,
		// Color:       color,
		// IsPersonal:  false,
		// IsWhole:     false,
		// ChatRoomID: entity.UUID{
		// 	Valid: true,
		// 	Bytes: cr.ChatRoomID,
		// },
	}
	e, err = m.DB.CreateMemberWithSd(ctx, sd, p)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to create organization: %w", err)
	}
	return e, nil
}

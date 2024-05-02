package service

import (
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// ManageAbsence 欠席管理サービス。
type ManageAbsence struct {
	DB store.Store
}

package repo

import (
	"context"
	"qn-base/app/admin/internal/biz/systemuser"
)

// SystemUserRepo is a SystemUser repo.
//go:generate mockgen -source=./system_user_repo.go -destination=./mocks/system_user_repo_mock.go -package=repo
type SystemUserRepo interface {
	Save(context.Context, *systemuser.SystemUser) (*systemuser.SystemUser, error)
	Update(context.Context, *systemuser.SystemUser) (*systemuser.SystemUser, error)
	Delete(context.Context, string) error
	BatchDelete(context.Context, []string) (int32, int32, []string, error)
	FindByID(context.Context, string) (*systemuser.SystemUser, error)
	FindByUsername(context.Context, string) (*systemuser.SystemUser, error)
	FindByEmail(context.Context, string) (*systemuser.SystemUser, error)
	FindByMobile(context.Context, string) (*systemuser.SystemUser, error)
	ListSystemUsers(context.Context, *systemuser.ListUserRequest) ([]*systemuser.SystemUser, int32, error)
	ChangeStatus(context.Context, string, int8) error
	GetUserStats(context.Context, string) (*systemuser.UserStats, error)
}

// ListUserRequest is a list user request.
type ListUserRequest struct {
	Page      int32
	PageSize  int32
	Username  string
	Email     string
	Mobile    string
	Status    *int8
	DeptID    string
	StartDate string
	EndDate   string
	TenantID  string
}
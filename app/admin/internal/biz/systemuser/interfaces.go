package systemuser

import (
	"context"
	"time"
)


//go:generate mockgen -source=interfaces.go -destination=./mocks/usecase_mock.go -package=mocks
type UserUsecase interface {
	CreateUser(ctx context.Context, u *SystemUser) (*SystemUser, error)
	GetUser(ctx context.Context, id string) (*SystemUser, error)
	UpdateUser(ctx context.Context, u *SystemUser) (*SystemUser, error)
	DeleteUser(ctx context.Context, id string) error
	BatchDeleteUsers(ctx context.Context, ids []string) (*BatchDeleteResult, error)
	ChangeUserStatus(ctx context.Context, id string, status int8) error
	ResetPassword(ctx context.Context, id, newPassword string) error
	CheckAccountExists(ctx context.Context, account string) (bool, error)
	GetUserStats(ctx context.Context, tenantID string) (*UserStats, error)
	ListUsers(ctx context.Context, req *ListUserRequest) ([]*SystemUser, int32, error)
}


// SystemUser is a SystemUser model.
type SystemUser struct {
	ID        *string    `json:"id,omitempty"`         // id
	CreateBy  *string    `json:"create_by,omitempty"`  // 创建人
	CreatedAt *time.Time `json:"created_at,omitempty"` // 创建时间
	UpdateBy  *string    `json:"update_by,omitempty"`  // 更新人
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // 更新时间
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // 删除时间
	TenantID  *string    `json:"tenant_id,omitempty"`  // 租户ID
	Account   *string    `json:"account,omitempty"`    // 用户名
	Password  *string    `json:"password,omitempty"`   // 密码
	Nickname  *string    `json:"nickname,omitempty"`   // 昵称
	Remark    *string    `json:"remark,omitempty"`     // 备注
	DeptID    *string    `json:"dept_id,omitempty"`    // 部门ID
	PostIds   *string    `json:"post_ids,omitempty"`   // 岗位ID
	Email     *string    `json:"email,omitempty"`      // 邮箱
	Mobile    *string    `json:"mobile,omitempty"`     // 手机
	Sex       *int8      `json:"sex,omitempty"`        // 用户性别(0:女 1:男)
	Avatar    *string    `json:"avatar,omitempty"`     // 头像地址
	Status    *int8      `json:"status,omitempty"`     // 帐号状态(0:停用 1:正常)
	LoginIP   *string    `json:"login_ip,omitempty"`   // 登录IP
	LoginDate *time.Time `json:"login_date,omitempty"` // 登录时间
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

// UserStats represents user statistics.
type UserStats struct {
	TotalUsers          int32 `json:"total_users"`
	ActiveUsers         int32 `json:"active_users"`
	InactiveUsers       int32 `json:"inactive_users"`
	TodayRegistered     int32 `json:"today_registered"`
	ThisWeekRegistered  int32 `json:"this_week_registered"`
	ThisMonthRegistered int32 `json:"this_month_registered"`
}

// BatchDeleteResult represents batch delete result.
type BatchDeleteResult struct {
	SuccessCount int32    `json:"success_count"`
	FailedCount  int32    `json:"failed_count"`
	FailedIDs    []string `json:"failed_ids"`
}
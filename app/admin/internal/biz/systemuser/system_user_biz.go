package systemuser

import (
	"context"
	adminV1 "qn-base/api/gen/go/admin/v1"
	"qn-base/pkg/lang/ptr"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound("USER_NOT_FOUND", "user not found")
	// ErrUserAlreadyExists is user already exists.
	ErrUserAlreadyExists = errors.Conflict("USER_ALREADY_EXISTS", "user already exists")
)

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

// SystemUserRepo is a SystemUser repo.
type SystemUserRepo interface {
	Save(context.Context, *SystemUser) (*SystemUser, error)
	Update(context.Context, *SystemUser) (*SystemUser, error)
	Delete(context.Context, string) error
	FindByID(context.Context, string) (*SystemUser, error)
	FindByUsername(context.Context, string) (*SystemUser, error)
	ListSystemUsers(context.Context, *ListUserRequest) ([]*SystemUser, int32, error)
}

// ListUserRequest is a list user request.
type ListUserRequest struct {
	Page     int32
	PageSize int32
	Username string
}

// UserUsecase is a SystemUser usecase.
type UserUsecase struct {
	repo SystemUserRepo
	log  *log.Helper
}

// NewUserUsecase new a SystemUser usecase.
func NewUserUsecase(repo SystemUserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateUser creates a SystemUser, and returns the new SystemUser.
func (uc *UserUsecase) CreateUser(ctx context.Context, u *SystemUser) (*SystemUser, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", u.Account)
	if u.Account == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 检查用户是否已存在
	existingUser, err := uc.repo.FindByUsername(ctx, ptr.From(u.Account))
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	return uc.repo.Save(ctx, u)
}

// GetUser gets a SystemUser by ID.
func (uc *UserUsecase) GetUser(ctx context.Context, id string) (*SystemUser, error) {
	uc.log.WithContext(ctx).Infof("GetUser: %s", id)
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateUser updates a SystemUser.
func (uc *UserUsecase) UpdateUser(ctx context.Context, u *SystemUser) (*SystemUser, error) {
	uc.log.WithContext(ctx).Infof("UpdateUser: %s", u.ID)
	return uc.repo.Update(ctx, u)
}

// DeleteUser deletes a SystemUser by ID.
func (uc *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteUser: %s", id)
	return uc.repo.Delete(ctx, id)
}

// ListUsers lists users.
func (uc *UserUsecase) ListUsers(ctx context.Context, req *ListUserRequest) ([]*SystemUser, int32, error) {
	uc.log.WithContext(ctx).Infof("ListSystemUsers: page=%d, page_size=%d, username=%s", req.Page, req.PageSize, req.Username)
	return uc.repo.ListSystemUsers(ctx, req)
}

package systemuser

import (
	"context"
	"fmt"
	"qn-base/app/admin/internal/biz/systemuser/repo"
	"qn-base/pkg/lang/ptr"
	"qn-base/pkg/util/pswd"
	"qn-base/pkg/util/validator"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound("USER_NOT_FOUND", "user not found")
	// ErrUserAlreadyExists is user already exists.
	ErrUserAlreadyExists = errors.Conflict("USER_ALREADY_EXISTS", "user already exists")
	// ErrEmailAlreadyExists is email already exists.
	ErrEmailAlreadyExists = errors.Conflict("EMAIL_ALREADY_EXISTS", "email already exists")
	// ErrMobileAlreadyExists is mobile already exists.
	ErrMobileAlreadyExists = errors.Conflict("MOBILE_ALREADY_EXISTS", "mobile already exists")
	// ErrInvalidParameter is invalid parameter.
	ErrInvalidParameter = errors.BadRequest("INVALID_PARAMETER", "invalid parameter")
	// ErrPasswordVerifyFailed is password verify failed.
	ErrPasswordVerifyFailed = errors.Unauthorized("PASSWORD_VERIFY_FAILED", "password verify failed")
)



// userUsecase 是 UserUsecase 接口的具体实现
type userUsecase struct {
	repo repo.SystemUserRepo
	log  *log.Helper
}

// 确保 userUsecase 实现了 UserUsecase 接口
var _ UserUsecase = (*userUsecase)(nil)

// NewUserUsecase new a SystemUser usecase.
func NewUserUsecase(repo repo.SystemUserRepo, logger log.Logger) UserUsecase {\r
	return &userUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateUser creates a SystemUser, and returns the new SystemUser.
func (uc *userUsecase) CreateUser(ctx context.Context, u *SystemUser) (*SystemUser, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", u.Account)

	// 参数校验
	if err := uc.validateCreateUser(u); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	existingUser, err := uc.repo.FindByUsername(ctx, ptr.From(u.Account))
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 检查邮箱是否已存在（如果提供了邮箱）
	if u.Email != nil && *u.Email != "" {
		existingEmail, err := uc.repo.FindByEmail(ctx, *u.Email)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
		if existingEmail != nil {
			return nil, ErrEmailAlreadyExists
		}
	}

	// 检查手机号是否已存在（如果提供了手机号）
	if u.Mobile != nil && *u.Mobile != "" {
		existingMobile, err := uc.repo.FindByMobile(ctx, *u.Mobile)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
		if existingMobile != nil {
			return nil, ErrMobileAlreadyExists
		}
	}

	// 密码加密
	if u.Password != nil && *u.Password != "" {
		hashedPassword, err := pswd.HashPassword(*u.Password)
		if err != nil {
			return nil, fmt.Errorf("密码加密失败: %w", err)
		}
		u.Password = &hashedPassword
	}

	// 设置默认状态
	if u.Status == nil {
		u.Status = ptr.Of(int8(1)) // 默认正常状态
	}

	return uc.repo.Save(ctx, u)
}

// GetUser gets a SystemUser by ID.
func (uc *userUsecase) GetUser(ctx context.Context, id string) (*SystemUser, error) {
	uc.log.WithContext(ctx).Infof("GetUser: %s", id)
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateUser updates a SystemUser.
func (uc *userUsecase) UpdateUser(ctx context.Context, u *SystemUser) (*SystemUser, error) {
	uc.log.WithContext(ctx).Infof("UpdateUser: %s", ptr.From(u.ID))

	// 参数校验
	if err := uc.validateUpdateUser(u); err != nil {
		return nil, err
	}

	// 检查用户是否存在
	existingUser, err := uc.repo.FindByID(ctx, ptr.From(u.ID))
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, ErrUserNotFound
	}

	// 检查邮箱是否被其他用户使用
	if u.Email != nil && *u.Email != "" && (existingUser.Email == nil || *existingUser.Email != *u.Email) {
		existingEmail, err := uc.repo.FindByEmail(ctx, *u.Email)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
		if existingEmail != nil && *existingEmail.ID != *u.ID {
			return nil, ErrEmailAlreadyExists
		}
	}

	// 检查手机号是否被其他用户使用
	if u.Mobile != nil && *u.Mobile != "" && (existingUser.Mobile == nil || *existingUser.Mobile != *u.Mobile) {
		existingMobile, err := uc.repo.FindByMobile(ctx, *u.Mobile)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
		if existingMobile != nil && *existingMobile.ID != *u.ID {
			return nil, ErrMobileAlreadyExists
		}
	}

	return uc.repo.Update(ctx, u)
}

// DeleteUser deletes a SystemUser by ID.
func (uc *userUsecase) DeleteUser(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteUser: %s", id)
	return uc.repo.Delete(ctx, id)
}

// ListUsers lists users.
func (uc *userUsecase) ListUsers(ctx context.Context, req *ListUserRequest) ([]*SystemUser, int32, error) {
	uc.log.WithContext(ctx).Infof("ListSystemUsers: page=%d, page_size=%d, username=%s", req.Page, req.PageSize, req.Username)

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 15
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	return uc.repo.ListSystemUsers(ctx, req)
}

// BatchDeleteUsers deletes multiple users.
func (uc *userUsecase) BatchDeleteUsers(ctx context.Context, ids []string) (*BatchDeleteResult, error) {
	uc.log.WithContext(ctx).Infof("BatchDeleteUsers: ids=%v", ids)

	// 参数校验
	if err := validator.ValidateIDs(ids); err != nil {
		return nil, errors.BadRequest("INVALID_PARAMETER", err.Error())
	}

	successCount, failedCount, failedIDs, err := uc.repo.BatchDelete(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &BatchDeleteResult{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		FailedIDs:    failedIDs,
	}, nil
}

// ChangeUserStatus changes user status.
func (uc *userUsecase) ChangeUserStatus(ctx context.Context, id string, status int8) error {
	uc.log.WithContext(ctx).Infof("ChangeUserStatus: id=%s, status=%d", id, status)

	// 参数校验
	if err := validator.ValidateRequiredString(id, "用户ID"); err != nil {
		return errors.BadRequest("INVALID_PARAMETER", err.Error())
	}
	if err := validator.ValidateStatus(status); err != nil {
		return errors.BadRequest("INVALID_PARAMETER", err.Error())
	}

	// 检查用户是否存在
	existingUser, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrUserNotFound
	}

	return uc.repo.ChangeStatus(ctx, id, status)
}

// ResetPassword resets user password.
func (uc *userUsecase) ResetPassword(ctx context.Context, id, newPassword string) error {
	uc.log.WithContext(ctx).Infof("ResetPassword: id=%s", id)

	// 参数校验
	if err := validator.ValidateRequiredString(id, "用户ID"); err != nil {
		return errors.BadRequest("INVALID_PARAMETER", err.Error())
	}
	if err := validator.ValidatePassword(newPassword); err != nil {
		return errors.BadRequest("INVALID_PARAMETER", err.Error())
	}

	// 检查用户是否存在
	existingUser, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrUserNotFound
	}

	// 密码加密
	hashedPassword, err := pswd.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	updateUser := &SystemUser{
		ID:       &id,
		Password: &hashedPassword,
	}

	_, err = uc.repo.Update(ctx, updateUser)
	return err
}

// CheckAccountExists checks if account exists.
func (uc *userUsecase) CheckAccountExists(ctx context.Context, account string) (bool, error) {
	uc.log.WithContext(ctx).Infof("CheckAccountExists: account=%s", account)

	// 参数校验
	if err := validator.ValidateUsername(account); err != nil {
		return false, errors.BadRequest("INVALID_PARAMETER", err.Error())
	}

	existingUser, err := uc.repo.FindByUsername(ctx, account)
	if err != nil && !errors.IsNotFound(err) {
		return false, err
	}

	return existingUser != nil, nil
}

// GetUserStats gets user statistics.
func (uc *userUsecase) GetUserStats(ctx context.Context, tenantID string) (*UserStats, error) {
	uc.log.WithContext(ctx).Infof("GetUserStats: tenantID=%s", tenantID)
	return uc.repo.GetUserStats(ctx, tenantID)
}

// validateCreateUser validates create user parameters.
func (uc *userUsecase) validateCreateUser(u *SystemUser) error {
	if u.Account == nil {
		return errors.BadRequest("INVALID_PARAMETER", "用户名不能为空")
	}
	if err := validator.ValidateUsername(*u.Account); err != nil {
		return errors.BadRequest("INVALID_PARAMETER", err.Error())
	}

	if u.Password == nil {
		return errors.BadRequest("INVALID_PARAMETER", "密码不能为空")
	}
	if err := validator.ValidatePassword(*u.Password); err != nil {
		return errors.BadRequest("INVALID_PARAMETER", err.Error())
	}

	if u.Nickname != nil {
		if err := validator.ValidateNickname(*u.Nickname); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	if u.Email != nil {
		if err := validator.ValidateEmail(*u.Email); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	if u.Mobile != nil {
		if err := validator.ValidateMobile(*u.Mobile); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	if u.Sex != nil {
		if err := validator.ValidateSex(*u.Sex); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	if u.Status != nil {
		if err := validator.ValidateStatus(*u.Status); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	if u.Remark != nil {
		if err := validator.ValidateRemark(*u.Remark); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	return nil
}

// validateUpdateUser validates update user parameters.
func (uc *userUsecase) validateUpdateUser(u *SystemUser) error {
	if u.ID == nil {
		return errors.BadRequest("INVALID_PARAMETER", "用户ID不能为空")
	}
	if err := validator.ValidateRequiredString(*u.ID, "用户ID"); err != nil {
		return errors.BadRequest("INVALID_PARAMETER", err.Error())
	}

	if u.Nickname != nil {
		if err := validator.ValidateNickname(*u.Nickname); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	if u.Email != nil {
		if err := validator.ValidateEmail(*u.Email); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	if u.Mobile != nil {
		if err := validator.ValidateMobile(*u.Mobile); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	if u.Sex != nil {
		if err := validator.ValidateSex(*u.Sex); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	if u.Status != nil {
		if err := validator.ValidateStatus(*u.Status); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	if u.Remark != nil {
		if err := validator.ValidateRemark(*u.Remark); err != nil {
			return errors.BadRequest("INVALID_PARAMETER", err.Error())
		}
	}

	return nil
}

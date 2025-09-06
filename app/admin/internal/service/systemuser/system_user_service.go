package systemuser

import (
	"context"
	v1 "qn-base/api/gen/go/admin/v1"
	"qn-base/app/admin/internal/biz/systemuser"
	"qn-base/pkg/lang/conv"
	"qn-base/pkg/lang/ptr"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// UserUsecase defines the interface for user business logic.
type UserUsecase interface {
	CreateUser(ctx context.Context, u *systemuser.SystemUser) (*systemuser.SystemUser, error)
	GetUser(ctx context.Context, id string) (*systemuser.SystemUser, error)
	UpdateUser(ctx context.Context, u *systemuser.SystemUser) (*systemuser.SystemUser, error)
	DeleteUser(ctx context.Context, id string) error
	BatchDeleteUsers(ctx context.Context, ids []string) (*systemuser.BatchDeleteResult, error)
	ChangeUserStatus(ctx context.Context, id string, status int8) error
	ResetPassword(ctx context.Context, id, newPassword string) error
	CheckAccountExists(ctx context.Context, account string) (bool, error)
	GetUserStats(ctx context.Context, tenantID string) (*systemuser.UserStats, error)
	ListUsers(ctx context.Context, req *systemuser.ListUserRequest) ([]*systemuser.SystemUser, int32, error)
}

// UserService is a user service.
type UserService struct {
	v1.UnimplementedUserServer

	uc  UserUsecase
	log *log.Helper
}

// NewUserService new a user service.
func NewUserService(logger log.Logger, uc UserUsecase) *UserService {
	l := log.NewHelper(log.With(logger, "module", "admin/service/user-service"))
	return &UserService{
		uc:  uc,
		log: l,
	}
}

// CreateUser implements admin.UserServer.
func (s *UserService) CreateUser(ctx context.Context, in *v1.CreateUserRequest) (*v1.CreateUserReply, error) {
	s.log.WithContext(ctx).Infof("CreateUser: %v", in.Account)

	// 转换时间字段
	var loginDate *time.Time
	if in.LoginDate != nil {
		toTime, err := conv.DefaultStrToTime(*in.LoginDate)
		if err != nil {
			return nil, v1.ErrorBadRequest("invalid parameter")
		}
		loginDate = ptr.Of(toTime)
	}

	user, err := s.uc.CreateUser(ctx, &systemuser.SystemUser{
		Account:   &in.Account,
		Password:  &in.Password,
		Nickname:  in.Nickname,
		Remark:    in.Remark,
		DeptID:    in.DeptId,
		PostIds:   in.PostIds,
		Email:     in.Email,
		Mobile:    in.Mobile,
		Sex:       ptr.Of(int8(ptr.From(in.Sex))),
		Avatar:    in.Avatar,
		LoginIP:   in.LoginIp,
		LoginDate: loginDate,
		TenantID:  in.TenantId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateUserReply{
		User: s.convertToUserInfo(user),
	}, nil
}

// GetUser implements admin.UserServer.
func (s *UserService) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.GetUserReply, error) {
	s.log.WithContext(ctx).Infof("GetUser: %v", in.Id)

	user, err := s.uc.GetUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &v1.GetUserReply{
		User: s.convertToUserInfo(user),
	}, nil
}

// UpdateUser implements admin.UserServer.
func (s *UserService) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UpdateUserReply, error) {
	s.log.WithContext(ctx).Infof("UpdateUser: %v", in.Id)

	// 转换时间字段
	var loginDate *time.Time
	if in.LoginDate != nil {
		toTime, err := conv.DefaultStrToTime(*in.LoginDate)
		if err != nil {
			return nil, v1.ErrorBadRequest("invalid login_date parameter")
		}
		loginDate = ptr.Of(toTime)
	}

	updateUser := &systemuser.SystemUser{
		ID:        &in.Id,
		Nickname:  in.Nickname,
		Remark:    in.Remark,
		DeptID:    in.DeptId,
		PostIds:   in.PostIds,
		Email:     in.Email,
		Mobile:    in.Mobile,
		Sex:       ptr.Of(int8(ptr.From(in.Sex))),
		Avatar:    in.Avatar,
		Status:    ptr.Of(int8(ptr.From(in.Status))),
		LoginIP:   in.LoginIp,
		LoginDate: loginDate,
		TenantID:  in.TenantId,
	}

	user, err := s.uc.UpdateUser(ctx, updateUser)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserReply{
		User: s.convertToUserInfo(user),
	}, nil
}

// DeleteUser implements admin.UserServer.
func (s *UserService) DeleteUser(ctx context.Context, in *v1.DeleteUserRequest) (*v1.DeleteUserReply, error) {
	s.log.WithContext(ctx).Infof("DeleteUser: %v", in.Id)

	err := s.uc.DeleteUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteUserReply{
		Success: true,
	}, nil
}

// BatchDeleteUsers implements admin.UserServer.
func (s *UserService) BatchDeleteUsers(ctx context.Context, in *v1.BatchDeleteUsersRequest) (*v1.BatchDeleteUsersReply, error) {
	s.log.WithContext(ctx).Infof("BatchDeleteUsers: ids=%v", in.Ids)

	result, err := s.uc.BatchDeleteUsers(ctx, in.Ids)
	if err != nil {
		return nil, err
	}

	return &v1.BatchDeleteUsersReply{
		SuccessCount: result.SuccessCount,
		FailedCount:  result.FailedCount,
		FailedIds:    result.FailedIDs,
	}, nil
}

// ChangeUserStatus implements admin.UserServer.
func (s *UserService) ChangeUserStatus(ctx context.Context, in *v1.ChangeUserStatusRequest) (*v1.ChangeUserStatusReply, error) {
	s.log.WithContext(ctx).Infof("ChangeUserStatus: id=%s, status=%d", in.Id, in.Status)

	err := s.uc.ChangeUserStatus(ctx, in.Id, int8(in.Status))
	if err != nil {
		return nil, err
	}

	return &v1.ChangeUserStatusReply{
		Success: true,
	}, nil
}

// ResetPassword implements admin.UserServer.
func (s *UserService) ResetPassword(ctx context.Context, in *v1.ResetPasswordRequest) (*v1.ResetPasswordReply, error) {
	s.log.WithContext(ctx).Infof("ResetPassword: id=%s", in.Id)

	err := s.uc.ResetPassword(ctx, in.Id, in.NewPassword)
	if err != nil {
		return nil, err
	}

	return &v1.ResetPasswordReply{
		Success: true,
	}, nil
}

// CheckAccountExists implements admin.UserServer.
func (s *UserService) CheckAccountExists(ctx context.Context, in *v1.CheckAccountExistsRequest) (*v1.CheckAccountExistsReply, error) {
	s.log.WithContext(ctx).Infof("CheckAccountExists: account=%s", in.Account)

	exists, err := s.uc.CheckAccountExists(ctx, in.Account)
	if err != nil {
		return nil, err
	}

	return &v1.CheckAccountExistsReply{
		Exists: exists,
	}, nil
}

// GetUserStats implements admin.UserServer.
func (s *UserService) GetUserStats(ctx context.Context, in *v1.GetUserStatsRequest) (*v1.GetUserStatsReply, error) {
	s.log.WithContext(ctx).Infof("GetUserStats: tenantID=%s", ptr.From(in.TenantId))

	stats, err := s.uc.GetUserStats(ctx, ptr.From(in.TenantId))
	if err != nil {
		return nil, err
	}

	return &v1.GetUserStatsReply{
		Stats: &v1.UserStats{
			TotalUsers:          stats.TotalUsers,
			ActiveUsers:         stats.ActiveUsers,
			InactiveUsers:       stats.InactiveUsers,
			TodayRegistered:     stats.TodayRegistered,
			ThisWeekRegistered:  stats.ThisWeekRegistered,
			ThisMonthRegistered: stats.ThisMonthRegistered,
		},
	}, nil
}

// convertToUserInfo converts SystemUser to UserInfo.
func (s *UserService) convertToUserInfo(user *systemuser.SystemUser) *v1.UserInfo {
	if user == nil {
		return nil
	}

	return &v1.UserInfo{
		Id:        ptr.From(user.ID),
		Account:   ptr.From(user.Account),
		Nickname:  ptr.From(user.Nickname),
		Remark:    ptr.From(user.Remark),
		DeptId:    ptr.From(user.DeptID),
		PostIds:   ptr.From(user.PostIds),
		Email:     ptr.From(user.Email),
		Mobile:    ptr.From(user.Mobile),
		Sex:       int32(ptr.From(user.Sex)),
		Avatar:    ptr.From(user.Avatar),
		Status:    int32(ptr.From(user.Status)),
		LoginIp:   ptr.From(user.LoginIP),
		LoginDate: conv.TimeToDefaultStr(ptr.From(user.LoginDate)),
		TenantId:  ptr.From(user.TenantID),
		CreatedAt: conv.TimeToDefaultStr(ptr.From(user.CreatedAt)),
		UpdatedAt: conv.TimeToDefaultStr(ptr.From(user.UpdatedAt)),
		CreatedBy: ptr.From(user.CreateBy),
		UpdatedBy: ptr.From(user.UpdateBy),
	}
}

// ListUsers implements admin.UserServer.
func (s *UserService) ListUsers(ctx context.Context, in *v1.ListUsersRequest) (*v1.ListUsersReply, error) {
	s.log.WithContext(ctx).Infof("ListUsers: page=%d, page_size=%d", ptr.From(in.Page), ptr.From(in.PageSize))

	req := &systemuser.ListUserRequest{
		Page:      ptr.From(in.Page),
		PageSize:  ptr.From(in.PageSize),
		Username:  ptr.From(in.Account),
		Email:     ptr.From(in.Email),
		Mobile:    ptr.From(in.Mobile),
		DeptID:    ptr.From(in.DeptId),
		StartDate: ptr.From(in.StartDate),
		EndDate:   ptr.From(in.EndDate),
	}

	if in.Status != nil {
		status := int8(*in.Status)
		req.Status = &status
	}

	users, total, err := s.uc.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	userInfos := make([]*v1.UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = s.convertToUserInfo(user)
	}

	return &v1.ListUsersReply{
		Users: userInfos,
		Total: total,
	}, nil
}

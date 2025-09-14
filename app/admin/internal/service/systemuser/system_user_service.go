package systemuser

import (
	"context"
	v1 "qn-base/api/gen/go/admin/v1"
	"qn-base/app/admin/internal/biz/systemuser"
	"qn-base/app/admin/internal/service/systemuser/convertor"
	"qn-base/pkg/lang/ptr"

	"github.com/go-kratos/kratos/v2/log"
)

// UserService is a user service.
type UserService struct {
	v1.UnimplementedUserServer

	uc  systemuser.UserUsecase
	log *log.Helper
}

// NewUserService new a user service.
func NewUserService(logger log.Logger, uc systemuser.UserUsecase) *UserService {
	l := log.NewHelper(log.With(logger, "module", "admin/service/user-service"))
	return &UserService{
		uc:  uc,
		log: l,
	}
}

// CreateUser implements admin.UserServer.
func (s *UserService) CreateUser(ctx context.Context, in *v1.CreateUserRequest) (*v1.CreateUserReply, error) {
	s.log.WithContext(ctx).Infof("CreateUser: %v", in.Account)
	v1.UserHTTPClientImpl
	// 使用转换函数将请求转换为biz层对象
	bizUser, err := convertor.ToCreateUserBiz(in)
	if err != nil {
		return nil, v1.ErrorBadRequest("invalid parameter")
	}

	user, err := s.uc.CreateUser(ctx, bizUser)
	if err != nil {
		return nil, err
	}

	return &v1.CreateUserReply{
		User: convertor.ToUserInfo(user),
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
		User: convertor.ToUserInfo(user),
	}, nil
}

// UpdateUser implements admin.UserServer.
func (s *UserService) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UpdateUserReply, error) {
	s.log.WithContext(ctx).Infof("UpdateUser: %v", in.Id)

	// 使用转换函数将请求转换为biz层对象
	bizUser, err := convertor.ToUpdateUserBiz(in)
	if err != nil {
		return nil, v1.ErrorBadRequest("invalid login_date parameter")
	}

	user, err := s.uc.UpdateUser(ctx, bizUser)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserReply{
		User: convertor.ToUserInfo(user),
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

	return convertor.ToBatchDeleteResult(result), nil
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
		Stats: convertor.ToUserStats(stats),
	}, nil
}

// ListUsers implements admin.UserServer.
func (s *UserService) ListUsers(ctx context.Context, in *v1.ListUsersRequest) (*v1.ListUsersReply, error) {
	s.log.WithContext(ctx).Infof("ListUsers: page=%d, page_size=%d", ptr.From(in.Page), ptr.From(in.PageSize))

	// 使用转换函数将请求转换为biz层对象
	req := convertor.ToListUserRequestBiz(in)

	users, total, err := s.uc.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	return &v1.ListUsersReply{
		Users: convertor.ToUserInfos(users),
		Total: total,
	}, nil
}

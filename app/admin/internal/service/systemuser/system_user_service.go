package systemuser

import (
	"context"
	"qn-base/api/gen/go/admin/v1"
	"qn-base/app/admin/internal/biz/systemuser"
	"qn-base/pkg/lang/conv"
	"qn-base/pkg/lang/ptr"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// UserService is a user service.
type UserService struct {
	v1.UnimplementedUserServer

	uc  *systemuser.UserUsecase
	log *log.Helper
}

// NewUserService new a user service.
func NewUserService(logger log.Logger, uc *systemuser.UserUsecase) *UserService {
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
		Account:   in.Account,
		Password:  in.Password,
		Nickname:  in.Nickname,
		Remark:    in.Remark,
		DeptID:    in.DeptId,
		PostIds:   in.PostIds,
		Email:     in.Email,
		Mobile:    in.Mobile,
		Sex:       ptr.Of(int8(*in.Sex)),
		Avatar:    in.Avatar,
		LoginIP:   in.LoginIp,
		LoginDate: loginDate,
		TenantID:  in.TenantId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateUserReply{
		User: &v1.UserInfo{
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
		},
	}, nil
}

// GetUser implements admin.UserServer.
func (s *UserService) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.GetUserReply, error) {
	panic("implement me")
}

// UpdateUser implements admin.UserServer.
func (s *UserService) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UpdateUserReply, error) {
	panic("implement me")
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

// ListUsers implements admin.UserServer.
func (s *UserService) ListUsers(ctx context.Context, in *v1.ListUsersRequest) (*v1.ListUsersReply, error) {
	panic("implement me")
}

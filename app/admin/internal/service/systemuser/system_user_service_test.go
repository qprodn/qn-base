package systemuser

import (
	"context"
	"testing"

	v1 "qn-base/api/gen/go/admin/v1"
	bizuser "qn-base/app/admin/internal/biz/systemuser"
	"qn-base/app/admin/internal/service/systemuser/mocks"
	"qn-base/pkg/lang/ptr"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mocks.NewMockUserUsecase(ctrl)
	logger := log.DefaultLogger
	service := NewUserService(logger, mockUc)

	ctx := context.Background()

	t.Run("成功创建用户", func(t *testing.T) {
		req := &v1.CreateUserRequest{
			Account:  "testuser",
			Password: "password123",
			Nickname: ptr.Of("Test User"),
			Email:    ptr.Of("test@example.com"),
			Mobile:   ptr.Of("13800138000"),
			Sex:      ptr.Of(int32(1)),
		}

		expectedUser := &bizuser.SystemUser{
			ID:       ptr.Of("user123"),
			Account:  ptr.Of("testuser"),
			Nickname: ptr.Of("Test User"),
			Email:    ptr.Of("test@example.com"),
			Mobile:   ptr.Of("13800138000"),
			Sex:      ptr.Of(int8(1)),
			Status:   ptr.Of(int8(1)),
		}

		// Mock 期望
		mockUc.EXPECT().
			CreateUser(ctx, gomock.Any()).
			Return(expectedUser, nil)

		// 执行测试
		result, err := service.CreateUser(ctx, req)

		// 断言
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.User)
		assert.Equal(t, "user123", result.User.Id)
		assert.Equal(t, "testuser", result.User.Account)
	})

	t.Run("参数验证失败", func(t *testing.T) {
		req := &v1.CreateUserRequest{
			// 缺少必要参数
		}

		// Mock 期望
		mockUc.EXPECT().
			CreateUser(ctx, gomock.Any()).
			Return(nil, bizuser.ErrInvalidParameter)

		// 执行测试
		result, err := service.CreateUser(ctx, req)

		// 断言
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mocks.NewMockUserUsecase(ctrl)
	logger := log.DefaultLogger
	service := NewUserService(logger, mockUc)

	ctx := context.Background()

	t.Run("成功获取用户", func(t *testing.T) {
		req := &v1.GetUserRequest{
			Id: "user123",
		}

		expectedUser := &bizuser.SystemUser{
			ID:      ptr.Of("user123"),
			Account: ptr.Of("testuser"),
		}

		// Mock 期望
		mockUc.EXPECT().
			GetUser(ctx, "user123").
			Return(expectedUser, nil)

		// 执行测试
		result, err := service.GetUser(ctx, req)

		// 断言
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.User)
		assert.Equal(t, "user123", result.User.Id)
	})

	t.Run("用户不存在", func(t *testing.T) {
		req := &v1.GetUserRequest{
			Id: "nonexistent",
		}

		// Mock 期望
		mockUc.EXPECT().
			GetUser(ctx, "nonexistent").
			Return(nil, bizuser.ErrUserNotFound)

		// 执行测试
		result, err := service.GetUser(ctx, req)

		// 断言
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mocks.NewMockUserUsecase(ctrl)
	logger := log.DefaultLogger
	service := NewUserService(logger, mockUc)

	ctx := context.Background()

	t.Run("成功删除用户", func(t *testing.T) {
		req := &v1.DeleteUserRequest{
			Id: "user123",
		}

		// Mock 期望
		mockUc.EXPECT().
			DeleteUser(ctx, "user123").
			Return(nil)

		// 执行测试
		result, err := service.DeleteUser(ctx, req)

		// 断言
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
	})

	t.Run("删除失败", func(t *testing.T) {
		req := &v1.DeleteUserRequest{
			Id: "user123",
		}

		// Mock 期望
		mockUc.EXPECT().
			DeleteUser(ctx, "user123").
			Return(bizuser.ErrUserNotFound)

		// 执行测试
		result, err := service.DeleteUser(ctx, req)

		// 断言
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUserService_BatchDeleteUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mocks.NewMockUserUsecase(ctrl)
	logger := log.DefaultLogger
	service := NewUserService(logger, mockUc)

	ctx := context.Background()

	t.Run("成功批量删除用户", func(t *testing.T) {
		req := &v1.BatchDeleteUsersRequest{
			Ids: []string{"user1", "user2", "user3"},
		}

		expectedResult := &bizuser.BatchDeleteResult{
			SuccessCount: 3,
			FailedCount:  0,
			FailedIDs:    []string{},
		}

		// Mock 期望
		mockUc.EXPECT().
			BatchDeleteUsers(ctx, req.Ids).
			Return(expectedResult, nil)

		// 执行测试
		result, err := service.BatchDeleteUsers(ctx, req)

		// 断言
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(3), result.SuccessCount)
		assert.Equal(t, int32(0), result.FailedCount)
	})
}

func TestUserService_ChangeUserStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mocks.NewMockUserUsecase(ctrl)
	logger := log.DefaultLogger
	service := NewUserService(logger, mockUc)

	ctx := context.Background()

	t.Run("成功修改用户状态", func(t *testing.T) {
		req := &v1.ChangeUserStatusRequest{
			Id:     "user123",
			Status: 0,
		}

		// Mock 期望
		mockUc.EXPECT().
			ChangeUserStatus(ctx, "user123", int8(0)).
			Return(nil)

		// 执行测试
		result, err := service.ChangeUserStatus(ctx, req)

		// 断言
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
	})
}

func TestUserService_CheckAccountExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mocks.NewMockUserUsecase(ctrl)
	logger := log.DefaultLogger
	service := NewUserService(logger, mockUc)

	ctx := context.Background()

	t.Run("用户名存在", func(t *testing.T) {
		req := &v1.CheckAccountExistsRequest{
			Account: "testuser",
		}

		// Mock 期望
		mockUc.EXPECT().
			CheckAccountExists(ctx, "testuser").
			Return(true, nil)

		// 执行测试
		result, err := service.CheckAccountExists(ctx, req)

		// 断言
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Exists)
	})

	t.Run("用户名不存在", func(t *testing.T) {
		req := &v1.CheckAccountExistsRequest{
			Account: "nonexistent",
		}

		// Mock 期望
		mockUc.EXPECT().
			CheckAccountExists(ctx, "nonexistent").
			Return(false, nil)

		// 执行测试
		result, err := service.CheckAccountExists(ctx, req)

		// 断言
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.Exists)
	})
}

package systemuser_test

import (
	"context"
	"testing"

	"qn-base/app/admin/internal/biz/systemuser"
	"qn-base/app/admin/internal/biz/systemuser/mocks"
	"qn-base/pkg/lang/ptr"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSystemUserRepo(ctrl)
	logger := log.DefaultLogger
	uc := systemuser.NewUserUsecase(mockRepo, logger)

	ctx := context.Background()

	t.Run("成功创建用户", func(t *testing.T) {
		user := &systemuser.SystemUser{
			Account:  ptr.Of("testuser"),
			Password: ptr.Of("password123"),
			Nickname: ptr.Of("Test User"),
			Email:    ptr.Of("test@example.com"),
			Mobile:   ptr.Of("13800138000"),
			Sex:      ptr.Of(int8(1)),
		}

		expectedUser := &systemuser.SystemUser{
			ID:       ptr.Of("user123"),
			Account:  ptr.Of("testuser"),
			Nickname: ptr.Of("Test User"),
			Email:    ptr.Of("test@example.com"),
			Mobile:   ptr.Of("13800138000"),
			Sex:      ptr.Of(int8(1)),
			Status:   ptr.Of(int8(1)),
		}

		// Mock 期望
		mockRepo.EXPECT().
			FindByUsername(ctx, "testuser").
			Return(nil, nil)

		mockRepo.EXPECT().
			FindByEmail(ctx, "test@example.com").
			Return(nil, nil)

		mockRepo.EXPECT().
			FindByMobile(ctx, "13800138000").
			Return(nil, nil)

		mockRepo.EXPECT().
			Save(ctx, gomock.Any()).
			Return(expectedUser, nil)

		// 执行测试
		result, err := uc.CreateUser(ctx, user)

		// 断言
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "user123", *result.ID)
		assert.Equal(t, "testuser", *result.Account)
	})

	t.Run("用户名为空", func(t *testing.T) {
		user := &systemuser.SystemUser{
			Password: ptr.Of("password123"),
		}

		// 执行测试
		result, err := uc.CreateUser(ctx, user)

		// 断言
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "用户名不能为空")
	})

	t.Run("用户名已存在", func(t *testing.T) {
		user := &systemuser.SystemUser{
			Account:  ptr.Of("testuser"),
			Password: ptr.Of("password123"),
		}

		existingUser := &systemuser.SystemUser{
			ID:      ptr.Of("existing123"),
			Account: ptr.Of("testuser"),
		}

		// Mock 期望
		mockRepo.EXPECT().
			FindByUsername(ctx, "testuser").
			Return(existingUser, nil)

		// 执行测试
		result, err := uc.CreateUser(ctx, user)

		// 断言
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, systemuser.ErrUserAlreadyExists))
	})
}

func TestUserUsecase_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSystemUserRepo(ctrl)
	logger := log.DefaultLogger
	uc := systemuser.NewUserUsecase(mockRepo, logger)

	ctx := context.Background()

	t.Run("成功获取用户", func(t *testing.T) {
		userID := "user123"
		expectedUser := &systemuser.SystemUser{
			ID:      ptr.Of(userID),
			Account: ptr.Of("testuser"),
		}

		// Mock 期望
		mockRepo.EXPECT().
			FindByID(ctx, userID).
			Return(expectedUser, nil)

		// 执行测试
		result, err := uc.GetUser(ctx, userID)

		// 断言
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, *result.ID)
	})

	t.Run("用户不存在", func(t *testing.T) {
		userID := "nonexistent"

		// Mock 期望
		mockRepo.EXPECT().
			FindByID(ctx, userID).
			Return(nil, nil)

		// 执行测试
		result, err := uc.GetUser(ctx, userID)

		// 断言
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, systemuser.ErrUserNotFound))
	})
}

func TestUserUsecase_BatchDeleteUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSystemUserRepo(ctrl)
	logger := log.DefaultLogger
	uc := systemuser.NewUserUsecase(mockRepo, logger)

	ctx := context.Background()

	t.Run("成功批量删除用户", func(t *testing.T) {
		ids := []string{"user1", "user2", "user3"}

		// Mock 期望
		mockRepo.EXPECT().
			BatchDelete(ctx, ids).
			Return(int32(3), int32(0), []string{}, nil)

		// 执行测试
		result, err := uc.BatchDeleteUsers(ctx, ids)

		// 断言
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(3), result.SuccessCount)
		assert.Equal(t, int32(0), result.FailedCount)
		assert.Empty(t, result.FailedIDs)
	})

	t.Run("ID列表为空", func(t *testing.T) {
		ids := []string{}

		// 执行测试
		result, err := uc.BatchDeleteUsers(ctx, ids)

		// 断言
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "ID列表不能为空")
	})
}

func TestUserUsecase_CheckAccountExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSystemUserRepo(ctrl)
	logger := log.DefaultLogger
	uc := systemuser.NewUserUsecase(mockRepo, logger)

	ctx := context.Background()

	t.Run("用户名存在", func(t *testing.T) {
		account := "testuser"
		existingUser := &systemuser.SystemUser{
			ID:      ptr.Of("user123"),
			Account: ptr.Of(account),
		}

		// Mock 期望
		mockRepo.EXPECT().
			FindByUsername(ctx, account).
			Return(existingUser, nil)

		// 执行测试
		exists, err := uc.CheckAccountExists(ctx, account)

		// 断言
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("用户名不存在", func(t *testing.T) {
		account := "nonexistent"

		// Mock 期望
		mockRepo.EXPECT().
			FindByUsername(ctx, account).
			Return(nil, nil)

		// 执行测试
		exists, err := uc.CheckAccountExists(ctx, account)

		// 断言
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("无效的用户名格式", func(t *testing.T) {
		account := "ab"

		// 执行测试
		exists, err := uc.CheckAccountExists(ctx, account)

		// 断言
		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), "用户名长度必须在3-50个字符之间")
	})
}

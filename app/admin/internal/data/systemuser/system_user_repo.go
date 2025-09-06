package systemuser

import (
	"context"
	"fmt"
	"qn-base/app/admin/internal/data/idgen"
	"time"

	bizsystemuser "qn-base/app/admin/internal/biz/systemuser"
	"qn-base/app/admin/internal/data/data"
	"qn-base/app/admin/internal/data/ent"
	"qn-base/app/admin/internal/data/ent/systemuser"

	"github.com/go-kratos/kratos/v2/log"
)

type systemUserRepo struct {
	data  *data.Data
	log   *log.Helper
	idGen *idgen.IDGenerator
}

// convertToBizUser 将数据库实体转换为业务对象
func (s systemUserRepo) convertToBizUser(entity *ent.SystemUser) *bizsystemuser.SystemUser {
	if entity == nil {
		return nil
	}

	return &bizsystemuser.SystemUser{
		ID:        &entity.ID,
		CreateBy:  entity.CreateBy,
		CreatedAt: entity.CreatedAt,
		UpdateBy:  entity.UpdateBy,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
		TenantID:  &entity.TenantID,
		Account:   &entity.Account,
		Password:  entity.Password,
		Nickname:  entity.Nickname,
		Remark:    entity.Remark,
		DeptID:    entity.DeptID,
		PostIds:   entity.PostIds,
		Email:     entity.Email,
		Mobile:    entity.Mobile,
		Sex:       entity.Sex,
		Avatar:    entity.Avatar,
		Status:    &entity.Status,
		LoginIP:   entity.LoginIP,
		LoginDate: entity.LoginDate,
	}
}

// NewSystemUserRepo .
func NewSystemUserRepo(data *data.Data, idGen *idgen.IDGenerator, logger log.Logger) bizsystemuser.SystemUserRepo {
	return &systemUserRepo{
		data:  data,
		log:   log.NewHelper(log.With(logger, "module", "systemuser/repo")),
		idGen: idGen,
	}
}

func (s systemUserRepo) Save(ctx context.Context, user *bizsystemuser.SystemUser) (*bizsystemuser.SystemUser, error) {
	// 创建系统用户
	create := s.data.DB.SystemUser(ctx).Create().
		SetAccount(*user.Account).
		SetNillablePassword(user.Password).
		SetNillableNickname(user.Nickname).
		SetNillableRemark(user.Remark).
		SetNillableDeptID(user.DeptID).
		SetNillablePostIds(user.PostIds).
		SetNillableEmail(user.Email).
		SetNillableMobile(user.Mobile).
		SetNillableSex(user.Sex).
		SetNillableAvatar(user.Avatar).
		SetNillableStatus(user.Status).
		SetNillableLoginIP(user.LoginIP).
		SetNillableLoginDate(user.LoginDate)

	// 设置租户ID
	if user.TenantID != nil {
		create.SetTenantID(*user.TenantID)
	}

	// 设置创建人
	if user.CreateBy != nil {
		create.SetCreateBy(*user.CreateBy)
	}

	result, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为业务对象
	return s.convertToBizUser(result), nil
}

func (s systemUserRepo) Update(ctx context.Context, user *bizsystemuser.SystemUser) (*bizsystemuser.SystemUser, error) {
	// 检查ID是否为空
	if user.ID == nil {
		return nil, fmt.Errorf("user ID cannot be nil")
	}

	// 更新系统用户
	update := s.data.DB.SystemUser(ctx).UpdateOneID(*user.ID).
		SetNillablePassword(user.Password).
		SetNillableNickname(user.Nickname).
		SetNillableRemark(user.Remark).
		SetNillableDeptID(user.DeptID).
		SetNillablePostIds(user.PostIds).
		SetNillableEmail(user.Email).
		SetNillableMobile(user.Mobile).
		SetNillableSex(user.Sex).
		SetNillableAvatar(user.Avatar).
		SetNillableStatus(user.Status).
		SetNillableLoginIP(user.LoginIP).
		SetNillableLoginDate(user.LoginDate)

	// 设置更新人
	if user.UpdateBy != nil {
		update.SetUpdateBy(*user.UpdateBy)
	}

	result, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为业务对象
	return s.convertToBizUser(result), nil
}

func (s systemUserRepo) Delete(ctx context.Context, id string) error {
	// 删除系统用户（软删除）
	// 使用 Update 方法加条件，确保只删除未删除的记录
	affected, err := s.data.DB.SystemUser(ctx).Update().
		Where(systemuser.ID(id)).
		Where(systemuser.DeletedAtIsNil()). // 只删除未删除的记录
		SetDeletedAt(time.Now()).
		Save(ctx)

	if err != nil {
		return err
	}

	// 检查是否有记录被影响
	if affected == 0 {
		return &ent.NotFoundError{}
	}

	return nil
}

func (s systemUserRepo) FindByID(ctx context.Context, id string) (*bizsystemuser.SystemUser, error) {
	// 根据ID查找系统用户
	result, err := s.data.DB.SystemUser(ctx).Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	// 转换为业务对象
	return s.convertToBizUser(result), nil
}

func (s systemUserRepo) FindByUsername(ctx context.Context, username string) (*bizsystemuser.SystemUser, error) {
	// 根据用户名查找系统用户
	result, err := s.data.DB.SystemUser(ctx).Query().
		Where(systemuser.Account(username)).
		Where(systemuser.DeletedAtIsNil()). // 只查找未删除的用户
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	// 转换为业务对象
	return s.convertToBizUser(result), nil
}

func (s systemUserRepo) ListSystemUsers(ctx context.Context, request *bizsystemuser.ListUserRequest) ([]*bizsystemuser.SystemUser, int32, error) {
	// 查询系统用户列表
	query := s.data.DB.SystemUser(ctx).Query().
		Where(systemuser.DeletedAtIsNil()) // 只查询未删除的用户

	// 根据用户名过滤
	if request.Username != "" {
		query = query.Where(systemuser.AccountContains(request.Username))
	}

	// 获取总数
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	query = query.
		Offset(int((request.Page - 1) * request.PageSize)).
		Limit(int(request.PageSize)).
		Order(ent.Desc(systemuser.FieldCreatedAt))

	results, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 转换为业务对象列表
	users := make([]*bizsystemuser.SystemUser, len(results))
	for i, result := range results {
		users[i] = s.convertToBizUser(result)
	}

	return users, int32(count), nil
}

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

	// 根据用户名过滤（模糊匹配）
	if request.Username != "" {
		query = query.Where(systemuser.AccountContains(request.Username))
	}

	// 根据邮箱过滤（模糊匹配）
	if request.Email != "" {
		query = query.Where(systemuser.EmailContains(request.Email))
	}

	// 根据手机号过滤（模糊匹配）
	if request.Mobile != "" {
		query = query.Where(systemuser.MobileContains(request.Mobile))
	}

	// 根据状态过滤
	if request.Status != nil {
		query = query.Where(systemuser.Status(*request.Status))
	}

	// 根据部门ID过滤
	if request.DeptID != "" {
		query = query.Where(systemuser.DeptID(request.DeptID))
	}

	// 根据租户ID过滤
	if request.TenantID != "" {
		query = query.Where(systemuser.TenantID(request.TenantID))
	}

	// 根据创建时间范围过滤
	if request.StartDate != "" {
		startTime, err := time.Parse("2006-01-02", request.StartDate)
		if err == nil {
			query = query.Where(systemuser.CreatedAtGTE(startTime))
		}
	}
	if request.EndDate != "" {
		endTime, err := time.Parse("2006-01-02", request.EndDate)
		if err == nil {
			// 结束时间设置为当天的23:59:59
			endTime = endTime.Add(24*time.Hour - time.Second)
			query = query.Where(systemuser.CreatedAtLTE(endTime))
		}
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

// BatchDelete implements batch delete users.
func (s systemUserRepo) BatchDelete(ctx context.Context, ids []string) (int32, int32, []string, error) {
	if len(ids) == 0 {
		return 0, 0, nil, fmt.Errorf("ids cannot be empty")
	}

	var successCount, failedCount int32
	var failedIDs []string

	for _, id := range ids {
		// 软删除单个用户
		affected, err := s.data.DB.SystemUser(ctx).Update().
			Where(systemuser.ID(id)).
			Where(systemuser.DeletedAtIsNil()).
			SetDeletedAt(time.Now()).
			Save(ctx)

		if err != nil || affected == 0 {
			failedCount++
			failedIDs = append(failedIDs, id)
		} else {
			successCount++
		}
	}

	return successCount, failedCount, failedIDs, nil
}

// FindByEmail finds user by email.
func (s systemUserRepo) FindByEmail(ctx context.Context, email string) (*bizsystemuser.SystemUser, error) {
	result, err := s.data.DB.SystemUser(ctx).Query().
		Where(systemuser.Email(email)).
		Where(systemuser.DeletedAtIsNil()).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return s.convertToBizUser(result), nil
}

// FindByMobile finds user by mobile.
func (s systemUserRepo) FindByMobile(ctx context.Context, mobile string) (*bizsystemuser.SystemUser, error) {
	result, err := s.data.DB.SystemUser(ctx).Query().
		Where(systemuser.Mobile(mobile)).
		Where(systemuser.DeletedAtIsNil()).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return s.convertToBizUser(result), nil
}

// ChangeStatus implements change user status.
func (s systemUserRepo) ChangeStatus(ctx context.Context, id string, status int8) error {
	affected, err := s.data.DB.SystemUser(ctx).Update().
		Where(systemuser.ID(id)).
		Where(systemuser.DeletedAtIsNil()).
		SetStatus(status).
		Save(ctx)

	if err != nil {
		return err
	}

	if affected == 0 {
		return &ent.NotFoundError{}
	}

	return nil
}

// GetUserStats implements get user statistics.
func (s systemUserRepo) GetUserStats(ctx context.Context, tenantID string) (*bizsystemuser.UserStats, error) {
	query := s.data.DB.SystemUser(ctx).Query().
		Where(systemuser.DeletedAtIsNil())

	// 如果提供了租户ID，则按租户过滤
	if tenantID != "" {
		query = query.Where(systemuser.TenantID(tenantID))
	}

	// 总用户数
	totalUsers, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	// 活跃用户数（状态为1）
	activeUsers, err := query.Where(systemuser.Status(1)).Count(ctx)
	if err != nil {
		return nil, err
	}

	// 停用用户数（状态为0）
	inactiveUsers, err := query.Where(systemuser.Status(0)).Count(ctx)
	if err != nil {
		return nil, err
	}

	// 今天注册的用户数
	todayStart := time.Now().Truncate(24 * time.Hour)
	todayEnd := todayStart.Add(24 * time.Hour)
	todayRegistered, err := query.
		Where(systemuser.CreatedAtGTE(todayStart)).
		Where(systemuser.CreatedAtLT(todayEnd)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// 本周注册的用户数
	weekStart := todayStart.AddDate(0, 0, -int(todayStart.Weekday()))
	weekEnd := weekStart.AddDate(0, 0, 7)
	thisWeekRegistered, err := query.
		Where(systemuser.CreatedAtGTE(weekStart)).
		Where(systemuser.CreatedAtLT(weekEnd)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// 本月注册的用户数
	monthStart := time.Date(todayStart.Year(), todayStart.Month(), 1, 0, 0, 0, 0, todayStart.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	thisMonthRegistered, err := query.
		Where(systemuser.CreatedAtGTE(monthStart)).
		Where(systemuser.CreatedAtLT(monthEnd)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	return &bizsystemuser.UserStats{
		TotalUsers:          int32(totalUsers),
		ActiveUsers:         int32(activeUsers),
		InactiveUsers:       int32(inactiveUsers),
		TodayRegistered:     int32(todayRegistered),
		ThisWeekRegistered:  int32(thisWeekRegistered),
		ThisMonthRegistered: int32(thisMonthRegistered),
	}, nil
}

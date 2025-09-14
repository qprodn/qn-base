package convertor

import (
	v1 "qn-base/api/gen/go/admin/v1"
	"qn-base/app/admin/internal/biz/systemuser"
	"qn-base/pkg/lang/conv"
	"qn-base/pkg/lang/ptr"
	"time"
)

// ToCreateUserBiz converts CreateUserRequest to SystemUser (biz).
func ToCreateUserBiz(req *v1.CreateUserRequest) (*systemuser.SystemUser, error) {
	if req == nil {
		return nil, nil
	}

	// 转换时间字段
	var loginDate *time.Time
	if req.LoginDate != nil {
		toTime, err := conv.DefaultStrToTime(*req.LoginDate)
		if err != nil {
			return nil, err
		}
		loginDate = ptr.Of(toTime)
	}

	return &systemuser.SystemUser{
		Account:   &req.Account,
		Password:  &req.Password,
		Nickname:  req.Nickname,
		Remark:    req.Remark,
		DeptID:    req.DeptId,
		PostIds:   req.PostIds,
		Email:     req.Email,
		Mobile:    req.Mobile,
		Sex:       ptr.Of(int8(ptr.From(req.Sex))),
		Avatar:    req.Avatar,
		LoginIP:   req.LoginIp,
		LoginDate: loginDate,
		TenantID:  req.TenantId,
	}, nil
}

// ToUpdateUserBiz converts UpdateUserRequest to SystemUser (biz).
func ToUpdateUserBiz(req *v1.UpdateUserRequest) (*systemuser.SystemUser, error) {
	if req == nil {
		return nil, nil
	}

	// 转换时间字段
	var loginDate *time.Time
	if req.LoginDate != nil {
		toTime, err := conv.DefaultStrToTime(*req.LoginDate)
		if err != nil {
			return nil, err
		}
		loginDate = ptr.Of(toTime)
	}

	return &systemuser.SystemUser{
		ID:        &req.Id,
		Nickname:  req.Nickname,
		Remark:    req.Remark,
		DeptID:    req.DeptId,
		PostIds:   req.PostIds,
		Email:     req.Email,
		Mobile:    req.Mobile,
		Sex:       ptr.Of(int8(ptr.From(req.Sex))),
		Avatar:    req.Avatar,
		Status:    ptr.Of(int8(ptr.From(req.Status))),
		LoginIP:   req.LoginIp,
		LoginDate: loginDate,
		TenantID:  req.TenantId,
	}, nil
}

// ToListUserRequestBiz converts ListUsersRequest to ListUserRequest (biz).
func ToListUserRequestBiz(req *v1.ListUsersRequest) *systemuser.ListUserRequest {
	if req == nil {
		return nil
	}

	bizReq := &systemuser.ListUserRequest{
		Page:      ptr.From(req.Page),
		PageSize:  ptr.From(req.PageSize),
		Username:  ptr.From(req.Account),
		Email:     ptr.From(req.Email),
		Mobile:    ptr.From(req.Mobile),
		DeptID:    ptr.From(req.DeptId),
		StartDate: ptr.From(req.StartDate),
		EndDate:   ptr.From(req.EndDate),
	}

	if req.Status != nil {
		status := int8(*req.Status)
		bizReq.Status = &status
	}

	return bizReq
}

// ToUserInfo converts SystemUser (biz) to UserInfo (proto).
func ToUserInfo(user *systemuser.SystemUser) *v1.UserInfo {
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

// ToUserInfos converts a slice of SystemUser (biz) to UserInfo (proto).
func ToUserInfos(users []*systemuser.SystemUser) []*v1.UserInfo {
	if users == nil {
		return nil
	}

	userInfos := make([]*v1.UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = ToUserInfo(user)
	}
	return userInfos
}

// ToUserStats converts UserStats (biz) to UserStats (proto).
func ToUserStats(stats *systemuser.UserStats) *v1.UserStats {
	if stats == nil {
		return nil
	}

	return &v1.UserStats{
		TotalUsers:          stats.TotalUsers,
		ActiveUsers:         stats.ActiveUsers,
		InactiveUsers:       stats.InactiveUsers,
		TodayRegistered:     stats.TodayRegistered,
		ThisWeekRegistered:  stats.ThisWeekRegistered,
		ThisMonthRegistered: stats.ThisMonthRegistered,
	}
}

// ToBatchDeleteResult converts BatchDeleteResult (biz) to BatchDeleteUsersReply (proto).
func ToBatchDeleteResult(result *systemuser.BatchDeleteResult) *v1.BatchDeleteUsersReply {
	if result == nil {
		return nil
	}

	return &v1.BatchDeleteUsersReply{
		SuccessCount: result.SuccessCount,
		FailedCount:  result.FailedCount,
		FailedIds:    result.FailedIDs,
	}
}
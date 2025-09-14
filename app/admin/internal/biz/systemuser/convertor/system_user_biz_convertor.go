package convertor


import (
	"qn-base/app/admin/internal/biz/systemuser"
)

// ToDataUser converts SystemUser (biz) to data layer model.
// This function would be used when converting from biz layer to data layer.
func ToDataUser(bizUser *systemuser.SystemUser) *systemuser.SystemUser {
	if bizUser == nil {
		return nil
	}
	
	// 由于目前biz层和data层使用相同的结构体，直接返回
	// 如果以后data层有独立的结构体，在这里进行转换
	return bizUser
}

// FromDataUser converts data layer model to SystemUser (biz).
// This function would be used when converting from data layer to biz layer.
func FromDataUser(dataUser *systemuser.SystemUser) *systemuser.SystemUser {
	if dataUser == nil {
		return nil
	}
	
	// 由于目前biz层和data层使用相同的结构体，直接返回
	// 如果以后data层有独立的结构体，在这里进行转换
	return dataUser
}

// ToDataUsers converts a slice of SystemUser (biz) to data layer models.
func ToDataUsers(bizUsers []*systemuser.SystemUser) []*systemuser.SystemUser {
	if bizUsers == nil {
		return nil
	}
	
	dataUsers := make([]*systemuser.SystemUser, len(bizUsers))
	for i, user := range bizUsers {
		dataUsers[i] = ToDataUser(user)
	}
	return dataUsers
}

// FromDataUsers converts a slice of data layer models to SystemUser (biz).
func FromDataUsers(dataUsers []*systemuser.SystemUser) []*systemuser.SystemUser {
	if dataUsers == nil {
		return nil
	}
	
	bizUsers := make([]*systemuser.SystemUser, len(dataUsers))
	for i, user := range dataUsers {
		bizUsers[i] = FromDataUser(user)
	}
	return bizUsers
}
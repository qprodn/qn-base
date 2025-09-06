package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// 邮箱正则表达式
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// 手机号正则表达式（中国大陆）
var mobileRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// 用户名正则表达式
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	if email == "" {
		return nil // 允许为空
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("无效的邮箱格式")
	}
	return nil
}

// ValidateMobile 验证手机号格式
func ValidateMobile(mobile string) error {
	if mobile == "" {
		return nil // 允许为空
	}
	if !mobileRegex.MatchString(mobile) {
		return fmt.Errorf("无效的手机号格式")
	}
	return nil
}

// ValidateUsername 验证用户名格式
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if len(username) < 3 || len(username) > 50 {
		return fmt.Errorf("用户名长度必须在3-50个字符之间")
	}
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("用户名只能包含字母、数字和下划线")
	}
	return nil
}

// ValidatePassword 验证密码格式
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("密码不能为空")
	}
	if len(password) < 6 || len(password) > 128 {
		return fmt.Errorf("密码长度必须在6-128个字符之间")
	}
	return nil
}

// ValidateNickname 验证昵称格式
func ValidateNickname(nickname string) error {
	if nickname == "" {
		return nil // 允许为空
	}
	if len(nickname) > 50 {
		return fmt.Errorf("昵称长度不能超过50个字符")
	}
	return nil
}

// ValidateRemark 验证备注格式
func ValidateRemark(remark string) error {
	if remark == "" {
		return nil // 允许为空
	}
	if len(remark) > 255 {
		return fmt.Errorf("备注长度不能超过255个字符")
	}
	return nil
}

// ValidateStatus 验证状态值
func ValidateStatus(status int8) error {
	if status != 0 && status != 1 {
		return fmt.Errorf("状态值只能是0（停用）或1（正常）")
	}
	return nil
}

// ValidateSex 验证性别值
func ValidateSex(sex int8) error {
	if sex != 0 && sex != 1 {
		return fmt.Errorf("性别值只能是0（女）或1（男）")
	}
	return nil
}

// ValidateStringLength 验证字符串长度
func ValidateStringLength(str string, fieldName string, minLen, maxLen int) error {
	if str == "" {
		return nil // 允许为空
	}
	if len(str) < minLen || len(str) > maxLen {
		return fmt.Errorf("%s长度必须在%d-%d个字符之间", fieldName, minLen, maxLen)
	}
	return nil
}

// ValidateRequiredString 验证必填字符串
func ValidateRequiredString(str string, fieldName string) error {
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("%s不能为空", fieldName)
	}
	return nil
}

// ValidateIDs 验证ID列表
func ValidateIDs(ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("ID列表不能为空")
	}
	if len(ids) > 100 {
		return fmt.Errorf("批量操作ID数量不能超过100个")
	}
	for _, id := range ids {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("ID不能为空")
		}
	}
	return nil
}

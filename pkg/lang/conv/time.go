package conv

import "time"

const (
	DefaultTimeLayout = "2006-01-02 15:04:05"
)

// DefaultStrToTime 默认模板转换
func DefaultStrToTime(v string) (time.Time, error) {
	return StrToTimeCustom(v, DefaultTimeLayout)
}

// StrToTimeCustom 自定义模板转换
func StrToTimeCustom(v string, layout string) (time.Time, error) {
	return time.ParseInLocation(layout, v, time.Local)
}

// TimeToStrCustom 自定义模板转换
func TimeToStrCustom(t time.Time, layout string) string {
	return t.Format(layout)
}

// TimeToDefaultStr 默认模板转换
func TimeToDefaultStr(t time.Time) string {
	return TimeToStrCustom(t, DefaultTimeLayout)
}

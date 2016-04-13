package utils

import (
	"time"
	"fmt"
)
//计算指定时间到现在的时长
func GetDuration(t time.Time)string{
	duration :=time.Now().Sub(t)
	//return  duration.String()
	return DurationToString(duration)
}

//时长转字符串
func DurationToString(d time.Duration) string {
	u := uint64(d)

	switch  {
	case u < uint64(time.Second):
		return "刚刚"
	case u < uint64(time.Minute):
		return  fmt.Sprintf("%d秒前",int(d.Seconds()))
	case u < uint64(time.Hour):
		return  fmt.Sprintf("%d分钟前",int(d.Minutes()))
	case u < uint64(time.Hour*24):
		return  fmt.Sprintf("%d小时前",int(d.Hours()))
	default:
		return  fmt.Sprintf("%d天前",int(d.Hours()/24))
	}
	return ""
}



package request

import (
	"github.com/KeSilent/study-hub/server/model/common/request"
	"github.com/KeSilent/study-hub/server/model/edu_user_course"
	"time"
)

type EduEnrollmentSearch struct {
	edu_user_course.EduEnrollment
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	UserName       string     `json:"userName" form:"userName"`       // 学员姓名搜索
	UserPhone      string     `json:"userPhone" form:"userPhone"`     // 手机号搜索
	request.PageInfo
}

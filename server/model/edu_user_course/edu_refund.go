package edu_user_course

import (
	"time"

	"github.com/KeSilent/study-hub/server/global"
)

// EduRefund 退费/转课记录
type EduRefund struct {
	global.GVA_MODEL
	EnrollmentId   *int      `json:"enrollmentId" form:"enrollmentId" gorm:"column:enrollment_id;comment:报名ID;size:10;"`
	UserId         *int      `json:"userId" form:"userId" gorm:"column:user_id;comment:学员ID;size:10;"`
	CourseId       *int      `json:"courseId" form:"courseId" gorm:"column:course_id;comment:课程ID;size:10;"`
	Action         string    `json:"action" form:"action" gorm:"column:action;comment:操作类型(refund/transfer);size:20;"`
	RefundSessions int       `json:"refundSessions" form:"refundSessions" gorm:"column:refund_sessions;comment:退/转课时数;"`
	RefundAmount   float64   `json:"refundAmount" form:"refundAmount" gorm:"column:refund_amount;type:decimal(10,2);comment:退费金额;default:0;"`
	Reason         string    `json:"reason" form:"reason" gorm:"column:reason;comment:原因;size:255;"`
	OperatorId     *int      `json:"operatorId" form:"operatorId" gorm:"column:operator_id;comment:操作人ID;size:10;"`
	OperatorName   string    `json:"operatorName" form:"operatorName" gorm:"column:operator_name;comment:操作人姓名;size:255;"`
	RefundTime     time.Time `json:"refundTime" form:"refundTime" gorm:"column:refund_time;comment:退费时间;"`
}

// TableName EduRefund 表名
func (EduRefund) TableName() string {
	return "edu_refund"
}

/*
 * @Author: Yang
 * @Date: 2023-04-10 18:35:56
 * @Description: 请填写简介
 */
// 自动生成模板EduEnrollment
package edu_user_course

import (
	"github.com/KeSilent/study-hub/server/global"
	"github.com/KeSilent/study-hub/server/model/edu_organization"
)

// EduEnrollment 结构体
type EduEnrollment struct {
	global.GVA_MODEL
	UserId            *int                       `json:"userId" form:"userId" gorm:"column:user_id;comment:用户ID;size:10;"`
	CourseId          *int                       `json:"courseId" form:"courseId" gorm:"column:course_id;comment:课程ID;size:10;"`
	TotalSessions     *int                       `json:"totalSessions" form:"totalSessions" gorm:"column:total_sessions;comment:总课时数;size:10;"`
	RemainingSessions *int                       `json:"remainingSessions" form:"remainingSessions" gorm:"column:remaining_sessions;comment:剩余课时数;size:10;"`
	PaidSessions      *int                       `json:"paidSessions" form:"paidSessions" gorm:"column:paid_sessions;comment:付费课时数;size:10;"`
	GiftSessions      *int                       `json:"giftSessions" form:"giftSessions" gorm:"column:gift_sessions;comment:赠送课时数;size:10;"`
	RemainingPaid     *int                       `json:"remainingPaidSessions" form:"remainingPaidSessions" gorm:"column:remaining_paid_sessions;comment:剩余付费课时;size:10;"`
	RemainingGift     *int                       `json:"remainingGiftSessions" form:"remainingGiftSessions" gorm:"column:remaining_gift_sessions;comment:剩余赠送课时;size:10;"`
	PricePerSession   float64                    `json:"pricePerSession" form:"pricePerSession" gorm:"column:price_per_session;type:decimal(10,2);comment:课时单价;default:0;"`
	DiscountAmount    float64                    `json:"discountAmount" form:"discountAmount" gorm:"column:discount_amount;type:decimal(10,2);comment:优惠金额;default:0;"`
	TotalAmount       float64                    `json:"totalAmount" form:"totalAmount" gorm:"column:total_amount;type:decimal(10,2);comment:应收总额;default:0;"`
	PaidAmount        float64                    `json:"paidAmount" form:"paidAmount" gorm:"column:paid_amount;type:decimal(10,2);comment:已收金额;default:0;"`
	BalanceAmount     float64                    `json:"balanceAmount" form:"balanceAmount" gorm:"column:balance_amount;type:decimal(10,2);comment:应收余额;default:0;"`
	EduCourse         edu_organization.EduCourse `json:"eduCourse" gorm:"foreignKey:CourseId"`
	UserName          string                     `json:"userName" gorm:"-"`  // 用户姓名（从 sys_users 表联查）
	UserPhone         string                     `json:"userPhone" gorm:"-"` // 用户手机号（从 sys_users 表联查）
}

// TableName EduEnrollment 表名
func (EduEnrollment) TableName() string {
	return "edu_enrollment"
}

// ConsumptionClassReq 消耗课时请求参数
type ConsumptionClassResp struct {
	UserId            int    `json:"userId"`            // 用户ID
	CourseId          int    `json:"courseId"`          // 课程ID
	SessionsToConsume int    `json:"sessionsToConsume"` // 消耗课时数
	Reason            string `json:"reason"`            // 消耗原因
	UseDate           string `json:"useDate"`           // 上课时间（YYYY-MM-DD 或 YYYY-MM-DD HH:mm）
	TeacherId         int    `json:"teacherId"`         // 教师ID（可选）
	TeacherName       string `json:"teacherName"`       // 教师姓名（可选）
	Chargeable        *bool  `json:"chargeable"`        // 是否计费（可选）
}

// AddSessionReq 增加课时请求参数
type AddSessionResp struct {
	UserId        int    `json:"userId"`        // 用户ID
	CourseId      int    `json:"courseId"`      // 课程ID
	SessionsToAdd int    `json:"sessionsToAdd"` // 增加课时数
	Reason        string `json:"reason"`        // 消耗原因
	UseDate       string `json:"useDate"`       // 日期（YYYY-MM-DD 或 YYYY-MM-DD HH:mm）
	TeacherId     int    `json:"teacherId"`     // 教师ID（可选）
	TeacherName   string `json:"teacherName"`   // 教师姓名（可选）
}

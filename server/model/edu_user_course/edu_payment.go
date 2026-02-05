package edu_user_course

import (
	"time"

	"github.com/KeSilent/study-hub/server/global"
)

// EduPayment 结构体
type EduPayment struct {
	global.GVA_MODEL
	EnrollmentId *int      `json:"enrollmentId" form:"enrollmentId" gorm:"column:enrollment_id;comment:报名ID;size:10;"`
	Amount       float64   `json:"amount" form:"amount" gorm:"column:amount;type:decimal(10,2);comment:收款金额;default:0;"`
	PayTime      time.Time `json:"payTime" form:"payTime" gorm:"column:pay_time;comment:收款时间;"`
	OperatorId   *int      `json:"operatorId" form:"operatorId" gorm:"column:operator_id;comment:操作人ID;size:10;"`
	OperatorName string    `json:"operatorName" form:"operatorName" gorm:"column:operator_name;comment:操作人姓名;size:255;"`
	Remark       string    `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:255;"`
}

// TableName EduPayment 表名
func (EduPayment) TableName() string {
	return "edu_payment"
}

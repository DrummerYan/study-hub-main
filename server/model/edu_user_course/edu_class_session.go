// 自动生成模板EduClassSession
package edu_user_course

import (
	"time"

	"github.com/KeSilent/study-hub/server/global"
)

// EduClassSession 结构体
type EduClassSession struct {
	global.GVA_MODEL
	EnrollmentId  *int          `json:"enrollmentId" form:"enrollmentId" gorm:"column:enrollment_id;comment:报名ID;size:10;"`
	Action        string        `json:"action" form:"action" gorm:"column:action;type:enum('add', 'subtract');comment:操作类型（增加或扣除）;"`
	Reason        string        `json:"reason" form:"reason" gorm:"column:reason;comment:操作原因;size:255;"`
	NumSessions   *int          `json:"numSessions" form:"numSessions" gorm:"column:num_sessions;comment:课时数量;size:10;"`
	CourseName    string        `json:"courseName" form:"courseName" gorm:"column:course_name;comment:课程名称;size:255;"`
	UserName      string        `json:"userName" form:"userName" gorm:"column:user_name;comment:学员姓名;size:255;"`
	TeacherId     *int          `json:"teacherId" form:"teacherId" gorm:"column:teacher_id;comment:教师ID;size:10;"`
	TeacherName   string        `json:"teacherName" form:"teacherName" gorm:"column:teacher_name;comment:教师姓名;size:255;"`
	UnitPrice     float64       `json:"unitPrice" form:"unitPrice" gorm:"column:unit_price;type:decimal(10,2);comment:课时单价;default:0;"`
	Amount        float64       `json:"amount" form:"amount" gorm:"column:amount;type:decimal(10,2);comment:本次金额;default:0;"`
	Chargeable    bool          `json:"chargeable" form:"chargeable" gorm:"column:chargeable;comment:是否计费;default:true;"`
	EduEnrollment EduEnrollment `json:"eduEnrollment" gorm:"foreignKey:EnrollmentId"` // 用户科目
	UseDate       time.Time     `json:"useDate" form:"useDate" gorm:"column:use_date;comment:使用日期;size:255;"`
}

// TableName EduClassSession 表名
func (EduClassSession) TableName() string {
	return "edu_class_session"
}

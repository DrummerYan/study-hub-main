/*
 * @Author: Yang
 * @Date: 2023-04-10 18:35:56
 * @Description: 学生绑定课程信息
 */
package edu_user_course

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/KeSilent/study-hub/server/global"
	"github.com/KeSilent/study-hub/server/model/common/request"
	"github.com/KeSilent/study-hub/server/model/edu_user_course"
	edu_user_courseReq "github.com/KeSilent/study-hub/server/model/edu_user_course/request"
)

type EduEnrollmentService struct {
}

var eduClassSessionService EduClassSessionService

// CreateEduEnrollment 创建EduEnrollment记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduEnrollmentService *EduEnrollmentService) CreateEduEnrollment(eduEnrollment *edu_user_course.EduEnrollment) (err error) {
	applyEnrollmentDefaults(eduEnrollment)
	applyEnrollmentFinance(eduEnrollment)
	err = global.GVA_DB.Create(eduEnrollment).Error
	return err
}

// DeleteEduEnrollment 删除EduEnrollment记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduEnrollmentService *EduEnrollmentService) DeleteEduEnrollment(eduEnrollment edu_user_course.EduEnrollment) (err error) {
	err = global.GVA_DB.Delete(&eduEnrollment).Error
	return err
}

// DeleteEduEnrollmentByIds 批量删除EduEnrollment记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduEnrollmentService *EduEnrollmentService) DeleteEduEnrollmentByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]edu_user_course.EduEnrollment{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateEduEnrollment 更新EduEnrollment记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduEnrollmentService *EduEnrollmentService) UpdateEduEnrollment(eduEnrollment edu_user_course.EduEnrollment) (err error) {
	applyEnrollmentDefaults(&eduEnrollment)
	applyEnrollmentFinance(&eduEnrollment)
	err = global.GVA_DB.Save(&eduEnrollment).Error
	return err
}

// UpdateEduEnrollment 更新课程选择
// Author [piexlmax](https://github.com/piexlmax)
func (eduEnrollmentService *EduEnrollmentService) UpdateCourseId(eduEnrollment edu_user_course.EduEnrollment) (err error) {
	err = global.GVA_DB.Model(&edu_user_course.EduEnrollment{}).Where("ID = ?", eduEnrollment.ID).Update("course_id", eduEnrollment.CourseId).Error

	return err
}

// GetEduEnrollment 根据id获取EduEnrollment记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduEnrollmentService *EduEnrollmentService) GetEduEnrollment(id uint) (eduEnrollment edu_user_course.EduEnrollment, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&eduEnrollment).Error
	return
}

// GetEduEnrollmentInfoList 分页获取EduEnrollment记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduEnrollmentService *EduEnrollmentService) GetEduEnrollmentInfoList(info edu_user_courseReq.EduEnrollmentSearch) (list []edu_user_course.EduEnrollment, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&edu_user_course.EduEnrollment{})
	var eduEnrollments []edu_user_course.EduEnrollment
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	// 按课程ID搜索
	if info.CourseId != nil && *info.CourseId > 0 {
		db = db.Where("course_id = ?", *info.CourseId)
	}

	// 按学员姓名或手机号搜索（需要关联查询用户表）
	if info.UserName != "" || info.UserPhone != "" {
		db = db.Joins("JOIN sys_users ON sys_users.id = edu_enrollment.user_id")
		if info.UserName != "" {
			db = db.Where("sys_users.nick_name LIKE ?", "%"+info.UserName+"%")
		}
		if info.UserPhone != "" {
			db = db.Where("sys_users.phone LIKE ?", "%"+info.UserPhone+"%")
		}
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 预加载课程信息,便于前端显示
	err = db.Preload("EduCourse").Limit(limit).Offset(offset).Find(&eduEnrollments).Error
	if err != nil {
		return
	}

	// 联查用户信息，填充用户姓名和手机号
	for i := range eduEnrollments {
		if eduEnrollments[i].UserId != nil {
			var user struct {
				NickName string
				Phone    string
			}
			err := global.GVA_DB.Table("sys_users").
				Select("nick_name, phone").
				Where("id = ?", *eduEnrollments[i].UserId).
				First(&user).Error
			if err == nil {
				eduEnrollments[i].UserName = user.NickName
				eduEnrollments[i].UserPhone = user.Phone
			}
		}
	}

	return eduEnrollments, total, err
}

// ConsumeSession 消耗课时
func (eduEnrollmentService *EduEnrollmentService) ConsumeSession(userID, courseID, sessionsToConsume int, reason string, useData string, teacherId int, teacherName string, chargeable bool) error {
	var enrollment edu_user_course.EduEnrollment
	db := global.GVA_DB

	// 查询用户在指定课程中的报名信息
	result := db.Where("user_id = ? AND course_id = ?", userID, courseID).Preload("EduCourse").First(&enrollment)

	// 如果找不到报名信息或发生错误，则返回错误
	if result.Error != nil {
		return errors.New("报名信息未找到或查询错误")
	}

	// 检查剩余课时是否足够
	if enrollment.RemainingSessions == nil || *enrollment.RemainingSessions < sessionsToConsume {
		return errors.New("剩余课时不足")
	}

	// 扣除课时
	*enrollment.RemainingSessions -= sessionsToConsume

	enenrollmentId := int(enrollment.ID)
	t, err := parseUseDate(useData)
	if err != nil {
		return err
	}

	// 获取学员姓名
	var user struct {
		NickName string
	}
	err = db.Table("sys_users").Select("nick_name").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return errors.New("查询学员信息失败")
	}

	// 获取教师姓名（可选）
	if teacherName == "" && teacherId > 0 {
		var teacher struct {
			NickName string
		}
		if err := db.Table("sys_users").Select("nick_name").Where("id = ?", teacherId).First(&teacher).Error; err == nil {
			teacherName = teacher.NickName
		}
	}

	// 更新数据库中的报名信息
	result = db.Save(&enrollment)

	if result.Error != nil {
		return errors.New("更新报名信息失败")
	}
	// 记录课时操作
	unitPrice := enrollment.PricePerSession
	amount := 0.0
	if chargeable {
		amount = float64(sessionsToConsume) * unitPrice
	}
	classSession := edu_user_course.EduClassSession{
		EnrollmentId: &enenrollmentId,
		Action:       "subtract",
		Reason:       reason,
		NumSessions:  &sessionsToConsume,
		CourseName:   enrollment.EduCourse.CourseName,
		UserName:     user.NickName,
		TeacherId:    intPointer(teacherId),
		TeacherName:  teacherName,
		UnitPrice:    unitPrice,
		Amount:       roundMoney(amount),
		Chargeable:   chargeable,
		UseDate:      t,
	}

	err = eduClassSessionService.CreateEduClassSession(&classSession)

	if err != nil {
		return errors.New("创建课时操作记录失败")
	}

	return nil
}

// AddSession 为用户的课程添加课时
func (eduEnrollmentService *EduEnrollmentService) AddSession(userID, courseID, sessionsToAdd int, reason string, useData string, teacherId int, teacherName string) error {
	var enrollment edu_user_course.EduEnrollment
	db := global.GVA_DB

	// 查询用户在指定课程中的报名信息
	result := db.Where("user_id = ? AND course_id = ?", userID, courseID).Preload("EduCourse").First(&enrollment)

	// 如果找不到报名信息或发生错误，则返回错误
	if result.Error != nil {
		return errors.New("报名信息未找到或查询错误")
	}

	// 增加总课时
	if enrollment.TotalSessions == nil {
		enrollment.TotalSessions = new(int)
	}
	*enrollment.TotalSessions += sessionsToAdd

	// 增加剩余课时
	if enrollment.RemainingSessions == nil {
		enrollment.RemainingSessions = new(int)
	}
	*enrollment.RemainingSessions += sessionsToAdd

	t, err := parseUseDate(useData)
	if err != nil {
		return err
	}

	// 获取学员姓名
	var user struct {
		NickName string
	}
	err = db.Table("sys_users").Select("nick_name").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return errors.New("查询学员信息失败")
	}

	// 获取教师姓名（可选）
	if teacherName == "" && teacherId > 0 {
		var teacher struct {
			NickName string
		}
		if err := db.Table("sys_users").Select("nick_name").Where("id = ?", teacherId).First(&teacher).Error; err == nil {
			teacherName = teacher.NickName
		}
	}

	// 更新数据库中的报名信息
	applyEnrollmentFinance(&enrollment)
	result = db.Save(&enrollment)

	if result.Error != nil {
		return errors.New("更新报名信息失败")
	}

	enrollmentId := int(enrollment.ID)
	// 记录课时操作
	classSession := edu_user_course.EduClassSession{
		EnrollmentId: &enrollmentId,
		Action:       "add",
		Reason:       reason,
		NumSessions:  &sessionsToAdd,
		CourseName:   enrollment.EduCourse.CourseName,
		UserName:     user.NickName,
		TeacherId:    intPointer(teacherId),
		TeacherName:  teacherName,
		UnitPrice:    enrollment.PricePerSession,
		Amount:       0,
		Chargeable:   false,
		UseDate:      t,
	}

	result = db.Create(&classSession)

	if result.Error != nil {
		return errors.New("创建课时操作记录失败")
	}

	return nil
}

func parseUseDate(useData string) (time.Time, error) {
	trimmed := strings.TrimSpace(useData)
	if trimmed == "" {
		return time.Now(), nil
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, trimmed, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("日期转换失败: 格式应为 YYYY-MM-DD 或 YYYY-MM-DD HH:mm")
}

func intPointer(val int) *int {
	if val <= 0 {
		return nil
	}
	return &val
}

func applyEnrollmentDefaults(eduEnrollment *edu_user_course.EduEnrollment) {
	if eduEnrollment.TotalSessions == nil {
		zero := 0
		eduEnrollment.TotalSessions = &zero
	}
	if eduEnrollment.RemainingSessions == nil {
		total := *eduEnrollment.TotalSessions
		eduEnrollment.RemainingSessions = &total
	}
}

func applyEnrollmentFinance(eduEnrollment *edu_user_course.EduEnrollment) {
	totalSessions := 0
	if eduEnrollment.TotalSessions != nil {
		totalSessions = *eduEnrollment.TotalSessions
	}
	totalAmount := float64(totalSessions)*eduEnrollment.PricePerSession - eduEnrollment.DiscountAmount
	if totalAmount < 0 {
		totalAmount = 0
	}
	eduEnrollment.TotalAmount = roundMoney(totalAmount)
	eduEnrollment.BalanceAmount = roundMoney(eduEnrollment.TotalAmount - eduEnrollment.PaidAmount)
}

func roundMoney(val float64) float64 {
	return math.Round(val*100) / 100
}

// GetEduEnrollmentByUser 根据用户ID和课程ID获取报名信息
func (eduEnrollmentService *EduEnrollmentService) GetEduEnrollmentByUser(userId uint, courseId uint) (eduEnrollment edu_user_course.EduEnrollment, err error) {
	err = global.GVA_DB.Where("user_id = ? and course_id=?", userId, courseId).First(&eduEnrollment).Error
	return
}

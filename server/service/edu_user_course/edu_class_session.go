/*
 * @Author: Yang
 * @Date: 2023-04-10 18:33:08
 * @Description: 学生交易记录
 */
package edu_user_course

import (
	"time"

	"github.com/KeSilent/study-hub/server/global"
	"github.com/KeSilent/study-hub/server/model/common/request"
	"github.com/KeSilent/study-hub/server/model/edu_user_course"
	edu_user_courseReq "github.com/KeSilent/study-hub/server/model/edu_user_course/request"
	studentWithRemainingSessionsRes "github.com/KeSilent/study-hub/server/model/edu_user_course/response"
)

type EduClassSessionService struct {
}

// CreateEduClassSession 创建EduClassSession记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduClassSessionService *EduClassSessionService) CreateEduClassSession(eduClassSession *edu_user_course.EduClassSession) (err error) {
	err = global.GVA_DB.Create(eduClassSession).Error
	return err
}

// DeleteEduClassSession 删除EduClassSession记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduClassSessionService *EduClassSessionService) DeleteEduClassSession(eduClassSession edu_user_course.EduClassSession) (err error) {
	err = global.GVA_DB.Delete(&eduClassSession).Error
	return err
}

// DeleteEduClassSessionByIds 批量删除EduClassSession记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduClassSessionService *EduClassSessionService) DeleteEduClassSessionByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]edu_user_course.EduClassSession{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateEduClassSession 更新EduClassSession记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduClassSessionService *EduClassSessionService) UpdateEduClassSession(eduClassSession edu_user_course.EduClassSession) (err error) {
	err = global.GVA_DB.Save(&eduClassSession).Error
	return err
}

// GetEduClassSession 根据id获取EduClassSession记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduClassSessionService *EduClassSessionService) GetEduClassSession(id uint) (eduClassSession edu_user_course.EduClassSession, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&eduClassSession).Error
	return
}

// GetEduClassSessionInfoList 分页获取EduClassSession记录
// Author [piexlmax](https://github.com/piexlmax)
func (eduClassSessionService *EduClassSessionService) GetEduClassSessionInfoList(info edu_user_courseReq.EduClassSessionSearch) (list []edu_user_course.EduClassSession, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&edu_user_course.EduClassSession{})
	var eduClassSessions []edu_user_course.EduClassSession
	joinedEnrollment := false
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("edu_class_session.created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	// 按报名ID筛选
	if info.EnrollmentId != nil && *info.EnrollmentId > 0 {
		db = db.Where("edu_class_session.enrollment_id = ?", *info.EnrollmentId)
	}
	// 按教师筛选
	if info.TeacherId != nil && *info.TeacherId > 0 {
		db = db.Where("edu_class_session.teacher_id = ?", *info.TeacherId)
	}
	// 按课程筛选
	if info.CourseId > 0 {
		db = db.Joins("LEFT JOIN edu_enrollment ON edu_enrollment.id = edu_class_session.enrollment_id").
			Where("edu_enrollment.course_id = ?", info.CourseId)
		joinedEnrollment = true
	}

	// 按学员姓名搜索
	if info.UserName != "" {
		if !joinedEnrollment {
			db = db.Joins("LEFT JOIN edu_enrollment ON edu_enrollment.id = edu_class_session.enrollment_id")
			joinedEnrollment = true
		}
		db = db.Joins("LEFT JOIN sys_users ON sys_users.id = edu_enrollment.user_id").
			Where("sys_users.nick_name LIKE ?", "%"+info.UserName+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 按使用日期倒序排列，最新的在前面
	err = db.Order("edu_class_session.use_date DESC, edu_class_session.created_at DESC").
		Limit(limit).Offset(offset).Find(&eduClassSessions).Error
	if err != nil {
		return
	}

	// 填充学员姓名（如果user_name字段为空则从关联表查询）
	for i := range eduClassSessions {
		if eduClassSessions[i].UserName == "" && eduClassSessions[i].EnrollmentId != nil {
			// 通过enrollment_id关联查询用户姓名
			var result struct {
				NickName string
			}
			err := global.GVA_DB.Table("edu_enrollment").
				Select("sys_users.nick_name").
				Joins("LEFT JOIN sys_users ON sys_users.id = edu_enrollment.user_id").
				Where("edu_enrollment.id = ?", *eduClassSessions[i].EnrollmentId).
				First(&result).Error
			if err == nil {
				eduClassSessions[i].UserName = result.NickName
			}
		}
	}

	return eduClassSessions, total, err
}

// GetStudentClassSessions 获取EduClassSession记录
func (eduClassSessionService *EduClassSessionService) GetStudentClassSessions(studentId int) (list []edu_user_course.EduClassSession, err error) {
	var classSessions []edu_user_course.EduClassSession
	// 查询指定学生的所有课时记录
	result := global.GVA_DB.
		Preload("EduEnrollment.EduCourse").
		Joins("JOIN edu_enrollment on edu_enrollment.id = edu_class_session.enrollment_id").
		Where("edu_enrollment.user_id = ?", studentId).
		Order("edu_class_session.created_at DESC").
		Find(&classSessions)

	if result.Error != nil {
		return nil, result.Error
	}

	return classSessions, nil
}

// GetStudentsWithLessThanFiveSessions 获取剩余课时少于5节的学生列表
func (eduClassSessionService *EduClassSessionService) GetStudentsWithLessThanFiveSessions(organizationID uint) ([]studentWithRemainingSessionsRes.StudentWithRemainingSessions, error) {
	var students []studentWithRemainingSessionsRes.StudentWithRemainingSessions

	// 查询剩余课时少于5节的学生列表
	db := global.GVA_DB.Table("sys_users").
		Select("sys_users.*, edu_enrollment.remaining_sessions").
		Joins("JOIN edu_enrollment on edu_enrollment.user_id = sys_users.id").
		Where("edu_enrollment.remaining_sessions < ? ", 5)
	if organizationID > 0 {
		db = db.Where("sys_users.edu_organization_id= ?", organizationID)
	}
	result := db.Scan(&students)

	if result.Error != nil {
		return nil, result.Error
	}

	return students, nil
}

// GetMonthlyChargeSummary 获取指定月份计费汇总
func (eduClassSessionService *EduClassSessionService) GetMonthlyChargeSummary(month string) (studentWithRemainingSessionsRes.MonthlyChargeSummary, error) {
	summary := studentWithRemainingSessionsRes.MonthlyChargeSummary{
		Month: month,
	}
	if month == "" {
		month = time.Now().Format("2006-01")
		summary.Month = month
	}
	start, err := time.ParseInLocation("2006-01", month, time.Local)
	if err != nil {
		return summary, err
	}
	end := start.AddDate(0, 1, 0)

	type agg struct {
		ChargeableSessions    int64
		ChargeableAmount      float64
		NonChargeableSessions int64
	}
	var result agg

	err = global.GVA_DB.Model(&edu_user_course.EduClassSession{}).
		Where("action = ?", "subtract").
		Where("use_date >= ? AND use_date < ?", start, end).
		Select(`
			COALESCE(SUM(CASE WHEN chargeable = 1 OR chargeable IS NULL THEN num_sessions ELSE 0 END), 0) AS chargeable_sessions,
			COALESCE(SUM(CASE WHEN chargeable = 1 OR chargeable IS NULL THEN amount ELSE 0 END), 0) AS chargeable_amount,
			COALESCE(SUM(CASE WHEN chargeable = 0 THEN num_sessions ELSE 0 END), 0) AS non_chargeable_sessions
		`).Scan(&result).Error
	if err != nil {
		return summary, err
	}

	summary.ChargeableSessions = result.ChargeableSessions
	summary.ChargeableAmount = result.ChargeableAmount
	summary.NonChargeableSessions = result.NonChargeableSessions
	return summary, nil
}

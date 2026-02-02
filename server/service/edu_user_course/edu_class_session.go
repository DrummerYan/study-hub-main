/*
 * @Author: Yang
 * @Date: 2023-04-10 18:33:08
 * @Description: 学生交易记录
 */
package edu_user_course

import (
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
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("edu_class_session.created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	// 按报名ID筛选
	if info.EnrollmentId != nil && *info.EnrollmentId > 0 {
		db = db.Where("edu_class_session.enrollment_id = ?", *info.EnrollmentId)
	}
	
	// 按学员姓名搜索
	if info.UserName != "" {
		db = db.Joins("LEFT JOIN edu_enrollment ON edu_enrollment.id = edu_class_session.enrollment_id").
			Joins("LEFT JOIN sys_users ON sys_users.id = edu_enrollment.user_id").
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
	result := global.GVA_DB.Table("sys_users").
		Select("sys_users.*, edu_enrollment.remaining_sessions").
		Joins("JOIN edu_enrollment on edu_enrollment.user_id = sys_users.id").
		Where("edu_enrollment.remaining_sessions < ? and sys_users.edu_organization_id= ? ", 5, organizationID).
		Scan(&students)

	if result.Error != nil {
		return nil, result.Error
	}

	return students, nil
}

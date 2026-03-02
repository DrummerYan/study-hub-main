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
	"gorm.io/gorm"
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
	paidRemain := getIntValue(enrollment.RemainingPaid, enrollment.RemainingSessions)
	giftRemain := getIntValue(enrollment.RemainingGift, nil)
	if paidRemain+giftRemain < sessionsToConsume {
		return errors.New("剩余课时不足")
	}

	paidConsume := sessionsToConsume
	if paidConsume > paidRemain {
		paidConsume = paidRemain
	}
	giftConsume := sessionsToConsume - paidConsume

	paidRemain -= paidConsume
	giftRemain -= giftConsume
	enrollment.RemainingPaid = intPointerAllowZero(paidRemain)
	enrollment.RemainingGift = intPointerAllowZero(giftRemain)
	totalRemain := paidRemain + giftRemain
	enrollment.RemainingSessions = intPointerAllowZero(totalRemain)

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
	// 记录课时操作（付费优先，赠课在后）
	unitPrice := enrollment.PricePerSession
	if paidConsume > 0 {
		amount := 0.0
		if chargeable {
			amount = float64(paidConsume) * unitPrice
		}
		classSession := edu_user_course.EduClassSession{
			EnrollmentId: &enenrollmentId,
			Action:       "subtract",
			Reason:       reason,
			NumSessions:  intPointerAllowZero(paidConsume),
			CourseName:   enrollment.EduCourse.CourseName,
			UserName:     user.NickName,
			TeacherId:    intPointer(teacherId),
			TeacherName:  teacherName,
			UnitPrice:    unitPrice,
			Amount:       roundMoney(amount),
			Chargeable:   chargeable,
			UseDate:      t,
		}
		if err = eduClassSessionService.CreateEduClassSession(&classSession); err != nil {
			return errors.New("创建课时操作记录失败")
		}
	}
	if giftConsume > 0 {
		reasonGift := reason
		if !strings.Contains(reasonGift, "赠课") {
			reasonGift = reasonGift + "（赠课）"
		}
		classSession := edu_user_course.EduClassSession{
			EnrollmentId: &enenrollmentId,
			Action:       "subtract",
			Reason:       reasonGift,
			NumSessions:  intPointerAllowZero(giftConsume),
			CourseName:   enrollment.EduCourse.CourseName,
			UserName:     user.NickName,
			TeacherId:    intPointer(teacherId),
			TeacherName:  teacherName,
			UnitPrice:    unitPrice,
			Amount:       0,
			Chargeable:   false,
			UseDate:      t,
		}
		if err = eduClassSessionService.CreateEduClassSession(&classSession); err != nil {
			return errors.New("创建课时操作记录失败")
		}
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
	if enrollment.PaidSessions == nil {
		enrollment.PaidSessions = new(int)
	}
	if enrollment.GiftSessions == nil {
		enrollment.GiftSessions = new(int)
	}
	if enrollment.RemainingPaid == nil {
		enrollment.RemainingPaid = new(int)
	}
	if enrollment.RemainingGift == nil {
		enrollment.RemainingGift = new(int)
	}
	if isGiftAddReason(reason) {
		*enrollment.GiftSessions += sessionsToAdd
		*enrollment.RemainingGift += sessionsToAdd
	} else {
		*enrollment.PaidSessions += sessionsToAdd
		*enrollment.RemainingPaid += sessionsToAdd
	}

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

func intPointerAllowZero(val int) *int {
	return &val
}

func getIntValue(val *int, fallback *int) int {
	if val != nil {
		return *val
	}
	if fallback != nil {
		return *fallback
	}
	return 0
}

func applyEnrollmentDefaults(eduEnrollment *edu_user_course.EduEnrollment) {
	if eduEnrollment.PaidSessions == nil {
		zero := 0
		eduEnrollment.PaidSessions = intPointerAllowZero(zero)
	}
	if eduEnrollment.GiftSessions == nil {
		zero := 0
		eduEnrollment.GiftSessions = intPointerAllowZero(zero)
	}
	// 兼容旧数据：仅有总课时时默认全部为付费课时
	if getIntValue(eduEnrollment.PaidSessions, nil) == 0 && getIntValue(eduEnrollment.GiftSessions, nil) == 0 {
		if eduEnrollment.TotalSessions != nil && *eduEnrollment.TotalSessions > 0 {
			eduEnrollment.PaidSessions = intPointerAllowZero(*eduEnrollment.TotalSessions)
		}
	}
	total := getIntValue(eduEnrollment.PaidSessions, nil) + getIntValue(eduEnrollment.GiftSessions, nil)
	eduEnrollment.TotalSessions = &total
	if eduEnrollment.RemainingPaid == nil {
		paid := getIntValue(eduEnrollment.PaidSessions, nil)
		eduEnrollment.RemainingPaid = intPointerAllowZero(paid)
	}
	if eduEnrollment.RemainingGift == nil {
		gift := getIntValue(eduEnrollment.GiftSessions, nil)
		eduEnrollment.RemainingGift = intPointerAllowZero(gift)
	}
	// 兼容旧数据：仅有剩余总课时时默认全部为付费剩余
	if getIntValue(eduEnrollment.RemainingPaid, nil) == 0 && getIntValue(eduEnrollment.RemainingGift, nil) == 0 {
		if eduEnrollment.RemainingSessions != nil && *eduEnrollment.RemainingSessions > 0 {
			eduEnrollment.RemainingPaid = intPointerAllowZero(*eduEnrollment.RemainingSessions)
		}
	}
	remainingTotal := getIntValue(eduEnrollment.RemainingPaid, nil) + getIntValue(eduEnrollment.RemainingGift, nil)
	eduEnrollment.RemainingSessions = intPointerAllowZero(remainingTotal)
}

func applyEnrollmentFinance(eduEnrollment *edu_user_course.EduEnrollment) {
	paidSessions := 0
	if eduEnrollment.PaidSessions != nil {
		paidSessions = *eduEnrollment.PaidSessions
	}
	totalAmount := float64(paidSessions)*eduEnrollment.PricePerSession - eduEnrollment.DiscountAmount
	if totalAmount < 0 {
		totalAmount = 0
	}
	eduEnrollment.TotalAmount = roundMoney(totalAmount)
	eduEnrollment.BalanceAmount = roundMoney(eduEnrollment.TotalAmount - eduEnrollment.PaidAmount)
}

func roundMoney(val float64) float64 {
	return math.Round(val*100) / 100
}

func isGiftAddReason(reason string) bool {
	return strings.Contains(reason, "赠")
}

// RefundOrTransfer 退费/转课（仅处理付费课时，赠课不计入可退范围）
func (eduEnrollmentService *EduEnrollmentService) RefundOrTransfer(req edu_user_courseReq.RefundTransferReq) (edu_user_courseReq.RefundTransferResp, error) {
	var resp edu_user_courseReq.RefundTransferResp
	if req.EnrollmentId == 0 {
		return resp, errors.New("报名信息不存在")
	}
	action := strings.ToLower(strings.TrimSpace(req.Action))
	if action == "" {
		action = "refund"
	}
	if action != "refund" && action != "transfer" {
		return resp, errors.New("操作类型不合法")
	}
	if action == "transfer" && req.TargetCourseId == 0 {
		return resp, errors.New("请选择转入课程")
	}

	tx := global.GVA_DB.Begin()
	if tx.Error != nil {
		return resp, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var enrollment edu_user_course.EduEnrollment
	if err := tx.Preload("EduCourse").Where("id = ?", req.EnrollmentId).First(&enrollment).Error; err != nil {
		tx.Rollback()
		return resp, errors.New("报名信息未找到")
	}
	applyEnrollmentDefaults(&enrollment)

	userName := ""
	if enrollment.UserId != nil {
		var user struct {
			NickName string
		}
		if err := tx.Table("sys_users").Select("nick_name").Where("id = ?", *enrollment.UserId).First(&user).Error; err == nil {
			userName = user.NickName
		}
	}

	paidTotal := getIntValue(enrollment.PaidSessions, nil)
	remainingPaid := getIntValue(enrollment.RemainingPaid, nil)
	if paidTotal <= 0 {
		tx.Rollback()
		return resp, errors.New("没有付费课时")
	}
	if remainingPaid <= 0 {
		tx.Rollback()
		return resp, errors.New("剩余付费课时为0")
	}
	refundSessions := req.Sessions
	if refundSessions <= 0 {
		refundSessions = remainingPaid
	}
	if refundSessions > remainingPaid {
		tx.Rollback()
		return resp, errors.New("退费课时不能大于剩余付费课时")
	}

	paidUnitPrice := 0.0
	if enrollment.TotalAmount > 0 && paidTotal > 0 {
		paidUnitPrice = enrollment.TotalAmount / float64(paidTotal)
	} else {
		paidUnitPrice = enrollment.PricePerSession
	}
	refundAmount := roundMoney(float64(refundSessions) * paidUnitPrice)

	// 更新报名信息（减少付费课时）
	consumedPaid := paidTotal - remainingPaid
	newPaidTotal := paidTotal - refundSessions
	if newPaidTotal < consumedPaid {
		tx.Rollback()
		return resp, errors.New("可退付费课时不足")
	}
	newRemainingPaid := remainingPaid - refundSessions

	enrollment.PaidSessions = intPointerAllowZero(newPaidTotal)
	enrollment.RemainingPaid = intPointerAllowZero(newRemainingPaid)

	// 退费/转课后，剩余赠课不再保留（已消耗赠课保留为历史）
	giftTotal := getIntValue(enrollment.GiftSessions, nil)
	giftRemain := getIntValue(enrollment.RemainingGift, nil)
	giftConsumed := giftTotal - giftRemain
	if giftRemain > 0 {
		enrollment.GiftSessions = intPointerAllowZero(giftConsumed)
		enrollment.RemainingGift = intPointerAllowZero(0)
	}
	remainingGift := getIntValue(enrollment.RemainingGift, nil)
	enrollment.RemainingSessions = intPointerAllowZero(newRemainingPaid + remainingGift)

	applyEnrollmentDefaults(&enrollment)
	applyEnrollmentFinance(&enrollment)
	if err := tx.Save(&enrollment).Error; err != nil {
		tx.Rollback()
		return resp, errors.New("更新报名信息失败")
	}

	// 记录退费/转课流水
	refundRecord := edu_user_course.EduRefund{
		EnrollmentId:   intPointerAllowZero(int(req.EnrollmentId)),
		UserId:         enrollment.UserId,
		CourseId:       enrollment.CourseId,
		Action:         action,
		RefundSessions: refundSessions,
		RefundAmount:   refundAmount,
		Reason:         strings.TrimSpace(req.Reason),
		RefundTime:     time.Now(),
	}
	if req.OperatorId > 0 {
		refundRecord.OperatorId = intPointerAllowZero(req.OperatorId)
	}
	if req.OperatorName != "" {
		refundRecord.OperatorName = req.OperatorName
	}
	if err := tx.Create(&refundRecord).Error; err != nil {
		tx.Rollback()
		return resp, errors.New("记录退费失败")
	}

	// 转课：将剩余付费课时转入目标课程（按当前付费单价）
	if action == "transfer" {
		if enrollment.UserId == nil {
			tx.Rollback()
			return resp, errors.New("报名用户不存在")
		}
		var target edu_user_course.EduEnrollment
		err := tx.Preload("EduCourse").Where("user_id = ? AND course_id = ?", *enrollment.UserId, req.TargetCourseId).First(&target).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return resp, errors.New("查询转入课程失败")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			target = edu_user_course.EduEnrollment{
				UserId:          enrollment.UserId,
				CourseId:        &req.TargetCourseId,
				PaidSessions:    intPointerAllowZero(refundSessions),
				GiftSessions:    intPointerAllowZero(0),
				RemainingPaid:   intPointerAllowZero(refundSessions),
				RemainingGift:   intPointerAllowZero(0),
				PricePerSession: paidUnitPrice,
			}
			applyEnrollmentDefaults(&target)
			applyEnrollmentFinance(&target)
			if err := tx.Create(&target).Error; err != nil {
				tx.Rollback()
				return resp, errors.New("创建转入报名失败")
			}
		} else {
			applyEnrollmentDefaults(&target)
			targetPaid := getIntValue(target.PaidSessions, nil) + refundSessions
			targetRemain := getIntValue(target.RemainingPaid, nil) + refundSessions
			target.PaidSessions = intPointerAllowZero(targetPaid)
			target.RemainingPaid = intPointerAllowZero(targetRemain)
			if target.PricePerSession <= 0 {
				target.PricePerSession = paidUnitPrice
			}
			applyEnrollmentDefaults(&target)
			applyEnrollmentFinance(&target)
			if err := tx.Save(&target).Error; err != nil {
				tx.Rollback()
				return resp, errors.New("更新转入报名失败")
			}
		}

		// 记录转课流水（不计费）
		now := time.Now()
		sourceEnrollmentId := int(enrollment.ID)
		targetEnrollmentId := int(target.ID)
		transferOut := edu_user_course.EduClassSession{
			EnrollmentId: &sourceEnrollmentId,
			Action:       "subtract",
			Reason:       "转课转出",
			NumSessions:  intPointerAllowZero(refundSessions),
			CourseName:   enrollment.EduCourse.CourseName,
			UserName:     userName,
			UnitPrice:    paidUnitPrice,
			Amount:       0,
			Chargeable:   false,
			UseDate:      now,
		}
		targetCourseName := target.EduCourse.CourseName
		if targetCourseName == "" && req.TargetCourseId != 0 {
			var course struct {
				CourseName string `gorm:"column:course_name"`
			}
			if err := tx.Table("edu_course").Select("course_name").Where("id = ?", req.TargetCourseId).First(&course).Error; err == nil {
				targetCourseName = course.CourseName
			}
		}
		transferIn := edu_user_course.EduClassSession{
			EnrollmentId: &targetEnrollmentId,
			Action:       "add",
			Reason:       "转课转入",
			NumSessions:  intPointerAllowZero(refundSessions),
			CourseName:   targetCourseName,
			UserName:     userName,
			UnitPrice:    target.PricePerSession,
			Amount:       0,
			Chargeable:   false,
			UseDate:      now,
		}
		if err := tx.Create(&transferOut).Error; err != nil {
			tx.Rollback()
			return resp, errors.New("记录转出失败")
		}
		if err := tx.Create(&transferIn).Error; err != nil {
			tx.Rollback()
			return resp, errors.New("记录转入失败")
		}
	}

	// 退费后若学员没有任何剩余课时，则自动停用（可在用户管理中恢复）
	if action == "refund" && enrollment.UserId != nil {
		var activeCount int64
		if err := tx.Table("edu_enrollment").
			Where("user_id = ?", *enrollment.UserId).
			Where("(COALESCE(remaining_paid_sessions,0) + COALESCE(remaining_gift_sessions,0) > 0) OR COALESCE(remaining_sessions,0) > 0").
			Count(&activeCount).Error; err == nil {
			if activeCount == 0 {
				_ = tx.Table("sys_users").Where("id = ?", *enrollment.UserId).Update("enable", 2).Error
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return resp, err
	}
	resp.RefundSessions = refundSessions
	resp.PaidUnitPrice = roundMoney(paidUnitPrice)
	resp.RefundAmount = refundAmount
	return resp, nil
}

// GetEduEnrollmentByUser 根据用户ID和课程ID获取报名信息
func (eduEnrollmentService *EduEnrollmentService) GetEduEnrollmentByUser(userId uint, courseId uint) (eduEnrollment edu_user_course.EduEnrollment, err error) {
	err = global.GVA_DB.Where("user_id = ? and course_id=?", userId, courseId).First(&eduEnrollment).Error
	return
}

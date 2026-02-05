package edu_user_course

import (
	"errors"
	"time"

	"github.com/KeSilent/study-hub/server/global"
	"github.com/KeSilent/study-hub/server/model/edu_user_course"
	edu_user_courseReq "github.com/KeSilent/study-hub/server/model/edu_user_course/request"
	"gorm.io/gorm"
)

type EduPaymentService struct{}

// CreateEduPayment 创建收款记录并更新报名已收/余额
func (eduPaymentService *EduPaymentService) CreateEduPayment(payment *edu_user_course.EduPayment) error {
	if payment.EnrollmentId == nil || *payment.EnrollmentId == 0 {
		return errors.New("报名信息不能为空")
	}
	if payment.Amount <= 0 {
		return errors.New("收款金额必须大于0")
	}
	if payment.PayTime.IsZero() {
		payment.PayTime = time.Now()
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 查询报名
		var enrollment edu_user_course.EduEnrollment
		if err := tx.Where("id = ?", *payment.EnrollmentId).First(&enrollment).Error; err != nil {
			return err
		}
		// 创建收款记录
		if err := tx.Create(payment).Error; err != nil {
			return err
		}
		// 更新已收款与余额
		enrollment.PaidAmount = enrollment.PaidAmount + payment.Amount
		enrollment.BalanceAmount = enrollment.TotalAmount - enrollment.PaidAmount
		return tx.Save(&enrollment).Error
	})
}

// GetEduPaymentInfoList 分页获取收款记录
func (eduPaymentService *EduPaymentService) GetEduPaymentInfoList(info edu_user_courseReq.EduPaymentSearch) (list []edu_user_course.EduPayment, total int64, err error) {
	limit := info.PageSize
	if limit == 0 {
		limit = 10
	}
	page := info.Page
	if page == 0 {
		page = 1
	}
	offset := limit * (page - 1)
	db := global.GVA_DB.Model(&edu_user_course.EduPayment{})

	if info.EnrollmentId > 0 {
		db = db.Where("enrollment_id = ?", info.EnrollmentId)
	}
	if info.StartPayTime != nil && info.EndPayTime != nil {
		db = db.Where("pay_time BETWEEN ? AND ?", info.StartPayTime, info.EndPayTime)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("pay_time DESC, created_at DESC").
		Limit(limit).Offset(offset).Find(&list).Error
	return
}

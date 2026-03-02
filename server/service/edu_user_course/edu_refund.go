package edu_user_course

import (
	"time"

	"github.com/KeSilent/study-hub/server/global"
	"github.com/KeSilent/study-hub/server/model/edu_user_course"
	edu_user_courseReq "github.com/KeSilent/study-hub/server/model/edu_user_course/request"
)

type EduRefundService struct{}

// CreateEduRefund 创建退费记录
func (eduRefundService *EduRefundService) CreateEduRefund(refund *edu_user_course.EduRefund) error {
	if refund.RefundTime.IsZero() {
		refund.RefundTime = time.Now()
	}
	return global.GVA_DB.Create(refund).Error
}

// GetEduRefundInfoList 分页获取退费记录
func (eduRefundService *EduRefundService) GetEduRefundInfoList(info edu_user_courseReq.EduRefundSearch) (list []edu_user_course.EduRefund, total int64, err error) {
	limit := info.PageSize
	if limit == 0 {
		limit = 10
	}
	page := info.Page
	if page == 0 {
		page = 1
	}
	offset := limit * (page - 1)
	db := global.GVA_DB.Model(&edu_user_course.EduRefund{})

	if info.EnrollmentId > 0 {
		db = db.Where("enrollment_id = ?", info.EnrollmentId)
	}
	if info.Action != "" {
		db = db.Where("action = ?", info.Action)
	}
	if info.StartRefundTime != nil && info.EndRefundTime != nil {
		db = db.Where("refund_time BETWEEN ? AND ?", info.StartRefundTime, info.EndRefundTime)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("refund_time DESC, created_at DESC").
		Limit(limit).Offset(offset).Find(&list).Error
	return
}

package edu_user_course

import (
	"errors"
	"strings"
	"time"

	"github.com/KeSilent/study-hub/server/global"
	"github.com/KeSilent/study-hub/server/model/common/response"
	"github.com/KeSilent/study-hub/server/model/edu_user_course"
	edu_user_courseReq "github.com/KeSilent/study-hub/server/model/edu_user_course/request"
	"github.com/KeSilent/study-hub/server/service"
	"github.com/KeSilent/study-hub/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EduPaymentApi struct{}

var eduPaymentService = service.ServiceGroupApp.Edu_user_courseServiceGroup.EduPaymentService

// CreateEduPayment 创建收款记录
// @Tags EduPayment
// @Summary 创建收款记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body edu_user_course.EduPayment true "创建收款记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /eduPayment/createEduPayment [post]
func (eduPaymentApi *EduPaymentApi) CreateEduPayment(c *gin.Context) {
	if !utils.IsSuperAdmin(c) {
		response.FailWithMessage("权限不足", c)
		return
	}
	var req edu_user_courseReq.EduPaymentCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	payTime, err := parsePaymentTime(req.PayTime)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	payment := edu_user_course.EduPayment{
		EnrollmentId: intPointer(req.EnrollmentId),
		Amount:       req.Amount,
		PayTime:      payTime,
		Remark:       req.Remark,
	}
	if claims := utils.GetUserInfo(c); claims != nil {
		operatorId := int(claims.BaseClaims.ID)
		payment.OperatorId = &operatorId
		payment.OperatorName = claims.BaseClaims.NickName
	}
	if err := eduPaymentService.CreateEduPayment(&payment); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetEduPaymentList 分页获取收款记录
// @Tags EduPayment
// @Summary 分页获取收款记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query edu_user_courseReq.EduPaymentSearch true "分页获取收款记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /eduPayment/getEduPaymentList [get]
func (eduPaymentApi *EduPaymentApi) GetEduPaymentList(c *gin.Context) {
	if !utils.IsSuperAdmin(c) {
		response.FailWithMessage("权限不足", c)
		return
	}
	var pageInfo edu_user_courseReq.EduPaymentSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := eduPaymentService.GetEduPaymentInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

func parsePaymentTime(value string) (time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Now(), nil
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		time.RFC3339,
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

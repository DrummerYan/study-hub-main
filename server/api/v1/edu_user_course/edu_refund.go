package edu_user_course

import (
	"github.com/KeSilent/study-hub/server/global"
	"github.com/KeSilent/study-hub/server/model/common/response"
	edu_user_courseReq "github.com/KeSilent/study-hub/server/model/edu_user_course/request"
	"github.com/KeSilent/study-hub/server/service"
	"github.com/KeSilent/study-hub/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EduRefundApi struct{}

var eduRefundService = service.ServiceGroupApp.Edu_user_courseServiceGroup.EduRefundService

// GetEduRefundList 分页获取退费记录
// @Tags EduRefund
// @Summary 分页获取退费记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query edu_user_courseReq.EduRefundSearch true "分页获取退费记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /eduRefund/getEduRefundList [get]
func (eduRefundApi *EduRefundApi) GetEduRefundList(c *gin.Context) {
	if !utils.IsSuperAdmin(c) {
		response.FailWithMessage("权限不足", c)
		return
	}
	var pageInfo edu_user_courseReq.EduRefundSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := eduRefundService.GetEduRefundInfoList(pageInfo); err != nil {
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

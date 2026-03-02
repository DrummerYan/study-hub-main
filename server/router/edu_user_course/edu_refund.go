package edu_user_course

import (
	v1 "github.com/KeSilent/study-hub/server/api/v1"
	"github.com/gin-gonic/gin"
)

type EduRefundRouter struct{}

// InitEduRefundRouter 初始化 EduRefund 路由信息
func (s *EduRefundRouter) InitEduRefundRouter(Router *gin.RouterGroup) {
	eduRefundRouterWithoutRecord := Router.Group("eduRefund")
	var eduRefundApi = v1.ApiGroupApp.Edu_user_courseApiGroup.EduRefundApi
	{
		eduRefundRouterWithoutRecord.GET("getEduRefundList", eduRefundApi.GetEduRefundList) // 获取退费记录
	}
}

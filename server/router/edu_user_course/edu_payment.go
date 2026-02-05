package edu_user_course

import (
	v1 "github.com/KeSilent/study-hub/server/api/v1"
	"github.com/KeSilent/study-hub/server/middleware"
	"github.com/gin-gonic/gin"
)

type EduPaymentRouter struct{}

// InitEduPaymentRouter 初始化 EduPayment 路由信息
func (s *EduPaymentRouter) InitEduPaymentRouter(Router *gin.RouterGroup) {
	eduPaymentRouter := Router.Group("eduPayment").Use(middleware.OperationRecord())
	eduPaymentRouterWithoutRecord := Router.Group("eduPayment")
	var eduPaymentApi = v1.ApiGroupApp.Edu_user_courseApiGroup.EduPaymentApi
	{
		eduPaymentRouter.POST("createEduPayment", eduPaymentApi.CreateEduPayment) // 新建收款记录
	}
	{
		eduPaymentRouterWithoutRecord.GET("getEduPaymentList", eduPaymentApi.GetEduPaymentList) // 获取收款记录
	}
}

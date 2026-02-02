package middleware

import (
	"strconv"
	"strings"

	"github.com/KeSilent/study-hub/server/global"
	"github.com/KeSilent/study-hub/server/model/common/response"
	"github.com/KeSilent/study-hub/server/service"
	"github.com/KeSilent/study-hub/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var casbinService = service.ServiceGroupApp.SystemServiceGroup.CasbinService

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.GVA_CONFIG.System.Env != "develop" {
			waitUse, _ := utils.GetClaims(c)
			//获取请求的PATH
			path := c.Request.URL.Path
			obj := strings.TrimPrefix(path, global.GVA_CONFIG.System.RouterPrefix)
			// 获取请求方法
			act := c.Request.Method
			// 获取用户的角色
			sub := strconv.Itoa(int(waitUse.AuthorityId))
			e := casbinService.Casbin() // 判断策略中是否存在
			success, _ := e.Enforce(sub, obj, act)
			
			// 添加调试日志
			global.GVA_LOG.Info("Casbin权限检查",
				zap.String("用户角色", sub),
				zap.String("请求路径", obj),
				zap.String("请求方法", act),
				zap.Bool("是否通过", success))
			
			if !success {
				global.GVA_LOG.Warn("权限不足",
					zap.String("用户角色", sub),
					zap.String("请求路径", obj),
					zap.String("请求方法", act))
				response.FailWithDetailed(gin.H{}, "权限不足", c)
				c.Abort()
				return
			}
		} else {
			global.GVA_LOG.Info("Casbin跳过检查（开发环境）", zap.String("path", c.Request.URL.Path))
		}
		c.Next()
	}
}

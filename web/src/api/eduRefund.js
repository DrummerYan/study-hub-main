import service from '@/utils/request'

// @Tags EduRefund
// @Summary 分页获取退费记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query object true "分页获取退费记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /eduRefund/getEduRefundList [get]
export const getEduRefundList = (params) => {
  return service({
    url: '/eduRefund/getEduRefundList',
    method: 'get',
    params
  })
}

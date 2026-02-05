import service from '@/utils/request'

// @Tags EduPayment
// @Summary 创建收款记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body object true "创建收款记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /eduPayment/createEduPayment [post]
export const createEduPayment = (data) => {
  return service({
    url: '/eduPayment/createEduPayment',
    method: 'post',
    data
  })
}

// @Tags EduPayment
// @Summary 分页获取收款记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query object true "分页获取收款记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /eduPayment/getEduPaymentList [get]
export const getEduPaymentList = (params) => {
  return service({
    url: '/eduPayment/getEduPaymentList',
    method: 'get',
    params
  })
}

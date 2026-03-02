package request

// RefundTransferReq 退费/转课请求
type RefundTransferReq struct {
	EnrollmentId   uint   `json:"enrollmentId"`             // 报名ID
	Action         string `json:"action"`                   // refund | transfer
	Sessions       int    `json:"sessions"`                 // 退费/转课的付费课时数（<=剩余付费课时，0表示全部）
	TargetCourseId int    `json:"targetCourseId,omitempty"` // 转课目标课程ID
	Reason         string `json:"reason"`                   // 退费/转课原因（可选）
	OperatorId     int    `json:"operatorId"`               // 操作人ID
	OperatorName   string `json:"operatorName"`             // 操作人姓名
}

// RefundTransferResp 退费/转课响应
type RefundTransferResp struct {
	RefundSessions int     `json:"refundSessions"`
	PaidUnitPrice  float64 `json:"paidUnitPrice"`
	RefundAmount   float64 `json:"refundAmount"`
}

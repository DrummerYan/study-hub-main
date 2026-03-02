package request

import (
	"time"

	"github.com/KeSilent/study-hub/server/model/common/request"
)

// EduRefundSearch 退费记录查询
type EduRefundSearch struct {
	EnrollmentId    uint       `json:"enrollmentId" form:"enrollmentId"`
	Action          string     `json:"action" form:"action"`
	StartRefundTime *time.Time `json:"startRefundTime" form:"startRefundTime"`
	EndRefundTime   *time.Time `json:"endRefundTime" form:"endRefundTime"`
	request.PageInfo
}

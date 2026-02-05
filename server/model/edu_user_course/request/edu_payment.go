package request

import (
	"time"

	"github.com/KeSilent/study-hub/server/model/common/request"
)

type EduPaymentSearch struct {
	EnrollmentId uint       `json:"enrollmentId" form:"enrollmentId"`
	StartPayTime *time.Time `json:"startPayTime" form:"startPayTime"`
	EndPayTime   *time.Time `json:"endPayTime" form:"endPayTime"`
	request.PageInfo
}

type EduPaymentCreate struct {
	EnrollmentId int     `json:"enrollmentId"`
	Amount       float64 `json:"amount"`
	PayTime      string  `json:"payTime"`
	Remark       string  `json:"remark"`
}

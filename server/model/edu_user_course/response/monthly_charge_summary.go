package response

type MonthlyChargeSummary struct {
	Month                 string  `json:"month"`
	ChargeableSessions    int64   `json:"chargeableSessions"`
	ChargeableAmount      float64 `json:"chargeableAmount"`
	NonChargeableSessions int64   `json:"nonChargeableSessions"`
}

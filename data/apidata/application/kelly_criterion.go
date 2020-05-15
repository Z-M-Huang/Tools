package application

//KellyCriterionRequest /api/kelly-criterion/simulate request
type KellyCriterionRequest struct {
	MaxWinChancePayout float64 `json:"maxWinChancePayout" xml:"maxWinChancePayout" form:"maxWinChancePayout" binding:"required"`
	MaxWinChance       float64 `json:"maxWinChance" xml:"maxWinChance" form:"maxWinChance" binding:"required"`
}

//KellyCriterionSimulationResponse /api/kelly-criterion/simulate response. Response will be a slice
type KellyCriterionSimulationResponse struct {
	Payout float64
	Chance float64
	Factor float64
}

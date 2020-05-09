package application

//KellyCriterionRequest /api/kelly-criterion/simulate request
type KellyCriterionRequest struct {
	MaxWinChancePayout string `json:"maxWinChancePayout"`
	MaxWinChance       string `json:"maxWinChance"`
}

//KellyCriterionSimulationResponse /api/kelly-criterion/simulate response. Response will be a slice
type KellyCriterionSimulationResponse struct {
	Payout float64
	Chance float64
	Factor float64
}

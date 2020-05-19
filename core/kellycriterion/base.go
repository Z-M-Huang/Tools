package kellycriterion

//Request /api/kelly-criterion/simulate request
type Request struct {
	MaxWinChancePayout float64 `json:"maxWinChancePayout" xml:"maxWinChancePayout" form:"maxWinChancePayout" binding:"required"`
	MaxWinChance       float64 `json:"maxWinChance" xml:"maxWinChance" form:"maxWinChance" binding:"required"`
}

//Response /api/kelly-criterion/simulate response. Response will be a slice
type Response struct {
	Payout float64
	Chance float64
	Factor float64
}

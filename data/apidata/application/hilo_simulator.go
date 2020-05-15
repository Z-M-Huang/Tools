package application

//HiLoSimulateRequest simulation request
type HiLoSimulateRequest struct {
	TotalStack             string `json:"totalStack"`
	WinChance              string `json:"winChance"`
	Odds                   string `json:"odds"`
	BaseBet                string `json:"baseBet"`
	RollAmount             string `json:"rollAmount"`
	OnWinReturnToBaseBet   string `json:"onWinReturnToBaseBet"`
	OnWinIncreaseBy        string `json:"onWinIncreaseBy"`
	OnWinChangeOdds        string `json:"onWinChangeOdds"`
	OnWinChangeOddsTo      string `json:"onWinChangeOddsTo"`
	OnWinNewOddsWinChance  string `json:"onWinNewOddsWinChance"`
	OnLossReturnToBaseBet  string `json:"onLossReturnToBaseBet"`
	OnLossIncreaseBy       string `json:"onLossIncreaseBy"`
	OnLossChangeOdds       string `json:"onLossChangeOdds"`
	OnLossChangeOddsTo     string `json:"onLossChangeOddsTo"`
	OnLossNewOddsWinChance string `json:"onLossNewOddsWinChance"`
	RandomClientSeed       string `json:"randomClientSeed"`
	AlternateHiLo          string `json:"alternateHiLo"`
}

//HiLoVerifyRequest verify request
type HiLoVerifyRequest struct {
	ServerSeed string
	ClientSeed string
	Nonce      uint64
	Roll       float64
}

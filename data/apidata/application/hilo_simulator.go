package application

//HiLoSimulateRequest simulation request
type HiLoSimulateRequest struct {
	TotalStack             string `json:"totalStack" xml:"totalStack" form:"totalStack" binding:"required"`
	WinChance              string `json:"winChance" xml:"winChance" form:"winChance" binding:"required"`
	Odds                   string `json:"odds" xml:"odds" form:"odds" binding:"required"`
	BaseBet                string `json:"baseBet" xml:"baseBet" form:"baseBet" binding:"required"`
	RollAmount             string `json:"rollAmount" xml:"rollAmount" form:"rollAmount" binding:"required"`
	OnWinReturnToBaseBet   string `json:"onWinReturnToBaseBet" xml:"onWinReturnToBaseBet" form:"onWinReturnToBaseBet"`
	OnWinIncreaseBy        string `json:"onWinIncreaseBy" xml:"onWinIncreaseBy" form:"onWinIncreaseBy"`
	OnWinChangeOdds        string `json:"onWinChangeOdds" xml:"onWinChangeOdds" form:"onWinChangeOdds"`
	OnWinChangeOddsTo      string `json:"onWinChangeOddsTo" xml:"onWinChangeOddsTo" form:"onWinChangeOddsTo"`
	OnWinNewOddsWinChance  string `json:"onWinNewOddsWinChance" xml:"onWinNewOddsWinChance" form:"onWinNewOddsWinChance"`
	OnLossReturnToBaseBet  string `json:"onLossReturnToBaseBet" xml:"onLossReturnToBaseBet" form:"onLossReturnToBaseBet"`
	OnLossIncreaseBy       string `json:"onLossIncreaseBy" xml:"onLossIncreaseBy" form:"onLossIncreaseBy"`
	OnLossChangeOdds       string `json:"onLossChangeOdds" xml:"onLossChangeOdds" form:"onLossChangeOdds"`
	OnLossChangeOddsTo     string `json:"onLossChangeOddsTo" xml:"onLossChangeOddsTo" form:"onLossChangeOddsTo"`
	OnLossNewOddsWinChance string `json:"onLossNewOddsWinChance" xml:"onLossNewOddsWinChance" form:"onLossNewOddsWinChance"`
	RandomClientSeed       string `json:"randomClientSeed" xml:"randomClientSeed" form:"randomClientSeed"`
	AlternateHiLo          string `json:"alternateHiLo" xml:"alternateHiLo" form:"alternateHiLo"`
}

//HiLoVerifyRequest verify request
type HiLoVerifyRequest struct {
	ServerSeed string
	ClientSeed string
	Nonce      uint64
	Roll       float64
}

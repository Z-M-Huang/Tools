package application

//HiLoSimulateRequest simulation request
type HiLoSimulateRequest struct {
	TotalStack             float64 `json:"totalStack" xml:"totalStack" form:"totalStack" binding:"required"`
	WinChance              float64 `json:"winChance" xml:"winChance" form:"winChance" binding:"required"`
	Odds                   float64 `json:"odds" xml:"odds" form:"odds" binding:"required"`
	BaseBet                float64 `json:"baseBet" xml:"baseBet" form:"baseBet" binding:"required"`
	RollAmount             int64   `json:"rollAmount" xml:"rollAmount" form:"rollAmount" binding:"required"`
	OnWinReturnToBaseBet   bool    `json:"onWinReturnToBaseBet" xml:"onWinReturnToBaseBet" form:"onWinReturnToBaseBet"`
	OnWinIncreaseBy        float64 `json:"onWinIncreaseBy" xml:"onWinIncreaseBy" form:"onWinIncreaseBy"`
	OnWinChangeOdds        bool    `json:"onWinChangeOdds" xml:"onWinChangeOdds" form:"onWinChangeOdds"`
	OnWinChangeOddsTo      float64 `json:"onWinChangeOddsTo" xml:"onWinChangeOddsTo" form:"onWinChangeOddsTo"`
	OnWinNewOddsWinChance  float64 `json:"onWinNewOddsWinChance" xml:"onWinNewOddsWinChance" form:"onWinNewOddsWinChance"`
	OnLossReturnToBaseBet  bool    `json:"onLossReturnToBaseBet" xml:"onLossReturnToBaseBet" form:"onLossReturnToBaseBet"`
	OnLossIncreaseBy       float64 `json:"onLossIncreaseBy" xml:"onLossIncreaseBy" form:"onLossIncreaseBy"`
	OnLossChangeOdds       bool    `json:"onLossChangeOdds" xml:"onLossChangeOdds" form:"onLossChangeOdds"`
	OnLossChangeOddsTo     float64 `json:"onLossChangeOddsTo" xml:"onLossChangeOddsTo" form:"onLossChangeOddsTo"`
	OnLossNewOddsWinChance float64 `json:"onLossNewOddsWinChance" xml:"onLossNewOddsWinChance" form:"onLossNewOddsWinChance"`
	RandomClientSeed       bool    `json:"randomClientSeed" xml:"randomClientSeed" form:"randomClientSeed"`
	AlternateHiLo          bool    `json:"alternateHiLo" xml:"alternateHiLo" form:"alternateHiLo"`
}

//HiLoVerifyRequest verify request
type HiLoVerifyRequest struct {
	ServerSeed string
	ClientSeed string
	Nonce      uint64
	Roll       float64
}

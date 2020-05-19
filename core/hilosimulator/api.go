package hilosimulator

import (
	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/utils"
	hilosimulator "github.com/Z-M-Huang/hilosimulator"
	"github.com/gin-gonic/gin"
)

//API hilo simulator
type API struct{}

//HILOSimulate /api/hilo-simulator/simulate
func (API) HILOSimulate(c *gin.Context) {
	response := core.GetResponseInContext(c.Keys)
	request := &SimulateRequest{}

	err := c.ShouldBind(&request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Invalid simulation request.",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	if request.RollAmount < 0 {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   "Roll Amount: Cannot be negative number",
		})
		core.WriteResponse(c, 400, response)
		return
	} else if request.RollAmount > 50000 {
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Requested Roll Amount is too large. Please do batches and keep the server health. Thank you",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	simConfig := &hilosimulator.Configuration{
		TotalStack: request.TotalStack,
		WinChance:  request.WinChance,
		Odds:       request.Odds,
		BaseBet:    request.BaseBet,
		RollAmount: uint64(request.RollAmount),
		OnWin: &hilosimulator.ConditionalChangeConfiguration{
			ReturnToBaseBet: request.OnWinReturnToBaseBet,
			IncreaseBet:     !request.OnWinReturnToBaseBet,
			IncreaseBetBy:   request.OnWinIncreaseBy / 100,
			ChangeOdds:      request.OnWinChangeOdds,
			NewWinChance:    request.OnWinNewOddsWinChance / 100,
		},
		OnLoss: &hilosimulator.ConditionalChangeConfiguration{
			ReturnToBaseBet: request.OnLossReturnToBaseBet,
			IncreaseBet:     !request.OnLossReturnToBaseBet,
			IncreaseBetBy:   request.OnLossIncreaseBy / 100,
			ChangeOdds:      request.OnLossChangeOdds,
			NewWinChance:    request.OnLossNewOddsWinChance / 100,
		},
	}

	simConfig.OnLoss.IncreaseBet = !simConfig.OnLoss.ReturnToBaseBet

	result, err := hilosimulator.Simulate(simConfig)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  err.Error(),
		})
		core.WriteResponse(c, 500, response)
		return
	}

	response.Data = result
	core.WriteResponse(c, 200, response)
}

//HILOVerify /api/hilo-simulator/verify
func (API) HILOVerify(c *gin.Context) {
	response := core.GetResponseInContext(c.Keys)
	request := &VerifyRequest{}

	err := c.ShouldBind(&request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Invalid simulation request.",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	valid, err := hilosimulator.Verify(request.ClientSeed, request.ServerSeed, request.Nonce, request.Roll)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  err.Error(),
		})
		core.WriteResponse(c, 400, response)
		return
	}
	response.Data = valid
	core.WriteResponse(c, 200, response)
}

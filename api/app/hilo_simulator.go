package app

import (
	"github.com/Z-M-Huang/Tools/api"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/apidata/application"
	"github.com/Z-M-Huang/Tools/data/constval"
	"github.com/Z-M-Huang/Tools/utils"
	hilosimulator "github.com/Z-M-Huang/hilosimulator"
	"github.com/gin-gonic/gin"
)

//HILOSimulate /api/hilo-simulator/simulate
func HILOSimulate(c *gin.Context) {
	response := c.Keys[constval.ResponseCtxKey].(*data.Response)
	request := &application.HiLoSimulateRequest{}

	err := c.ShouldBind(&request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid simulation request.",
		})
		api.WriteResponse(c, 200, response)
		return
	}

	if request.RollAmount < 0 {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Roll Amount: Cannot be negative number",
		})
		api.WriteResponse(c, 200, response)
		return
	} else if request.RollAmount > 50000 {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Requested Roll Amount is too large. Please do batches and keep the server health. Thank you",
		})
		api.WriteResponse(c, 200, response)
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
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}

	response.Data = result
	api.WriteResponse(c, 200, response)
}

//HILOVerify /api/hilo-simulator/verify
func HILOVerify(c *gin.Context) {
	response := c.Keys[constval.ResponseCtxKey].(*data.Response)
	request := &application.HiLoVerifyRequest{}

	err := c.ShouldBind(&request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid simulation request.",
		})
		api.WriteResponse(c, 200, response)
		return
	}

	valid, err := hilosimulator.Verify(request.ClientSeed, request.ServerSeed, request.Nonce, request.Roll)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}
	response.Data = valid
	api.WriteResponse(c, 200, response)
}

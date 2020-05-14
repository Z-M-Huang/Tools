package app

import (
	"github.com/Z-M-Huang/Tools/api"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/apidata/application"
	"github.com/Z-M-Huang/Tools/utils"
	hilosimulator "github.com/Z-M-Huang/hilosimulator"
	"github.com/gin-gonic/gin"
)

//HILOSimulate /api/hilo-simulator/simulate
func HILOSimulate(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*data.Response)
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

	simConfig := &hilosimulator.Configuration{
		OnWin:  &hilosimulator.ConditionalChangeConfiguration{},
		OnLoss: &hilosimulator.ConditionalChangeConfiguration{},
	}

	//Total Stack
	simConfig.TotalStack, err = api.ParseFloat(request.TotalStack, 64, true)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Total Stack: " + err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}

	//Win Chance
	simConfig.WinChance, err = api.ParseFloat(request.WinChance, 64, true)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Win Chance: " + err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}

	//Odds
	simConfig.Odds, err = api.ParseFloat(request.Odds, 64, true)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Odds: " + err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}

	//Base Bet
	simConfig.BaseBet, err = api.ParseFloat(request.BaseBet, 64, true)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Base Bet: " + err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}

	//Rolls Amount
	simConfig.RollAmount, err = api.ParseUint(request.RollAmount, 10, 64, true)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Roll Amount: " + err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}

	if simConfig.RollAmount > 50000 {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Requested Roll Amount is too large. Please do batches and keep the server health. Thank you",
		})
		api.WriteResponse(c, 200, response)
		return
	}

	//OnWin Return to Base Bet
	simConfig.OnWin.ReturnToBaseBet, err = api.ParseBool(request.OnWinReturnToBaseBet, false)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "On Win Return to Base Bet: " + err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}
	simConfig.OnWin.IncreaseBet = !simConfig.OnWin.ReturnToBaseBet

	if simConfig.OnWin.IncreaseBet {
		//OnWin Increate By
		simConfig.OnWin.IncreaseBetBy, err = api.ParseFloat(request.OnWinIncreaseBy, 64, true)
		if err != nil {
			response.SetAlert(&data.AlertData{
				IsWarning: true,
				Message:   "On Win Increase Bet By: " + err.Error(),
			})
			api.WriteResponse(c, 200, response)
			return
		}
		simConfig.OnWin.IncreaseBetBy = simConfig.OnWin.IncreaseBetBy / 100
	}

	//OnWin Change Odds
	simConfig.OnWin.ChangeOdds, err = api.ParseBool(request.OnWinChangeOdds, false)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "On Win Change Odds: " + err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}

	if simConfig.OnWin.ChangeOdds {
		//OnWin Change Odds To
		simConfig.OnWin.ChangeOddsTo, err = api.ParseFloat(request.OnWinChangeOddsTo, 64, true)
		if err != nil {
			response.SetAlert(&data.AlertData{
				IsWarning: true,
				Message:   "On Win Change Odds to: " + err.Error(),
			})
			api.WriteResponse(c, 200, response)
			return
		}

		//OnWin New Win Chance
		simConfig.OnWin.NewWinChance, err = api.ParseFloat(request.OnWinNewOddsWinChance, 64, true)
		if err != nil {
			response.SetAlert(&data.AlertData{
				IsWarning: true,
				Message:   "On Win New Win Chance: " + err.Error(),
			})
			api.WriteResponse(c, 200, response)
			return
		}
	}

	//OnLoss Return to Base Bet
	simConfig.OnLoss.ReturnToBaseBet, err = api.ParseBool(request.OnLossReturnToBaseBet, false)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "On Loss Return to Base Bet: " + err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}

	simConfig.OnLoss.IncreaseBet = !simConfig.OnLoss.ReturnToBaseBet

	if simConfig.OnLoss.IncreaseBet {
		//OnLoss Increate By
		simConfig.OnLoss.IncreaseBetBy, err = api.ParseFloat(request.OnLossIncreaseBy, 64, true)
		if err != nil {
			response.SetAlert(&data.AlertData{
				IsWarning: true,
				Message:   "On Loss Increase Bet By: " + err.Error(),
			})
			api.WriteResponse(c, 200, response)
			return
		}
		simConfig.OnLoss.IncreaseBetBy = simConfig.OnLoss.IncreaseBetBy / 100
	}

	//OnLoss Change Odds
	simConfig.OnLoss.ChangeOdds, err = api.ParseBool(request.OnLossChangeOdds, false)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "On Loss Change Odds: " + err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}

	if simConfig.OnLoss.ChangeOdds {
		//OnLoss Change Odds To
		simConfig.OnLoss.ChangeOddsTo, err = api.ParseFloat(request.OnLossChangeOddsTo, 64, true)
		if err != nil {
			response.SetAlert(&data.AlertData{
				IsWarning: true,
				Message:   "On Loss Change Odds to: " + err.Error(),
			})
			api.WriteResponse(c, 200, response)
			return
		}

		//OnLoss New Win Chance
		simConfig.OnLoss.NewWinChance, err = api.ParseFloat(request.OnLossNewOddsWinChance, 64, true)
		if err != nil {
			response.SetAlert(&data.AlertData{
				IsWarning: true,
				Message:   "On Loss New Win Chance: " + err.Error(),
			})
			api.WriteResponse(c, 200, response)
			return
		}
	}

	simConfig.RandomClientSeed, err = api.ParseBool(request.RandomClientSeed, false)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Random Client Seed: " + err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}

	simConfig.AlternateHiLo, err = api.ParseBool(request.AlternateHiLo, false)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Alternate Bet Hi/Low: " + err.Error(),
		})
		api.WriteResponse(c, 200, response)
		return
	}

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
	response := c.Keys[utils.ResponseCtxKey].(*data.Response)
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

package hilosimulator

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	hilosimulator "github.com/Z-M-Huang/hilosimulator"
	"github.com/gin-gonic/gin"
)

//API hilo simulator
type API struct{}

func simulate(request *SimulateRequest) (int, *data.APIResponse) {
	response := &data.APIResponse{}
	if request.RollAmount < 0 {
		response.Message = "Roll Amount: Cannot be negative number"
		return http.StatusBadRequest, response
	} else if request.RollAmount > 50000 {
		response.Message = "Requested Roll Amount is too large. Please do batches and keep the server health. Thank you"
		return http.StatusBadRequest, response
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
		response.Message = err.Error()
		return http.StatusBadRequest, response
	}

	response.Data = result

	return http.StatusOK, response
}

func verify(request *VerifyRequest) (int, *data.APIResponse) {
	response := &data.APIResponse{}
	valid, err := hilosimulator.Verify(request.ClientSeed, request.ServerSeed, request.Nonce, request.Roll)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Message = err.Error()
		response.Data = false
		return http.StatusBadRequest, response
	}
	response.Data = valid
	return http.StatusOK, response
}

//HILOSimulate /api/hilo-simulator/simulate
func (API) HILOSimulate(c *gin.Context) {
	request := &SimulateRequest{}
	err := c.ShouldBind(&request)
	if err != nil {
		utils.Logger.Error(err.Error())
		core.WriteResponse(c, 400, &data.APIResponse{
			Message: "Invalid simulation request.",
		})
		return
	}

	status, response := simulate(request)

	core.WriteResponse(c, status, response)
}

//HILOVerify /api/hilo-simulator/verify
func (API) HILOVerify(c *gin.Context) {
	request := &VerifyRequest{}
	err := c.ShouldBind(&request)
	if err != nil {
		utils.Logger.Error(err.Error())
		core.WriteResponse(c, 400, &data.APIResponse{
			Message: "Invalid verify request.",
		})
		return
	}

	status, response := verify(request)

	core.WriteResponse(c, status, response)
}

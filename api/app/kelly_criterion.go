package app

import (
	"math"
	"strconv"

	"github.com/Z-M-Huang/Tools/api"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/apidata/application"
	"github.com/Z-M-Huang/Tools/utils"
	kellycriterion "github.com/Z-M-Huang/kelly-criterion"
	"github.com/gin-gonic/gin"
)

//KellyCriterionSimulate /api/kelly-criterion/simulate
func KellyCriterionSimulate(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*data.Response)
	var simulationResult []*application.KellyCriterionSimulationResponse
	request := &application.KellyCriterionRequest{}
	err := c.ShouldBind(&request)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid simulation request.",
		})
		api.WriteResponse(c, 200, response)
		return
	}

	maxPayout, err := strconv.ParseFloat(request.MaxWinChancePayout, 64)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Max Payout needs to be a number.",
		})
		api.WriteResponse(c, 200, response)
		return
	}

	maxChance, err := strconv.ParseFloat(request.MaxWinChance, 64)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Max Chance needs to be a number.",
		})
		api.WriteResponse(c, 200, response)
		return
	}

	total := maxChance * maxPayout / 100
	for i := 0; i < 1000; i++ {
		payout := float64(maxPayout) + (float64(i) * 0.01)
		chance := float64(total / payout)
		simulationResult = append(simulationResult, &application.KellyCriterionSimulationResponse{
			Payout: math.Round(payout*100) / 100,
			Chance: math.Round(chance*100) / 100,
			Factor: math.Round(kellycriterion.Calculate(chance, payout)*10000) / 10000,
		})
	}
	response.Data = simulationResult
	api.WriteResponse(c, 200, response)
}

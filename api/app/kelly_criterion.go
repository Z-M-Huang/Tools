package app

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/Z-M-Huang/Tools/api"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/apidata/application"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/julienschmidt/httprouter"
)

//KellyCriterionSimulate /api/kelly-criterion/simulate
func KellyCriterionSimulate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	var simulationResult []*application.KellyCriterionSimulationResponse
	request := &application.KellyCriterionRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Alert.IsDanger = true
		response.Alert.Message = "Invalid simulation request."
		api.WriteResponse(w, response)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Alert.IsDanger = true
		response.Alert.Message = "Invalid simulation request."
		api.WriteResponse(w, response)
		return
	}

	maxPayout, err := strconv.ParseFloat(request.MaxWinChancePayout, 64)
	if err != nil {
		response.Alert.IsDanger = true
		response.Alert.Message = "Invalid simulation request."
		api.WriteResponse(w, response)
		return
	}

	maxChance, err := strconv.ParseFloat(request.MaxWinChance, 64)
	if err != nil {
		response.Alert.IsDanger = true
		response.Alert.Message = "Invalid simulation request."
		api.WriteResponse(w, response)
		return
	}

	total := maxChance * maxPayout / 100
	for i := 0; i < 1000; i++ {
		payout := float64(maxPayout) + (float64(i) * float64(0.01))
		chance := float64(total / payout)
		simulationResult = append(simulationResult, &application.KellyCriterionSimulationResponse{
			Payout: math.Round(payout*100) / 100,
			Chance: math.Round(chance*100) / 100,
			Factor: math.Round(((chance*payout-(1-chance))/payout)*10000) / 10000,
		})
	}
	response.Data = simulationResult
	api.WriteResponse(w, response)
}

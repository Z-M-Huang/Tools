package kellycriterion

import (
	"math"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/utils"
	kellycriterion "github.com/Z-M-Huang/kelly-criterion"
	"github.com/gin-gonic/gin"
)

//API kelly criterion
type API struct{}

//Simulate /api/kelly-criterion/simulate
func (API) Simulate(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*core.Response)
	var simulationResult []*Response
	request := &Request{}
	err := c.ShouldBind(&request)
	if err != nil {
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Invalid simulation request.",
		})
		core.WriteResponse(c, 200, response)
		return
	}

	total := request.MaxWinChance * request.MaxWinChancePayout / 100
	for i := 0; i < 1000; i++ {
		payout := float64(request.MaxWinChancePayout) + (float64(i) * 0.01)
		chance := float64(total / payout)
		simulationResult = append(simulationResult, &Response{
			Payout: math.Round(payout*100) / 100,
			Chance: math.Round(chance*100) / 100,
			Factor: math.Round(kellycriterion.Calculate(chance, payout)*10000) / 10000,
		})
	}
	response.Data = simulationResult
	core.WriteResponse(c, 200, response)
}

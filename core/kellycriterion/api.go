package kellycriterion

import (
	"math"
	"net/http"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	kellycriterion "github.com/Z-M-Huang/kelly-criterion"
	"github.com/gin-gonic/gin"
)

//API kelly criterion
type API struct{}

func simualte(request *Request) (int, *data.APIResponse) {
	response := &data.APIResponse{}
	var simulationResult []*Response
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
	return http.StatusOK, response
}

//Simulate /api/kelly-criterion/simulate
func (API) Simulate(c *gin.Context) {
	request := &Request{}
	err := c.ShouldBind(&request)
	if err != nil {
		core.WriteResponse(c, http.StatusBadRequest, &data.APIResponse{
			Message: "Invalid simulation request.",
		})
		return
	}

	status, response := simualte(request)

	core.WriteResponse(c, status, response)
}

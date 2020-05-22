package hilosimulator

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimulate(t *testing.T) {
	requests := []*SimulateRequest{
		{
			TotalStack: 100000,
			BaseBet:    1,
			WinChance:  47.5,
			Odds:       2,
			RollAmount: 500,
		},
		{
			TotalStack:            100000,
			BaseBet:               1,
			WinChance:             47.5,
			Odds:                  2,
			RollAmount:            500,
			OnWinReturnToBaseBet:  true,
			OnLossIncreaseBy:      100,
			OnLossReturnToBaseBet: false,
		},
	}

	for _, r := range requests {
		status, response := simulate(r)

		assert.Equal(t, http.StatusOK, status)
		assert.NotEmpty(t, response.Data)
	}
}

func TestSimulateFail(t *testing.T) {
	requests := []*SimulateRequest{
		{
			TotalStack: 100000,
			BaseBet:    1,
			WinChance:  47.5,
			Odds:       2,
			RollAmount: -1,
		},
		{
			TotalStack:            100000,
			BaseBet:               1,
			WinChance:             47.5,
			Odds:                  2,
			RollAmount:            60000,
			OnWinReturnToBaseBet:  true,
			OnLossIncreaseBy:      100,
			OnLossReturnToBaseBet: false,
		},
	}

	for _, r := range requests {
		status, response := simulate(r)

		assert.Equal(t, http.StatusBadRequest, status)
		assert.Empty(t, response.Data)
		assert.NotEmpty(t, response.Message)
	}
}

func TestVerify(t *testing.T) {
	request := &VerifyRequest{
		ServerSeed: "74767e76b301c53b3388c0d32f35acb86a263b33ea341c004000db5649c7fd67",
		ClientSeed: "42fe4b33e7f3bdd03d379ea5e239ed089b05f9d089e5063b9150bd598fc46e5c",
		Nonce:      0,
		Roll:       45.53,
	}

	status, response := verify(request)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, true, response.Data)
}

func TestVerifyFail(t *testing.T) {
	request := &VerifyRequest{
		ServerSeed: "a",
		ClientSeed: "cd",
		Nonce:      0,
		Roll:       45.53,
	}

	status, response := verify(request)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, false, response.Data)
}

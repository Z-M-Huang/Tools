package kellycriterion

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimulate(t *testing.T) {
	request := &Request{
		MaxWinChancePayout: 1,
		MaxWinChance:       95,
	}

	status, response := simualte(request)

	assert.Equal(t, status, http.StatusOK)
	assert.NotEmpty(t, response.Data)
}

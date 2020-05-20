package dnslookup

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	request := &Request{
		DomainName: "google.com",
	}

	status, response := lookup(request)

	assert.Equal(t, http.StatusOK, status)
	assert.NotEmpty(t, response)
	assert.NotEmpty(t, response.Data)
}

func TestLookupFail(t *testing.T) {
	request := &Request{
		DomainName: "google",
	}

	status, response := lookup(request)

	assert.NotEqual(t, http.StatusOK, status)
	assert.NotEmpty(t, response)
}

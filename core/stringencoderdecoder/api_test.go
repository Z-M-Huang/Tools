package stringencoderdecoder

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	testStr := "https://encoding test"
	requests := []*Request{
		{
			RequestString: testStr,
			Type:          "Base32",
			Action:        "encode",
		},
		{
			RequestString: testStr,
			Type:          "Base64",
			Action:        "encode",
		},
		{
			RequestString: testStr,
			Type:          "Binary",
			Action:        "encode",
		},
		{
			RequestString: testStr,
			Type:          "URL",
			Action:        "encode",
		},
	}

	for _, r := range requests {
		status, response := encodeDecode(r)
		assert.Equal(t, http.StatusOK, status)
		assert.NotEmpty(t, response)
		assert.NotEmpty(t, response.Data)
		assert.Empty(t, response.Message)
	}
}

func TestDecode(t *testing.T) {
	testStr := "https://encoding test"
	requests := []*Request{
		{
			RequestString: "NB2HI4DTHIXS6ZLOMNXWI2LOM4QHIZLTOQ======",
			Type:          "Base32",
			Action:        "decode",
		},
		{
			RequestString: "aHR0cHM6Ly9lbmNvZGluZyB0ZXN0",
			Type:          "Base64",
			Action:        "decode",
		},
		{
			RequestString: "1101000 1110100 1110100 1110000 1110011 111010 101111 101111 1100101 1101110 1100011 1101111 1100100 1101001 1101110 1100111 100000 1110100 1100101 1110011 1110100",
			Type:          "Binary",
			Action:        "decode",
		},
		{
			RequestString: "https%3A%2F%2Fencoding+test",
			Type:          "URL",
			Action:        "decode",
		},
	}

	for _, r := range requests {
		status, response := encodeDecode(r)
		assert.Equal(t, http.StatusOK, status)
		assert.NotEmpty(t, response)
		assert.NotEmpty(t, response.Data)
		assert.Empty(t, response.Message)
		assert.Equal(t, testStr, response.Data.([]string)[0])
	}
}

func TestDecodeFailure(t *testing.T) {
	requests := []*Request{
		{
			RequestString: "a",
			Type:          "Base32",
			Action:        "decode",
		},
		{
			RequestString: "a",
			Type:          "Base64",
			Action:        "decode",
		},
		{
			RequestString: "%",
			Type:          "URL",
			Action:        "decode",
		},
	}

	for _, r := range requests {
		status, response := encodeDecode(r)
		assert.Equal(t, http.StatusBadRequest, status)
		assert.NotEmpty(t, response)
		assert.NotEmpty(t, response.Data)
		assert.NotEmpty(t, response.Message)
		assert.Equal(t, response.Data, response.Message)
	}
}

func TestInvalidActionAndEncoding(t *testing.T) {
	requests := []*Request{
		{
			RequestString: "a",
			Type:          "Base33",
			Action:        "decode",
		},
		{
			RequestString: "a",
			Type:          "Base64",
			Action:        "test",
		},
	}

	for _, r := range requests {
		status, response := encodeDecode(r)
		assert.Equal(t, http.StatusBadRequest, status)
		assert.NotEmpty(t, response)
		assert.NotEmpty(t, response.Data)
		assert.NotEmpty(t, response.Message)
		assert.Equal(t, response.Data, response.Message)
	}
}

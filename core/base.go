package core

import (
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
)

//GetResponseInContext get response struct from context
func GetResponseInContext(contextKey map[string]interface{}) *data.Response {
	response := contextKey[utils.ResponseCtxKey]
	if response == nil {
		return nil
	}
	return response.(*data.Response)
}

//GetClaimInContext get claim struct from context
func GetClaimInContext(contextKey map[string]interface{}) *data.JWTClaim {
	claim := contextKey[utils.ClaimCtxKey]
	if claim == nil {
		return nil
	}
	return claim.(*data.JWTClaim)
}

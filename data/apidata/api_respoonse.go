package apidata

//APIResponse base api response
type APIResponse struct {
	Alert AlertResponse
	Data  interface{}
}

//AlertResponse alert response
type AlertResponse struct {
	IsInfo    bool
	IsSuccess bool
	IsWarning bool
	IsAlert   bool
	Message   string
}

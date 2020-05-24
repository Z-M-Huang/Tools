package portchecker

//Request port checker request
type Request struct {
	Host     string `json:"host" xml:"host" form:"host" binding:"required"`
	Port     int    `json:"port" xml:"port" form:"port" binding:"required"`
	PortType string `json:"type" xml:"type" form:"type" binding:"required"`
}

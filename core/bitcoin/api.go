package bitcoin

import (
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/btcsuite/btcd/rpcclient"
)

//API bitcoin
type API struct{}

var client *rpcclient.Client

func init() {
	connCfg := &rpcclient.ConnConfig{
		Host:         data.BitcoinRPCConfig.Host,
		User:         data.BitcoinRPCConfig.User,
		Pass:         data.BitcoinRPCConfig.Password,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	newClient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	client = newClient
}

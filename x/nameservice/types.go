package nameservice

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "nameservice"
	TypeMsgSetName = "set_name"
	TypeMsgBuyName = "buy_name"
)

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string              `json:"value"`
	Owner sdktypes.AccAddress `json:"owner"`
	Price sdktypes.Coins      `json:"price"`
}

// MinNamePrice Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdktypes.Coins{sdktypes.NewInt64Coin("nametoken", 1)}

// NewWhois Returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

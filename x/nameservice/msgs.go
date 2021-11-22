package nameservice

import (
	"encoding/json"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSetName defines a SetName message
type MsgSetName struct {
	Name  string
	Value string
	Owner sdktypes.AccAddress
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgSetName(name string, value string, owner sdktypes.AccAddress) MsgSetName {
	return MsgSetName{
		Name:  name,
		Value: value,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgSetName) Route() string {
	return ModuleName
}

// Type should return the action
func (msg MsgSetName) Type() string {
	return TypeMsgSetName
}

func (msg MsgSetName) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Value cannot be empty")
	}

	return nil
}

func (msg MsgSetName) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdktypes.MustSortJSON(b)
}

func (msg MsgSetName) GetSigners() []sdktypes.AccAddress {
	return []sdktypes.AccAddress{msg.Owner}
}

type MsgBuyName struct {
	Name  string
	Bid   sdktypes.Coins
	Buyer sdktypes.AccAddress
}

func NewMsgBuyName(name string, bid sdktypes.Coins, buyer sdktypes.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}

// Route should return the name of the module
func (msg MsgBuyName) Route() string {
	return ModuleName
}

// Type should return the action
func (msg MsgBuyName) Type() string {
	return TypeMsgBuyName
}

func (msg MsgBuyName) ValidateBasic() error {
	if msg.Buyer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Bids must be positive")
	}
	return nil
}

func (msg MsgBuyName) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdktypes.MustSortJSON(b)
}

func (msg MsgBuyName) GetSigners() []sdktypes.AccAddress {
	return []sdktypes.AccAddress{msg.Buyer}
}

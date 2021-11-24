package nameservice

import (
	"fmt"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(keeper Keeper) sdktypes.Handler {
	return func(ctx sdktypes.Context, msg sdktypes.Msg) (*sdktypes.Result, error) {
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBugName(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgSetName(ctx sdktypes.Context, keeper Keeper, msg MsgSetName) (*sdktypes.Result, error) {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	keeper.SetName(ctx, msg.Name, msg.Value)
	return &sdktypes.Result{}, nil
}

func handleMsgBugName(ctx sdktypes.Context, keeper Keeper, msg MsgBuyName) (*sdktypes.Result, error) {
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Bid not high enough")
	}

	if keeper.HasOwner(ctx, msg.Name) {
		if err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid); err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Buyer does not have enough coins")
		}
	} else {
		if _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid); err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Buyer does not have enough coins")
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return &sdktypes.Result{}, nil
}

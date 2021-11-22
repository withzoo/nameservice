package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdktypes.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdktypes.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// GetNamesIterator Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetNamesIterator(ctx sdktypes.Context) sdktypes.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdktypes.KVStorePrefixIterator(store, []byte{})
}

func (k Keeper) SetWhois(ctx sdktypes.Context, name string, whois Whois) {
	if whois.Owner.Empty() {
		return
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
}

func (k Keeper) GetWhois(ctx sdktypes.Context, name string) Whois {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(name)) {
		return NewWhois()
	}

	bz := store.Get([]byte(name))
	var whois Whois
	k.cdc.MustUnmarshalBinaryBare(bz, &whois)
	return whois
}

func (k Keeper) ResolveName(ctx sdktypes.Context, name string) string {
	return k.GetWhois(ctx, name).Value
}

func (k Keeper) SetName(ctx sdktypes.Context, name string, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

func (k Keeper) HasOwner(ctx sdktypes.Context, name string) bool {
	return !k.GetWhois(ctx, name).Owner.Empty()
}

func (k Keeper) GetOwner(ctx sdktypes.Context, name string) sdktypes.AccAddress {
	return k.GetWhois(ctx, name).Owner
}

func (k Keeper) SetOwner(ctx sdktypes.Context, name string, owner sdktypes.AccAddress) {
	whois := k.GetWhois(ctx, name)
	whois.Owner = owner
	k.SetWhois(ctx, name, whois)
}

func (k Keeper) GetPrice(ctx sdktypes.Context, name string) sdktypes.Coins {
	return k.GetWhois(ctx, name).Price
}

func (k Keeper) SetPrice(ctx sdktypes.Context, name string, price sdktypes.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.Price = price
	k.SetWhois(ctx, name, whois)
}

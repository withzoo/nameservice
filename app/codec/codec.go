package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func MakeCodec(bm module.BasicManager) *codec.Codec {
	cdc := codec.New()
	return cdc
}

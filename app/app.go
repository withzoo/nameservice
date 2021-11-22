package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"github.com/withzoo/nameservice/app/codec"
)

const appName = "nameservice"

var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
	)
)

type nameServiceApp struct {
	*baseapp.BaseApp
}

func NewNameServiceApp(logger log.Logger, db dbm.DB) *nameServiceApp {
	cdc := codec.MakeCodec(ModuleBasics)
	app := baseapp.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	return &nameServiceApp{
		BaseApp: app,
		cdc:     cdc,
	}
}

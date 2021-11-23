package main

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	auth "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	bank "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	appcdc "github.com/withzoo/nameservice/app/codec"
	nsclient "github.com/withzoo/nameservice/x/nameservice/client"
	nsrest "github.com/withzoo/nameservice/x/nameservice/client/rest"
)

const (
	storeAcc = "acc"
	storeNS  = "nameservice"
)

var defaultCLIHome = os.ExpandEnv("$HOME/.nscli")

func main() {
	cobra.EnableCommandSorting = false

	cdc := appcdc.MakeCodec()

	// Read in the configuration file for the sdk
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount(sdktypes.Bech32PrefixAccAddr, sdktypes.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdktypes.Bech32PrefixValAddr, sdktypes.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdktypes.Bech32PrefixConsAddr, sdktypes.Bech32PrefixConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "nscli",
		Short: "nameservice Client",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(defaultCLIHome),
		queryCmd(cdc, nsclient.NewModuleClient(storeNS, cdc)),
		txCmd(cdc, nsclient.NewModuleClient(storeNS, cdc)),
		flags.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
	)

	executor := cli.PrepareMainCmd(rootCmd, "NS", defaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func registerRoutes(rs *lcd.RestServer) {
	rpc.RegisterRPCRoutes(rs.CliCtx, rs.Mux)
	auth.RegisterRoutes(rs.CliCtx, rs.Mux, storeAcc)
	bank.RegisterRoutes(rs.CliCtx, rs.Mux)
	nsrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.CliCtx.Codec, storeNS)
}

func queryCmd(cdc *amino.Codec, mc nsclient.ModuleClient) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxCmd(cdc),
		flags.LineBreak,
		authcmd.GetAccountCmd(cdc),
		mc.GetQueryCmd(),
	)

	return queryCmd
}

func txCmd(cdc *amino.Codec, mc nsclient.ModuleClient) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		flags.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetBroadcastCommand(cdc),
		flags.LineBreak,
		mc.GetTxCmd(),
	)

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}

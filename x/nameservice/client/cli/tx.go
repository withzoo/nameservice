package cli

import (
	"bufio"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"

	"github.com/withzoo/nameservice/x/nameservice"
)

// GetCmdBuyName is the CLI command for sending a BuyName transaction
func GetCmdBuyName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "buy-name [name] [amount]",
		Short: "bid for existing name or claim new name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			//cliCtx := context.NewCLIContext().WithCodec(cdc)
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			coins, err := sdktypes.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := nameservice.NewMsgBuyName(args[0], coins, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdktypes.Msg{msg})
		},
	}
}

// GetCmdSetName is the CLI command for sending a SetName transaction
func GetCmdSetName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-name [name] [value]",
		Short: "set the value associated with a name that you own",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			//cliCtx := context.NewCLIContext().WithCodec(cdc)
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := nameservice.NewMsgSetName(args[0], args[1], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdktypes.Msg{msg})
		},
	}
}

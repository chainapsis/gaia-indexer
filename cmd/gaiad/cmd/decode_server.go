package cmd

// DONTCOVER

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/gaia/v10/app/params"
	"github.com/cosmos/gaia/v10/cmd/gaiad/cmd/decoder"
	"github.com/spf13/cobra"
)

const (
	decodeServerPort = "decodeServerPort"
)

// decoderCmd gets cmd to run decode server
func decoderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decoder",
		Short: "Example osmosisd decode -p 8888, which would run decoder server on specified port",
		Long:  "decoder command runs decoder server to decode byte array to General Cosmos messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			decodeServerFlag, err := cmd.Flags().GetString(decodeServerPort)
			if err != nil {
				return err
			}

			amino := codec.NewLegacyAmino()
			interfaceRegistry := types.NewInterfaceRegistry()
			marshaler := codec.NewProtoCodec(interfaceRegistry)
			txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

			d := decoder.Decoder{
				EncodingConfig: params.EncodingConfig{
					InterfaceRegistry: interfaceRegistry,
					Codec:             marshaler,
					TxConfig:          txCfg,
					Amino:             amino,
				},
			}
			err = d.ListenAndServe(decodeServerFlag)
			if err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringP(decodeServerPort, "p", "", "port to listen to")
	cmd.MarkFlagRequired(decodeServerPort)
	return cmd
}

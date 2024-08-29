package gaia

import (
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/gaia/v19/app/params"

	db "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

var emptyWasmOption []wasmkeeper.Option

func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	app := NewGaiaApp(
		log.NewNopLogger(),
		db.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,
		EmptyAppOptions{},
		emptyWasmOption,
	)

	moduleBasics := app.ModuleBasics

	moduleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	moduleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

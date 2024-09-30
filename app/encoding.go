package gaia

import (
	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	db "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/gaia/v20/app/params"
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

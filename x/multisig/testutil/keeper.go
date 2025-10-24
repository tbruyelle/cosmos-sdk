package testutil

import (
	"testing"

	"github.com/golang/mock/gomock"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/multisig/keeper"
	"github.com/cosmos/cosmos-sdk/x/multisig/types"
)

type Mocks struct {
	Router *MockRouter
}

func SetupMsgServer(t *testing.T) (types.MsgServer, *keeper.Keeper, Mocks, sdk.Context) {
	t.Helper()
	k, m, ctx := SetupMultisigKeeper(t)
	return keeper.NewMsgServerImpl(*k), k, m, ctx
}

func SetupMultisigKeeper(t *testing.T) (
	*keeper.Keeper,
	Mocks,
	sdk.Context,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	m := Mocks{
		Router: NewMockRouter(ctrl),
	}

	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	banktypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	return keeper.NewKeeper(encCfg.Codec, storeService, m.Router, authority), m, ctx
}

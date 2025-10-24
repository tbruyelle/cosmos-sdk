package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/x/multisig/testutil"
	"github.com/cosmos/cosmos-sdk/x/multisig/types"
)

func TestParams(t *testing.T) {
	k, _, ctx := testutil.SetupMultisigKeeper(t)
	params := types.DefaultParams()

	err := k.Params.Set(ctx, params)
	require.NoError(t, err)

	params2, err := k.Params.Get(ctx)
	require.NoError(t, err)
	require.EqualValues(t, params, params2)
}

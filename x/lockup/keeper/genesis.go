package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/althea-net/althea-chain/x/lockup/types"
)

// InitGenesis starts a chain from a genesis state
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetChainLocked(ctx, data.GetLocked())
	k.SetLockExemptAddresses(ctx, data.GetLockExempt())
	k.SetLockedMessageTypes(ctx, data.GetLockedMessageTypes())
}

// ExportGenesis exports all the state needed to restart the chain
// from the current state of the chain
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	return types.GenesisState{
		Locked:             k.GetChainLocked(ctx),
		LockExempt:         k.GetLockExemptAddresses(ctx),
		LockedMessageTypes: k.GetLockedMessageTypes(ctx),
	}
}

package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/althea-net/althea-chain/x/lockup/types"
)

type Keeper struct {
	storeKey   sdk.StoreKey
	paramSpace paramstypes.Subspace
	cdc        codec.BinaryMarshaler
}

func NewKeeper(cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, paramSpace paramstypes.Subspace) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	k := Keeper{
		cdc:        cdc,
		paramSpace: paramSpace,
		storeKey:   storeKey,
	}

	return k
}

func (k Keeper) GetChainIsLocked(ctx sdk.Context) bool {
	var locked bool
	k.paramSpace.Get(ctx, types.LockedKey, &locked)
	return locked
}

func (k Keeper) SetChainIsLocked(ctx sdk.Context, locked bool) {
	k.paramSpace.Set(ctx, types.LockedKey, &locked)
}

func (k Keeper) GetLockExemptAddresses(ctx sdk.Context) []string {
	var lockExempt []string
	k.paramSpace.Get(ctx, types.LockExemptKey, &lockExempt)
	return lockExempt
}

func (k Keeper) SetLockExemptAddresses(ctx sdk.Context, lockExempt []string) {
	k.paramSpace.Set(ctx, types.LockExemptKey, &lockExempt)
}

func (k Keeper) GetLockedMessageTypes(ctx sdk.Context) []string {
	var lockedMessageTypes []string
	k.paramSpace.Get(ctx, types.LockedMessageTypesKey, &lockedMessageTypes)
	return lockedMessageTypes
}

func (k Keeper) SetLockedMessageTypes(ctx sdk.Context, lockedMessageTypes []string) {
	k.paramSpace.Get(ctx, types.LockedMessageTypesKey, &lockedMessageTypes)
}

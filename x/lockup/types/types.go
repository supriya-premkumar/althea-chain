package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "lockup"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	// LockedKey Indexes the Locked boolean, indicating that the lockup module is active or inactive
	// In other words Locked -> The chain is "locked up"
	LockedKey = []byte{0x1}

	// LockExemptKey Indexes the LockExempt addresses, who will be able to initiate transactions even
	// when the chain is locked
	LockExemptKey = []byte{0x2}

	// LockedMessageTypesKey Indexes the LockedMessageTypes array, the collection of messages which
	// will be blocked when the chain is locked up and not sent from a LockExempt address
	LockedMessageTypesKey = []byte{0x3}
)

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&GenesisState{})
}

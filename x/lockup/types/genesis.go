package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/althea-net/althea-chain/types"
)

// DefaultGenesisState creates a simple GenesisState suitible for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Locked: true,
		LockExempt: []string{
			"0x0000000000000000000000000000000000000000",
		},
		LockedMessageTypes: []string{
			banktypes.TypeMsgSend,
			banktypes.TypeMsgMultiSend,
		},
	}
}

func (s GenesisState) ValidateBasic() error {
	if err := ValidateLockExempt(s.LockExempt); err != nil {
		return sdkerrors.Wrap(err, "Invalid LockExempt GenesisState")
	}
	if err := ValidateLockedMessageTypes(s.LockedMessageTypes); err != nil {
		return sdkerrors.Wrap(err, "Invalid LockedMessageTypes GenesisState")
	}
	return nil
}

func ValidateLocked(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid locked type: %T", i)
	}
	return nil
}

func ValidateLockExempt(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid lock exempt type: %T", i)
	}
	if len(v) == 0 {
		return fmt.Errorf("no lock exempt addresses")
	}
	return nil
}

func ValidateLockedMessageTypes(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid locked message types type: %T", i)
	}
	if len(v) == 0 {
		return fmt.Errorf("no locked message types %v", s)
	}
	return nil
}

func (s *GenesisState) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(types.LockedKey, &s.Locked, ValidateLocked),
		paramtypes.NewParamSetPair(types.LockExemptKey, &s.LockExempt, ValidateLockExempt),
		paramtypes.NewParamSetPair(types.LockedMessageTypesKey, &s.LockedMessageTypes, ValidateLockedMessageTypes),
	}
}

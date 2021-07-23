package lockup

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/althea-net/althea-chain/x/lockup/keeper"
	"github.com/althea-net/althea-chain/x/lockup/types"
)

// WrappedAnteHandle A AnteDecorator used to wrap any AnteHandler for decorator chaining
type WrappedAnteHandler struct {
	anteHandler sdk.AnteHandler
}

// AnteHandle calls wad.anteHandler and then the next one in the chain
// This is necessary to use the Cosmos SDK's NewAnteHandler() output with a LockupAnteHandler
func (wad WrappedAnteHandler) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	modCtx, ok := wad.anteHandler(ctx, tx, simulate)
	if ok != nil {
		return modCtx, err
	}
	return next(modCtx, tx, simulate)
}

// LockupWrappedAnteHandler wraps the input AnteHandler arount a LockupAnteHandler
func WrappedLockupAnteHandler(
	anteHandler sdk.AnteHandler,
	lockupKeeper keeper.Keeper,
) sdk.AnteHandler {
	wrapped := WrappedAnteHandler{anteHandler}
	lad := NewLockupAnteDecorator(lockupKeeper)

	return sdk.ChainAnteDecorators(wrapped, lad)
}

// LockupAnteDecorator Ensures that any transaction under a locked chain originates from a LockExempt address
type LockupAnteDecorator struct {
	lockupKeeper      keeper.Keeper
	exemptSet         map[string]struct{}
	lockedMsgTypesSet map[string]struct{}
}

// AnteHandle Ensures that any transaction under a locked chain originates from a LockExempt address
func (lad LockupAnteDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if lad.lockupKeeper.ChainIsLocked(ctx) {
		for _, msg := range tx.GetMsgs() {
			if lad.lockedMsgTypesSet == nil {
				lad.initLockedMsgTypesSet(ctx)
			}
			if _, typePresent := lad.lockedMsgTypesSet[msg.Type()]; typePresent {
				if lad.exemptSet == nil {
					lad.initExemptSet(ctx)
				}
				if allow, err := allowMessage(msg, lad.exemptSet); !allow {
					return ctx, sdkerrors.Wrap(err,
						fmt.Sprintf("Transaction %v blocked because of message %v", tx, msg))
				}
			}
		}
	}
	return next(ctx, tx, simulate)
}

// initLockedMsgTypesSet initializes the lockedMsgTypesSet
func (lad LockupAnteDecorator) initLockedMsgTypesSet(ctx sdk.Context) {
	lockedMsgTypes := lad.lockupKeeper.GetLockedMessageTypes(ctx)
	lad.lockedMsgTypesSet = createSet(lockedMsgTypes)
}

func (lad LockupAnteDecorator) initExemptSet(ctx sdk.Context) {
	exemptAddresses := lad.lockupKeeper.GetLockExemptAddresses(ctx)
	lad.exemptSet = createSet(exemptAddresses)
}

// NewAnteHandler returns an AnteHandler that ensures any transaction under a locked chain
// originates from a LockExempt address
func NewLockupAnteHandler(lockupKeeper keeper.Keeper) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(NewLockupAnteDecorator(lockupKeeper))
}

// NewLockupAnteDecorator initializes a LockupAnteDecorator for locking messages
// based on the settings stored in lockupKeeper
// Note: Cannot init exemptSet nor lockedMsgTypesSet as the context may not be available yet
func NewLockupAnteDecorator(lockupKeeper keeper.Keeper) LockupAnteDecorator {
	return LockupAnteDecorator{lockupKeeper, nil, nil}
}

func createSet(strings []string) map[string]struct{} {
	type void struct{}
	var member void
	set := make(map[string]struct{})

	for _, str := range strings {
		if _, present := set[str]; present {
			continue
		}
		set[str] = member
	}

	return set
}

// allowMessage checks that an input `msg` was sent by only addresses in `exemptSet`
// returns true if `msg` is either permissible or not a type of message this module blocks
func allowMessage(msg sdk.Msg, exemptSet map[string]struct{}) (bool, error) {
	switch msg.Type() {
	case banktypes.TypeMsgSend:
		msgSend := msg.(*banktypes.MsgSend)
		if _, present := exemptSet[msgSend.FromAddress]; !present {
			// Message sent from a non-exempt address while the chain is locked up, returning error
			return false, sdkerrors.Wrap(types.ErrLocked,
				"The chain is locked, only exempt addresses may be the FromAddress in a Send message")
		}
		return true, nil
	case banktypes.TypeMsgMultiSend:
		msgMultiSend := msg.(*banktypes.MsgMultiSend)
		for _, input := range msgMultiSend.Inputs {
			if _, present := exemptSet[input.Address]; !present {
				// Multi-send Message sent with a non-exempt input address while the chain is locked up, returning error
				return false, sdkerrors.Wrap(types.ErrLocked,
					"The chain is locked, only exempt addresses may be inputs in a MultiSend message")
			}
		}
		return true, nil
	default:
		return true, nil
	}
}

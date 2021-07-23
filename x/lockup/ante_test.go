package lockup

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/althea-net/althea-chain/x/lockup/keeper"
	"github.com/althea-net/althea-chain/x/lockup/types"
)

func TestAnteHandler(t *testing.T) {
	// TODO: Create a new ante decorator with default state and test that it
	input := keeper.CreateTestEnv(t)
	ctx := input.Context
	appCodec := input.Marshaler
	txCfg := tx.NewTxConfig(appCodec, tx.DefaultSignModes)
	txBuilder = txCfg.NewTxBuilder()
	keys := sdk.NewKVStoreKeys(types.StoreKey)
	subspace, _ := input.ParamsKeeper.GetSubspace(types.ModuleName)
	keeper := keeper.NewKeeper(
		appCodec, keys[types.StoreKey], subspace,
	)

	handler := NewLockupAnteHandler(keeper)
	allowedMsgSend := GetAllowedMsgSendTx()
	// blocks a transaction coming from 0x1 but not one coming from 0x0.

	// TODO: Test that messages not of the right type aren't blocked

}

// TODO: Test that empty lock exempt fails validation

// TODO: Test that empty locked message types fails validation

// TODO: Test that nothing is blocked when locked is false

func GetAllowedMsgSendTx(txCfg client.TxConfig) sdk.Tx {
	fromAddr, _ := sdk.AccAddressFromHex("0x0000000000000000000000000000000000000000")
	toAddr, _ := sdk.AccAddressFromHex("0x1111111111111111111111111111111111111111")
	amount := sdk.NewCoins(sdk.NewCoin("ualtg", sdk.NewInt(1000000000000000000)))
	msgSend := banktypes.NewMsgSend(fromAddr, toAddr, amount)
	msgs := []sdk.Msg{msgSend}

	return txV
}

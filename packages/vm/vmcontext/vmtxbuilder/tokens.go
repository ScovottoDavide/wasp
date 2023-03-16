package vmtxbuilder

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/vmcontext/vmexceptions"
)

// nativeTokenBalance represents on-chain account of the specific native token
type nativeTokenBalance struct {
	nativeTokenID iotago.NativeTokenID
	outputID      iotago.OutputID     // if in != nil, otherwise zeroOutputID
	in            *iotago.BasicOutput // if nil it means output does not exist, this is new account for the token_id
	out           *iotago.BasicOutput // current balance of the token_id on the chain
}

func (n *nativeTokenBalance) Clone() *nativeTokenBalance {
	nativeTokenID := iotago.NativeTokenID{}
	copy(nativeTokenID[:], n.nativeTokenID[:])

	outputID := iotago.OutputID{}
	copy(outputID[:], n.outputID[:])

	return &nativeTokenBalance{
		nativeTokenID: nativeTokenID,
		outputID:      outputID,
		in:            cloneInternalBasicOutputOrNil(n.in),
		out:           cloneInternalBasicOutputOrNil(n.out),
	}
}

// producesOutput if value update produces UTXO of the corresponding total native token balance
func (n *nativeTokenBalance) producesOutput() bool {
	if n.identicalInOut() {
		// value didn't change
		return false
	}
	if util.IsZeroBigInt(n.getOutValue()) {
		// end value is 0
		return false
	}
	return true
}

// requiresInput returns if value change requires input in the transaction
func (n *nativeTokenBalance) requiresInput() bool {
	if n.identicalInOut() {
		// value didn't change
		return false
	}
	return n.in != nil
}

func (n *nativeTokenBalance) getOutValue() *big.Int {
	return n.out.NativeTokens[0].Amount
}

func (n *nativeTokenBalance) add(delta *big.Int) *nativeTokenBalance {
	amount := new(big.Int).Add(n.getOutValue(), delta)
	if amount.Sign() < 0 {
		panic(fmt.Errorf("(id: %s, delta: %d): %v",
			n.nativeTokenID, delta, vm.ErrNotEnoughNativeAssetBalance))
	}
	if amount.Cmp(util.MaxUint256) > 0 {
		panic(vm.ErrOverflow)
	}
	n.out.NativeTokens[0].Amount = amount
	return n
}

// updateMinSD uptates the resulting output to have the minimum SD
func (n *nativeTokenBalance) updateMinSD() {
	minSD := parameters.L1().Protocol.RentStructure.MinRent(n.out)
	if minSD > n.out.Amount {
		// sd for internal output can only ever increase
		n.out.Amount = minSD
	}
}

func (n *nativeTokenBalance) identicalInOut() bool {
	switch {
	case n.in == n.out:
		panic("identicalBasicOutputs: internal inconsistency 1")
	case n.in == nil || n.out == nil:
		return false
	case !n.in.Ident().Equal(n.out.Ident()):
		return false
	case n.in.Amount != n.out.Amount:
		return false
	case !n.in.NativeTokens.Equal(n.out.NativeTokens):
		return false
	case !n.in.Features.Equal(n.out.Features):
		return false
	case len(n.in.NativeTokens) != 1:
		panic("identicalBasicOutputs: internal inconsistency 2")
	case len(n.out.NativeTokens) != 1:
		panic("identicalBasicOutputs: internal inconsistency 3")
	case n.in.NativeTokens[0].ID != n.nativeTokenID:
		panic("identicalBasicOutputs: internal inconsistency 4")
	case n.out.NativeTokens[0].ID != n.nativeTokenID:
		panic("identicalBasicOutputs: internal inconsistency 5")
	}
	return true
}

func cloneInternalBasicOutputOrNil(o *iotago.BasicOutput) *iotago.BasicOutput {
	if o == nil {
		return nil
	}
	return o.Clone().(*iotago.BasicOutput)
}

func (txb *AnchorTransactionBuilder) newInternalTokenOutput(aliasID iotago.AliasID, nativeTokenID iotago.NativeTokenID) *iotago.BasicOutput {
	out := &iotago.BasicOutput{
		Amount: 0,
		NativeTokens: iotago.NativeTokens{{
			ID:     nativeTokenID,
			Amount: big.NewInt(0),
		}},
		Conditions: iotago.UnlockConditions{
			&iotago.AddressUnlockCondition{Address: aliasID.ToAddress()},
		},
		Features: iotago.Features{
			&iotago.SenderFeature{
				Address: aliasID.ToAddress(),
			},
		},
	}
	return out
}

func (txb *AnchorTransactionBuilder) nativeTokenOutputsSorted() []*nativeTokenBalance {
	ret := make([]*nativeTokenBalance, 0, len(txb.balanceNativeTokens))
	for _, f := range txb.balanceNativeTokens {
		if !f.requiresInput() && !f.producesOutput() {
			continue
		}
		ret = append(ret, f)
	}
	sort.Slice(ret, func(i, j int) bool {
		return bytes.Compare(ret[i].nativeTokenID[:], ret[j].nativeTokenID[:]) < 0
	})
	return ret
}

func (txb *AnchorTransactionBuilder) NativeTokenRecordsToBeUpdated() ([]iotago.NativeTokenID, []iotago.NativeTokenID) {
	toBeUpdated := make([]iotago.NativeTokenID, 0, len(txb.balanceNativeTokens))
	toBeRemoved := make([]iotago.NativeTokenID, 0, len(txb.balanceNativeTokens))
	for _, nt := range txb.nativeTokenOutputsSorted() {
		if nt.producesOutput() {
			toBeUpdated = append(toBeUpdated, nt.nativeTokenID)
		} else if nt.requiresInput() {
			toBeRemoved = append(toBeRemoved, nt.nativeTokenID)
		}
	}
	return toBeUpdated, toBeRemoved
}

func (txb *AnchorTransactionBuilder) NativeTokenOutputsByTokenIDs(ids []iotago.NativeTokenID) map[iotago.NativeTokenID]*iotago.BasicOutput {
	ret := make(map[iotago.NativeTokenID]*iotago.BasicOutput)
	for _, id := range ids {
		ret[id] = txb.balanceNativeTokens[id].out
	}
	return ret
}

// addNativeTokenBalanceDelta adds delta to the token balance. Use negative delta to subtract.
// The call may result in adding new token ID to the ledger or disappearing one
// This impacts storage deposit amount locked in the internal UTXOs which keep respective balances
// Returns delta of required storage deposit
func (txb *AnchorTransactionBuilder) addNativeTokenBalanceDelta(nativeTokenID iotago.NativeTokenID, delta *big.Int) int64 {
	if util.IsZeroBigInt(delta) {
		return 0
	}
	nt := txb.ensureNativeTokenBalance(nativeTokenID).add(delta)

	if nt.identicalInOut() {
		return 0
	}

	if util.IsZeroBigInt(nt.getOutValue()) {
		// 0 native tokens on the output side
		if nt.in == nil {
			// in this case the internar accounting output that would be created is not needed anymore, reiburse the SD
			return int64(nt.out.Amount)
		}
		return int64(nt.in.Amount)
	}

	// update the SD in case the storage deposit has changed from the last time this output was used
	oldSD := nt.out.Amount
	nt.updateMinSD()
	updatedSD := nt.out.Amount

	return int64(oldSD) - int64(updatedSD)
}

// ensureNativeTokenBalance makes sure that cached output is in the builder
// if not, it asks for the in balance by calling the loader function
// Panics if the call results to exceeded limits
func (txb *AnchorTransactionBuilder) ensureNativeTokenBalance(nativeTokenID iotago.NativeTokenID) *nativeTokenBalance {
	if nativeTokenBalance, exists := txb.balanceNativeTokens[nativeTokenID]; exists {
		return nativeTokenBalance
	}

	basicOutputIn, outputID := txb.accountsView.NativeTokenOutput(nativeTokenID) // output will be nil if no such token id accounted yet
	if basicOutputIn != nil {
		if txb.InputsAreFull() {
			panic(vmexceptions.ErrInputLimitExceeded)
		}
		if txb.outputsAreFull() {
			panic(vmexceptions.ErrOutputLimitExceeded)
		}
	}

	var basicOutputOut *iotago.BasicOutput
	if basicOutputIn == nil {
		basicOutputOut = txb.newInternalTokenOutput(util.AliasIDFromAliasOutput(txb.anchorOutput, txb.anchorOutputID), nativeTokenID)
	} else {
		basicOutputOut = cloneInternalBasicOutputOrNil(basicOutputIn)
	}

	nativeTokenBalance := &nativeTokenBalance{
		nativeTokenID: nativeTokenID,
		outputID:      outputID,
		in:            basicOutputIn,
		out:           basicOutputOut,
	}
	txb.balanceNativeTokens[nativeTokenID] = nativeTokenBalance
	return nativeTokenBalance
}

package vmcontext

import (
	"errors"
	"fmt"
	"time"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/blob"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/migrations"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/execution"
	"github.com/iotaledger/wasp/packages/vm/gas"
	"github.com/iotaledger/wasp/packages/vm/processors"
	"github.com/iotaledger/wasp/packages/vm/vmcontext/vmtxbuilder"
)

// VMContext represents state of the chain during one run of the VM while processing
// a batch of requests. VMContext object mutates with each request in the batch.
// The VMContext is created from immutable vm.VMTask object and UTXO state of the
// chain address contained in the statetxbuilder.Builder
type VMContext struct {
	task *vm.VMTask
	// same for the block
	chainOwnerID        isc.AgentID
	finalStateTimestamp time.Time
	blockContext        map[isc.Hname]interface{}
	txbuilder           *vmtxbuilder.AnchorTransactionBuilder
	txsnapshot          *vmtxbuilder.AnchorTransactionBuilder
	gasBurnedTotal      uint64
	gasFeeChargedTotal  uint64

	// ---- request context
	chainInfo          *isc.ChainInfo
	req                isc.Request
	NumPostedOutputs   int // how many outputs has been posted in the request
	requestIndex       uint16
	requestEventIndex  uint16
	currentStateUpdate *StateUpdate
	entropy            hashing.HashValue
	callStack          []*callContext
	// --- gas related
	// max tokens that can be charged for gas fee
	gasMaxTokensToSpendForGasFee uint64
	// final gas budget set for the run
	gasBudgetAdjusted uint64
	// is gas bur enabled
	gasBurnEnabled bool
	// gas already burned
	gasBurned uint64
	// tokens charged
	gasFeeCharged uint64
	// burn history. If disabled, it is nil
	gasBurnLog *gas.BurnLog

	// used to set caller = nil when executing "open/close block context" funcs (meaning caller is the VM itself)
	callerIsVM bool
}

var _ execution.WaspContext = &VMContext{}

type callContext struct {
	caller             isc.AgentID // calling agent
	contract           isc.Hname   // called contract
	params             isc.Params  // params passed
	allowanceAvailable *isc.Assets // MUTABLE: allowance budget left after TransferAllowedFunds
}

// CreateVMContext creates a context for the whole batch run
func CreateVMContext(task *vm.VMTask) *VMContext {
	// assert consistency. It is a bit redundant double check
	if len(task.Requests) == 0 {
		// should never happen
		panic(errors.New("CreateVMContext.invalid params: must be at least 1 request"))
	}
	prevL1Commitment, err := transaction.L1CommitmentFromAliasOutput(task.AnchorOutput)
	if err != nil {
		// should never happen
		panic(fmt.Errorf("CreateVMContext: can't parse state data as L1Commitment from chain input %w", err))
	}

	task.StateDraft, err = task.Store.NewStateDraft(task.TimeAssumption, prevL1Commitment)
	if err != nil {
		// should never happen
		panic(err)
	}

	ret := &VMContext{
		task:                task,
		finalStateTimestamp: task.TimeAssumption.Add(time.Duration(len(task.Requests)+1) * time.Nanosecond),
		blockContext:        make(map[isc.Hname]interface{}),
		entropy:             task.Entropy,
		callStack:           make([]*callContext, 0),
	}
	if task.EnableGasBurnLogging {
		ret.gasBurnLog = gas.NewGasBurnLog()
	}
	// at the beginning of each block
	l1Commitment, err := transaction.L1CommitmentFromAliasOutput(task.AnchorOutput)
	if err != nil {
		// should never happen
		panic(err)
	}

	var totalL2Funds *isc.Assets
	ret.withStateUpdate(func() {
		ret.runMigrations(migrations.BaseSchemaVersion, migrations.Migrations)

		// save the anchor tx ID of the current state
		ret.callCore(blocklog.Contract, func(s kv.KVStore) {
			blocklog.UpdateLatestBlockInfo(
				s,
				ret.task.AnchorOutputID.TransactionID(),
				isc.NewAliasOutputWithID(ret.task.AnchorOutput, ret.task.AnchorOutputID),
				l1Commitment,
			)
		})
		// get the total L2 funds in accounting
		totalL2Funds = ret.loadTotalFungibleTokens()
	})

	task.AnchorOutputStorageDeposit = task.AnchorOutput.Amount - totalL2Funds.BaseTokens

	ret.txbuilder = vmtxbuilder.NewAnchorTransactionBuilder(
		task.AnchorOutput,
		task.AnchorOutputID,
		task.AnchorOutputStorageDeposit,
		vmtxbuilder.AccountsContractRead{
			NativeTokenOutput:   ret.loadNativeTokenOutput,
			FoundryOutput:       ret.loadFoundry,
			NFTOutput:           ret.loadNFT,
			TotalFungibleTokens: ret.loadTotalFungibleTokens,
		},
	)

	return ret
}

func (vmctx *VMContext) withStateUpdate(f func()) {
	vmctx.currentStateUpdate = NewStateUpdate()
	f()
	vmctx.currentStateUpdate.Mutations.ApplyTo(vmctx.task.StateDraft)
	vmctx.currentStateUpdate = nil
}

// CloseVMContext does the closing actions on the block
// return nil for normal block and rotation address for rotation block
func (vmctx *VMContext) CloseVMContext(numRequests, numSuccess, numOffLedger uint16) (uint32, *state.L1Commitment, time.Time, iotago.Address) {
	vmctx.GasBurnEnable(false)
	var rotationAddr iotago.Address
	vmctx.withStateUpdate(func() {
		rotationAddr = vmctx.saveBlockInfo(numRequests, numSuccess, numOffLedger)
		vmctx.closeBlockContexts()
		vmctx.saveInternalUTXOs()
	})

	block := vmctx.task.Store.ExtractBlock(vmctx.task.StateDraft)

	l1Commitment := block.L1Commitment()

	blockIndex := vmctx.task.StateDraft.BlockIndex()
	timestamp := vmctx.task.StateDraft.Timestamp()

	return blockIndex, l1Commitment, timestamp, rotationAddr
}

func (vmctx *VMContext) checkRotationAddress() (ret iotago.Address) {
	vmctx.callCore(governance.Contract, func(s kv.KVStore) {
		ret = governance.GetRotationAddress(s)
	})
	return
}

// saveBlockInfo is in the blocklog partition context. Returns rotation address if this block is a rotation block
func (vmctx *VMContext) saveBlockInfo(numRequests, numSuccess, numOffLedger uint16) iotago.Address {
	if rotationAddress := vmctx.checkRotationAddress(); rotationAddress != nil {
		// block was marked fake by the governance contract because it is a committee rotation.
		// There was only on request in the block
		// We skip saving block information in order to avoid inconsistencies
		return rotationAddress
	}

	blockInfo := &blocklog.BlockInfo{
		SchemaVersion:         blocklog.BlockInfoLatestSchemaVersion,
		Timestamp:             vmctx.task.StateDraft.Timestamp(),
		TotalRequests:         numRequests,
		NumSuccessfulRequests: numSuccess,
		NumOffLedgerRequests:  numOffLedger,
		PreviousAliasOutput:   isc.NewAliasOutputWithID(vmctx.task.AnchorOutput, vmctx.task.AnchorOutputID),
		GasBurned:             vmctx.gasBurnedTotal,
		GasFeeCharged:         vmctx.gasFeeChargedTotal,
	}

	// TODO this "SaveControlAddressesIfNecessary" call is saving potentially outdated info.
	// Regardless the "control addresses in the state" can be completely removed, these are useless as the info can be derived from the AO
	vmctx.callCore(blocklog.Contract, func(s kv.KVStore) {
		blocklog.SaveNextBlockInfo(s, blockInfo)
		blocklog.SaveControlAddressesIfNecessary(
			s,
			vmctx.task.AnchorOutput.StateController(),
			vmctx.task.AnchorOutput.GovernorAddress(),
			vmctx.task.AnchorOutput.StateIndex,
		)
	})
	vmctx.task.Log.Debugf("saved blockinfo: %s", blockInfo)
	return nil
}

// OpenBlockContexts calls the block context open function for all subscribed core contracts
func (vmctx *VMContext) OpenBlockContexts() {
	if vmctx.gasBurnEnabled {
		panic("expected gasBurnEnabled == false")
	}

	vmctx.withStateUpdate(func() {
		vmctx.loadChainConfig()

		var subs []root.BlockContextSubscription
		vmctx.callCore(root.Contract, func(s kv.KVStore) {
			subs = root.GetBlockContextSubscriptions(s)
		})
		vmctx.callerIsVM = true
		for _, sub := range subs {
			vmctx.callProgram(sub.Contract, sub.OpenFunc, nil, nil)
		}
		vmctx.callerIsVM = false
	})
}

// closeBlockContexts closes block contexts in deterministic FIFO sequence
func (vmctx *VMContext) closeBlockContexts() {
	if vmctx.gasBurnEnabled {
		panic("expected gasBurnEnabled == false")
	}
	var subs []root.BlockContextSubscription
	vmctx.callCore(root.Contract, func(s kv.KVStore) {
		subs = root.GetBlockContextSubscriptions(s)
	})
	vmctx.callerIsVM = true
	for i := len(subs) - 1; i >= 0; i-- {
		vmctx.callProgram(subs[i].Contract, subs[i].CloseFunc, nil, nil)
	}
	vmctx.callerIsVM = false
}

// saveInternalUTXOs relies on the order of the outputs in the anchor tx. If that order changes, this will be broken.
// Anchor Transaction outputs order must be:
// 0. Anchor Output
// 1. NativeTokens
// 2. Foundries
// 3. NFTs
func (vmctx *VMContext) saveInternalUTXOs() {
	// create a mock AO, with a nil statecommitment, just to calculate changes in the minimum SD
	mockAO := vmctx.txbuilder.CreateAnchorOutput(vmctx.StateMetadata(state.L1CommitmentNil))
	newMinSD := parameters.L1().Protocol.RentStructure.MinRent(mockAO)
	oldMinSD := vmctx.task.AnchorOutputStorageDeposit
	changeInSD := int64(oldMinSD) - int64(newMinSD)

	if changeInSD != 0 {
		vmctx.task.Log.Debugf("adjusting commonAccount because AO SD cost changed, old:%d new:%d", oldMinSD, newMinSD)
		// update the commonAccount with the change in SD cost
		vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
			accounts.AdjustAccountBaseTokens(s, accounts.CommonAccount(), changeInSD)
		})
	}

	nativeTokenIDs, nativeTokensToBeRemoved := vmctx.txbuilder.NativeTokenRecordsToBeUpdated()
	nativeTokensOutputsToBeUpdated := vmctx.txbuilder.NativeTokenOutputsByTokenIDs(nativeTokenIDs)

	foundryIDs, foundriesToBeRemoved := vmctx.txbuilder.FoundriesToBeUpdated()
	foundrySNToBeUpdated := vmctx.txbuilder.FoundryOutputsBySN(foundryIDs)

	NFTOutputsToBeAdded, NFTOutputsToBeRemoved := vmctx.txbuilder.NFTOutputsToBeUpdated()

	blockIndex := vmctx.task.AnchorOutput.StateIndex + 1
	outputIndex := uint16(1)

	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		// update native token outputs
		for _, out := range nativeTokensOutputsToBeUpdated {
			vmctx.task.Log.Debugf("saving NT %s, outputIndex: %d", out.NativeTokens[0].ID, outputIndex)
			accounts.SaveNativeTokenOutput(s, out, blockIndex, outputIndex)
			outputIndex++
		}
		for _, id := range nativeTokensToBeRemoved {
			vmctx.task.Log.Debugf("deleting NT %s", id)
			accounts.DeleteNativeTokenOutput(s, id)
		}

		// update foundry UTXOs
		for _, out := range foundrySNToBeUpdated {
			vmctx.task.Log.Debugf("saving foundry %d, outputIndex: %d", out.SerialNumber, outputIndex)
			accounts.SaveFoundryOutput(s, out, blockIndex, outputIndex)
			outputIndex++
		}
		for _, sn := range foundriesToBeRemoved {
			vmctx.task.Log.Debugf("deleting foundry %s", sn)
			accounts.DeleteFoundryOutput(s, sn)
		}

		// update NFT Outputs
		for _, out := range NFTOutputsToBeAdded {
			vmctx.task.Log.Debugf("saving NFT %s, outputIndex: %d", out.NFTID, outputIndex)
			accounts.SaveNFTOutput(s, out, blockIndex, outputIndex)
			outputIndex++
		}
		for _, out := range NFTOutputsToBeRemoved {
			vmctx.task.Log.Debugf("deleting NFT %s", out.NFTID)
			accounts.DeleteNFTOutput(s, out.NFTID)
		}
	})
}

func (vmctx *VMContext) AssertConsistentGasTotals() {
	var sumGasBurned, sumGasFeeCharged uint64

	for _, r := range vmctx.task.Results {
		sumGasBurned += r.Receipt.GasBurned
		sumGasFeeCharged += r.Receipt.GasFeeCharged
	}
	if vmctx.gasBurnedTotal != sumGasBurned {
		panic("vmctx.gasBurnedTotal != sumGasBurned")
	}
	if vmctx.gasFeeChargedTotal != sumGasFeeCharged {
		panic("vmctx.gasFeeChargedTotal != sumGasFeeCharged")
	}
}

func (vmctx *VMContext) LocateProgram(programHash hashing.HashValue) (vmtype string, binary []byte, err error) {
	vmctx.callCore(blob.Contract, func(s kv.KVStore) {
		vmtype, binary, err = blob.LocateProgram(vmctx.State(), programHash)
	})
	return vmtype, binary, err
}

func (vmctx *VMContext) Processors() *processors.Cache {
	return vmctx.task.Processors
}

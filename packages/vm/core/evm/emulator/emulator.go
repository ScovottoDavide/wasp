// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package emulator

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/params"
	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/iotaledger/wasp/packages/evm/evmtypes"
	"github.com/iotaledger/wasp/packages/evm/evmutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/subrealm"
	"github.com/iotaledger/wasp/packages/util/panicutil"
	"github.com/iotaledger/wasp/packages/vm/core/evm"
	"github.com/iotaledger/wasp/packages/vm/vmcontext/vmexceptions"
)

type EVMEmulator struct {
	timestamp   uint64
	gasLimits   GasLimits
	chainConfig *params.ChainConfig
	kv          kv.KVStore
	vmConfig    vm.Config
	l2Balance   L2Balance
}

type GasLimits struct {
	Block uint64
	Call  uint64
}

var configCache *lru.Cache[int, *params.ChainConfig]

func init() {
	var err error
	configCache, err = lru.New[int, *params.ChainConfig](100)
	if err != nil {
		panic(err)
	}
}

func getConfig(chainID int) *params.ChainConfig {
	if c, ok := configCache.Get(chainID); ok {
		return c
	}
	c := &params.ChainConfig{
		ChainID:             big.NewInt(int64(chainID)),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.Hash{},
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		Ethash:              &params.EthashConfig{},
	}
	configCache.Add(chainID, c)
	return c
}

const (
	keyStateDB      = "s"
	keyBlockchainDB = "b"
)

func newStateDB(store kv.KVStore, l2Balance L2Balance) *StateDB {
	return NewStateDB(subrealm.New(store, keyStateDB), l2Balance)
}

func newBlockchainDB(store kv.KVStore, blockGasLimit uint64) *BlockchainDB {
	return NewBlockchainDB(subrealm.New(store, keyBlockchainDB), blockGasLimit)
}

// Init initializes the EVM state with the provided genesis allocation parameters
func Init(
	store kv.KVStore,
	chainID uint16,
	blockKeepAmount int32,
	gasLimits GasLimits,
	timestamp uint64,
	alloc core.GenesisAlloc,
) {
	bdb := newBlockchainDB(store, gasLimits.Block)
	if bdb.Initialized() {
		panic("evm state already initialized in kvstore")
	}
	bdb.Init(chainID, blockKeepAmount, timestamp)

	statedb := newStateDB(store, nil)
	for addr, account := range alloc {
		statedb.CreateAccount(addr)
		if account.Balance != nil {
			panic("balances must be 0 at genesis")
		}
		if account.Code != nil {
			statedb.SetCode(addr, account.Code)
		}
		for k, v := range account.Storage {
			statedb.SetState(addr, k, v)
		}
		statedb.SetNonce(addr, account.Nonce)
	}
}

func NewEVMEmulator(
	store kv.KVStore,
	timestamp uint64,
	gasLimits GasLimits,
	magicContracts map[common.Address]vm.ISCMagicContract,
	l2Balance L2Balance,
) *EVMEmulator {
	bdb := newBlockchainDB(store, gasLimits.Block)
	if !bdb.Initialized() {
		panic("must initialize genesis block first")
	}

	return &EVMEmulator{
		timestamp:   timestamp,
		gasLimits:   gasLimits,
		chainConfig: getConfig(int(bdb.GetChainID())),
		kv:          store,
		vmConfig:    vm.Config{MagicContracts: magicContracts},
		l2Balance:   l2Balance,
	}
}

func (e *EVMEmulator) StateDB() *StateDB {
	return newStateDB(e.kv, e.l2Balance)
}

func (e *EVMEmulator) BlockchainDB() *BlockchainDB {
	return newBlockchainDB(e.kv, e.gasLimits.Block)
}

func (e *EVMEmulator) BlockGasLimit() uint64 {
	return e.gasLimits.Block
}

func (e *EVMEmulator) CallGasLimit() uint64 {
	return e.gasLimits.Call
}

func (e *EVMEmulator) ChainContext() core.ChainContext {
	return &chainContext{
		engine: ethash.NewFaker(),
	}
}

// CallContract executes a contract call, without committing changes to the state
func (e *EVMEmulator) CallContract(call ethereum.CallMsg, gasBurnEnable func(bool)) (*core.ExecutionResult, error) {
	// Ensure message is initialized properly.
	if call.Gas == 0 {
		call.Gas = e.gasLimits.Call
	}
	if call.Value == nil {
		call.Value = big.NewInt(0)
	}

	pendingHeader := e.BlockchainDB().GetPendingHeader()

	// run the EVM code on a buffered state (so that writes are not committed)
	statedb := e.StateDB().Buffered().StateDB()

	return e.applyMessage(callMsg{call}, statedb, pendingHeader, gasBurnEnable, nil)
}

func (e *EVMEmulator) applyMessage(
	msg callMsg,
	statedb vm.StateDB,
	header *types.Header,
	gasBurnEnable func(bool),
	tracer tracers.Tracer,
) (res *core.ExecutionResult, err error) {
	blockContext := core.NewEVMBlockContext(header, e.ChainContext(), nil)
	txContext := core.NewEVMTxContext(msg)

	vmConfig := e.vmConfig
	vmConfig.Tracer = tracer
	vmConfig.Debug = vmConfig.Tracer != nil

	vmEnv := vm.NewEVM(blockContext, txContext, statedb, e.chainConfig, vmConfig)

	if msg.CallMsg.Gas > e.gasLimits.Call {
		msg.CallMsg.Gas = e.gasLimits.Call
	}

	gasPool := core.GasPool(msg.Gas())
	vmEnv.Reset(txContext, statedb)
	if gasBurnEnable != nil {
		gasBurnEnable(true)
		defer gasBurnEnable(false)
	}

	caughtErr := panicutil.CatchAllExcept(
		func() {
			// catch any exceptions during the execution, so that an EVM receipt is produced
			res, err = core.ApplyMessage(vmEnv, msg, &gasPool)
		},
		vmexceptions.AllProtocolLimits...,
	)
	if caughtErr != nil {
		return nil, caughtErr
	}
	return res, err
}

func (e *EVMEmulator) SendTransaction(
	tx *types.Transaction,
	gasBurnEnable func(bool),
	tracer tracers.Tracer,
) (*types.Receipt, *core.ExecutionResult, error) {
	buf := e.StateDB().Buffered()
	statedb := buf.StateDB()
	pendingHeader := e.BlockchainDB().GetPendingHeader()

	sender, err := types.Sender(e.Signer(), tx)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid transaction: %w", err)
	}
	nonce := e.StateDB().GetNonce(sender)
	if tx.Nonce() != nonce {
		return nil, nil, fmt.Errorf("invalid transaction nonce: got %d, want %d", tx.Nonce(), nonce)
	}

	msg, err := tx.AsMessage(types.MakeSigner(e.chainConfig, pendingHeader.Number), pendingHeader.BaseFee)
	if err != nil {
		return nil, nil, err
	}

	msgWithZeroGasPrice := callMsg{
		CallMsg: ethereum.CallMsg{
			From:       msg.From(),
			To:         msg.To(),
			Gas:        msg.Gas(),
			GasPrice:   big.NewInt(0),
			GasFeeCap:  big.NewInt(0),
			GasTipCap:  big.NewInt(0),
			Value:      msg.Value(),
			Data:       msg.Data(),
			AccessList: msg.AccessList(),
		},
	}

	result, err := e.applyMessage(
		msgWithZeroGasPrice,
		statedb,
		pendingHeader,
		gasBurnEnable,
		tracer,
	)

	gasUsed := uint64(0)
	if result != nil {
		gasUsed = result.UsedGas
	}

	cumulativeGasUsed := gasUsed
	index := uint(0)
	latest := e.BlockchainDB().GetLatestPendingReceipt()
	if latest != nil {
		cumulativeGasUsed += latest.CumulativeGasUsed
		index = latest.TransactionIndex + 1
	}

	receipt := &types.Receipt{
		Type:              tx.Type(),
		CumulativeGasUsed: cumulativeGasUsed,
		TxHash:            tx.Hash(),
		GasUsed:           gasUsed,
		Logs:              statedb.GetLogs(tx.Hash()),
		BlockNumber:       pendingHeader.Number,
		TransactionIndex:  index,
	}
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	if result == nil || result.Failed() {
		receipt.Status = types.ReceiptStatusFailed
	} else {
		receipt.Status = types.ReceiptStatusSuccessful
	}

	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(msg.From(), tx.Nonce())
	}

	buf.Commit()
	e.BlockchainDB().AddTransaction(tx, receipt)

	return receipt, result, err
}

func (e *EVMEmulator) MintBlock() {
	e.BlockchainDB().MintBlock(e.timestamp)
}

// FilterLogs executes a log filter operation, blocking during execution and
// returning all the results in one batch.
func (e *EVMEmulator) FilterLogs(query *ethereum.FilterQuery) ([]*types.Log, error) {
	receipts, err := e.getReceiptsInFilterRange(query)
	if err != nil {
		return nil, err
	}
	return e.filterLogs(query, receipts)
}

const (
	maxBlocksInFilterRange = 1_000
	maxLogsInResult        = 10_000
)

func (e *EVMEmulator) getReceiptsInFilterRange(query *ethereum.FilterQuery) ([]*types.Receipt, error) {
	bc := e.BlockchainDB()

	if query.BlockHash != nil {
		blockNumber, ok := bc.GetBlockNumberByBlockHash(*query.BlockHash)
		if !ok {
			return nil, nil
		}
		return bc.GetReceiptsByBlockNumber(blockNumber), nil
	}

	// Initialize unset filter boundaries to run from genesis to chain head
	first := big.NewInt(1) // skip genesis since it has no logs
	last := new(big.Int).SetUint64(bc.GetNumber())
	from := first
	if query.FromBlock != nil && query.FromBlock.Cmp(first) >= 0 && query.FromBlock.Cmp(last) <= 0 {
		from = query.FromBlock
	}
	to := last
	if query.ToBlock != nil && query.ToBlock.Cmp(first) >= 0 && query.ToBlock.Cmp(last) <= 0 {
		to = query.ToBlock
	}

	if !from.IsUint64() || !to.IsUint64() {
		return nil, errors.New("block number is too large")
	}
	var receipts []*types.Receipt
	{
		from := from.Uint64()
		to := to.Uint64()
		if to > from && to-from > maxBlocksInFilterRange {
			return nil, errors.New("too many blocks in filter range")
		}
		for i := from; i <= to; i++ {
			receipts = append(receipts, bc.GetReceiptsByBlockNumber(i)...)
		}
	}
	return receipts, nil
}

func (e *EVMEmulator) filterLogs(query *ethereum.FilterQuery, receipts []*types.Receipt) ([]*types.Log, error) {
	var logs []*types.Log
	for _, r := range receipts {
		if !evmtypes.BloomFilter(r.Bloom, query.Addresses, query.Topics) {
			continue
		}
		for _, log := range r.Logs {
			if !evmtypes.LogMatches(log, query.Addresses, query.Topics) {
				continue
			}
			if len(logs) >= maxLogsInResult {
				return nil, errors.New("too many logs in result")
			}
			logs = append(logs, log)
		}
	}
	return logs, nil
}

func (e *EVMEmulator) Signer() types.Signer {
	return evmutil.Signer(e.chainConfig.ChainID)
}

// callMsg implements core.Message to allow passing it as a transaction simulator.
type callMsg struct {
	ethereum.CallMsg
}

func (m callMsg) From() common.Address         { return m.CallMsg.From }
func (m callMsg) Nonce() uint64                { return 0 }
func (m callMsg) IsFake() bool                 { return true }
func (m callMsg) To() *common.Address          { return m.CallMsg.To }
func (m callMsg) GasPrice() *big.Int           { return evm.GasPrice } // we ignore the gas price set by the sender
func (m callMsg) GasFeeCap() *big.Int          { return m.CallMsg.GasFeeCap }
func (m callMsg) GasTipCap() *big.Int          { return m.CallMsg.GasTipCap }
func (m callMsg) Gas() uint64                  { return m.CallMsg.Gas }
func (m callMsg) Value() *big.Int              { return m.CallMsg.Value }
func (m callMsg) Data() []byte                 { return m.CallMsg.Data }
func (m callMsg) AccessList() types.AccessList { return m.CallMsg.AccessList }

type chainContext struct {
	engine consensus.Engine
}

var _ core.ChainContext = &chainContext{}

func (c *chainContext) Engine() consensus.Engine {
	return c.engine
}

func (c *chainContext) GetHeader(common.Hash, uint64) *types.Header {
	panic("not implemented")
}

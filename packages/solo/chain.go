// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package solo

import (
	"bytes"
	"fmt"
	"math"
	"os"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/coreutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/kv/kvdecoder"
	"github.com/iotaledger/wasp/packages/kv/trie"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/blob"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/core/errors"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/gas"
	"github.com/iotaledger/wasp/packages/vm/vmtypes"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

func (ch *Chain) L1Params() *parameters.L1 {
	return ch.Env.L1Params()
}

// String is string representation for main parameters of the chain
//goland:noinspection ALL
func (ch *Chain) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Chain ID: %s\n", ch.ChainID)
	fmt.Fprintf(&buf, "Chain state controller: %s\n", ch.StateControllerAddress)
	fmt.Fprintf(&buf, "Root commitment: %s\n", trie.RootCommitment(ch.State.TrieNodeStore()))
	fmt.Fprintf(&buf, "UTXODB genesis address: %s\n", ch.Env.utxoDB.GenesisAddress())
	return buf.String()
}

// DumpAccounts dumps all account balances into the human-readable string
func (ch *Chain) DumpAccounts() string {
	_, chainOwnerID, _ := ch.GetInfo()
	ret := fmt.Sprintf("ChainID: %s\nChain owner: %s\n",
		ch.ChainID.String(),
		chainOwnerID.String(ch.L1Params().Protocol.Bech32HRP),
	)
	acc := ch.L2Accounts()
	for i := range acc {
		aid := acc[i]
		ret += fmt.Sprintf("  %s:\n", aid.String(ch.L1Params().Protocol.Bech32HRP))
		bals := ch.L2Assets(aid)
		ret += fmt.Sprintf("%s\n", bals.String())
	}
	return ret
}

// FindContract is a view call to the 'root' smart contract on the chain.
// It returns blobCache record of the deployed smart contract with the given name
func (ch *Chain) FindContract(scName string) (*root.ContractRecord, error) {
	retDict, err := ch.CallView(root.Contract.Name, root.ViewFindContract.Name,
		root.ParamHname, iscp.Hn(scName),
	)
	if err != nil {
		return nil, err
	}
	retBin, err := retDict.Get(root.ParamContractRecData)
	if err != nil {
		return nil, err
	}
	if retBin == nil {
		return nil, fmt.Errorf("smart contract '%s' not found", scName)
	}
	record, err := root.ContractRecordFromBytes(retBin)
	if err != nil {
		return nil, err
	}
	if record.Name != scName {
		return nil, fmt.Errorf("smart contract '%s' not found", scName)
	}
	return record, err
}

// GetBlobInfo return info about blob with the given hash with existence flag
// The blob information is returned as a map of pairs 'blobFieldName': 'fieldDataLength'
func (ch *Chain) GetBlobInfo(blobHash hashing.HashValue) (map[string]uint32, bool) {
	res, err := ch.CallView(blob.Contract.Name, blob.ViewGetBlobInfo.Name, blob.ParamHash, blobHash)
	require.NoError(ch.Env.T, err)
	if res.IsEmpty() {
		return nil, false
	}
	ret, err := blob.DecodeSizesMap(res)
	require.NoError(ch.Env.T, err)
	return ret, true
}

func (ch *Chain) GetGasFeePolicy() *gas.GasFeePolicy {
	res, err := ch.CallView(governance.Contract.Name, governance.ViewGetFeePolicy.Name)
	require.NoError(ch.Env.T, err)
	fpBin := res.MustGet(governance.ParamFeePolicyBytes)
	feePolicy, err := gas.FeePolicyFromBytes(fpBin)
	require.NoError(ch.Env.T, err)
	return feePolicy
}

// UploadBlob calls core 'blob' smart contract blob.FuncStoreBlob entry point to upload blob
// data to the chain. It returns hash of the blob, the unique identifier of it.
// The parameters must be either a dict.Dict, or a sequence of pairs 'fieldName': 'fieldValue'
// Requires at least 2 x gasFeeEstimate to be on sender's L2 account
func (ch *Chain) UploadBlob(user *cryptolib.KeyPair, params ...interface{}) (ret hashing.HashValue, err error) {
	if user == nil {
		user = ch.OriginatorPrivateKey
	}

	blobAsADict := parseParams(params)
	expectedHash := blob.MustGetBlobHash(blobAsADict)
	if _, ok := ch.GetBlobInfo(expectedHash); ok {
		// blob exists, return hash of existing
		return expectedHash, nil
	}
	req := NewCallParams(blob.Contract.Name, blob.FuncStoreBlob.Name, params...)
	g, _, err := ch.EstimateGasOffLedger(req, nil, true)
	if err != nil {
		return [32]byte{}, err
	}
	req.WithGasBudget(g)
	res, err := ch.PostRequestOffLedger(req, user)
	if err != nil {
		return
	}
	resBin := res.MustGet(blob.ParamHash)
	if resBin == nil {
		err = xerrors.Errorf("internal error: no hash returned")
		return
	}
	ret, err = codec.DecodeHashValue(resBin)
	if err != nil {
		return
	}
	require.EqualValues(ch.Env.T, expectedHash, ret)
	return ret, err
}

// UploadBlobFromFile uploads blob from file data in the specified blob field plus optional other fields
func (ch *Chain) UploadBlobFromFile(keyPair *cryptolib.KeyPair, fileName string, fieldName string, params ...interface{}) (hashing.HashValue, error) {
	fileBinary, err := os.ReadFile(fileName)
	if err != nil {
		return hashing.HashValue{}, err
	}
	par := parseParams(params)
	par.Set(kv.Key(fieldName), fileBinary)
	return ch.UploadBlob(keyPair, par)
}

// UploadWasm is a syntactic sugar of the UploadBlob used to upload Wasm binary to the chain.
//  parameter 'binaryCode' is the binary of Wasm smart contract program
//
// The blob for the Wasm binary used fixed field names which are statically known by the
// 'root' smart contract which is responsible for the deployment of contracts on the chain
func (ch *Chain) UploadWasm(keyPair *cryptolib.KeyPair, binaryCode []byte) (ret hashing.HashValue, err error) {
	return ch.UploadBlob(keyPair,
		blob.VarFieldVMType, vmtypes.WasmTime,
		blob.VarFieldProgramBinary, binaryCode,
	)
}

// UploadWasmFromFile is a syntactic sugar to upload file content as blob data to the chain
func (ch *Chain) UploadWasmFromFile(keyPair *cryptolib.KeyPair, fileName string) (hashing.HashValue, error) {
	var binary []byte
	binary, err := os.ReadFile(fileName)
	if err != nil {
		return hashing.HashValue{}, err
	}
	return ch.UploadWasm(keyPair, binary)
}

// GetWasmBinary retrieves program binary in the format of Wasm blob from the chain by hash.
func (ch *Chain) GetWasmBinary(progHash hashing.HashValue) ([]byte, error) {
	res, err := ch.CallView(blob.Contract.Name, blob.ViewGetBlobField.Name,
		blob.ParamHash, progHash,
		blob.ParamField, blob.VarFieldVMType,
	)
	if err != nil {
		return nil, err
	}
	require.EqualValues(ch.Env.T, vmtypes.WasmTime, string(res.MustGet(blob.ParamBytes)))

	res, err = ch.CallView(blob.Contract.Name, blob.ViewGetBlobField.Name,
		blob.ParamHash, progHash,
		blob.ParamField, blob.VarFieldProgramBinary,
	)
	if err != nil {
		return nil, err
	}
	binary := res.MustGet(blob.ParamBytes)
	return binary, nil
}

// DeployContract deploys contract with the given name by its 'programHash'. 'sigScheme' represents
// the private key of the creator (nil defaults to chain originator). The 'creator' becomes an immutable
// property of the contract instance.
// The parameter 'programHash' can be one of the following:
//   - it is and ID of  the blob stored on the chain in the format of Wasm binary
//   - it can be a hash (ID) of the example smart contract ("hardcoded"). The "hardcoded"
//     smart contract must be made available with the call examples.AddProcessor
func (ch *Chain) DeployContract(user *cryptolib.KeyPair, name string, programHash hashing.HashValue, params ...interface{}) error {
	par := codec.MakeDict(map[string]interface{}{
		root.ParamProgramHash: programHash,
		root.ParamName:        name,
	})
	for k, v := range parseParams(params) {
		par[k] = v
	}
	_, err := ch.PostRequestSync(
		NewCallParams(root.Contract.Name, root.FuncDeployContract.Name, par).
			WithGasBudget(math.MaxUint64),
		user,
	)
	return err
}

// DeployWasmContract is syntactic sugar for uploading Wasm binary from file and
// deploying the smart contract in one call
func (ch *Chain) DeployWasmContract(keyPair *cryptolib.KeyPair, name, fname string, params ...interface{}) error {
	hprog, err := ch.UploadWasmFromFile(keyPair, fname)
	if err != nil {
		return err
	}
	return ch.DeployContract(keyPair, name, hprog, params...)
}

// GetInfo return main parameters of the chain:
//  - chainID
//  - agentID of the chain owner
//  - blobCache of contract deployed on the chain in the form of map 'contract hname': 'contract record'
func (ch *Chain) GetInfo() (*iscp.ChainID, *iscp.AgentID, map[iscp.Hname]*root.ContractRecord) {
	res, err := ch.CallView(governance.Contract.Name, governance.ViewGetChainInfo.Name)
	require.NoError(ch.Env.T, err)

	chainID, err := codec.DecodeChainID(res.MustGet(governance.VarChainID))
	require.NoError(ch.Env.T, err)

	chainOwnerID, err := codec.DecodeAgentID(res.MustGet(governance.VarChainOwnerID))
	require.NoError(ch.Env.T, err)

	res, err = ch.CallView(root.Contract.Name, root.ViewGetContractRecords.Name)
	require.NoError(ch.Env.T, err)

	contracts, err := root.DecodeContractRegistry(collections.NewMapReadOnly(res, root.StateVarContractRegistry))
	require.NoError(ch.Env.T, err)
	return chainID, chainOwnerID, contracts
}

type DustInfo struct {
	TotalIotasInL2Accounts uint64
	TotalDustDeposit       uint64
	NumNativeTokens        int
}

func (d *DustInfo) Total() uint64 {
	return d.TotalIotasInL2Accounts + d.TotalDustDeposit*uint64(d.NumNativeTokens)
}

func (ch *Chain) GetTotalIotaInfo() *DustInfo {
	bi := ch.GetLatestBlockInfo()
	return &DustInfo{
		TotalIotasInL2Accounts: bi.TotalIotasInL2Accounts,
		TotalDustDeposit:       bi.TotalDustDeposit,
		NumNativeTokens:        len(ch.GetOnChainTokenIDs()),
	}
}

func eventsFromViewResult(t TestContext, viewResult dict.Dict) []string {
	recs := collections.NewArray16ReadOnly(viewResult, blocklog.ParamEvent)
	ret := make([]string, recs.MustLen())
	for i := range ret {
		data, err := recs.GetAt(uint16(i))
		require.NoError(t, err)
		ret[i] = string(data)
	}
	return ret
}

// GetEventsForContract calls the view in the  'blocklog' core smart contract to retrieve events for a given smart contract.
func (ch *Chain) GetEventsForContract(name string) ([]string, error) {
	viewResult, err := ch.CallView(
		blocklog.Contract.Name, blocklog.ViewGetEventsForContract.Name,
		blocklog.ParamContractHname, iscp.Hn(name),
	)
	if err != nil {
		return nil, err
	}

	return eventsFromViewResult(ch.Env.T, viewResult), nil
}

// GetEventsForRequest calls the view in the  'blocklog' core smart contract to retrieve events for a given request.
func (ch *Chain) GetEventsForRequest(reqID iscp.RequestID) ([]string, error) {
	viewResult, err := ch.CallView(
		blocklog.Contract.Name, blocklog.ViewGetEventsForRequest.Name,
		blocklog.ParamRequestID, reqID,
	)
	if err != nil {
		return nil, err
	}
	return eventsFromViewResult(ch.Env.T, viewResult), nil
}

// GetEventsForBlock calls the view in the 'blocklog' core smart contract to retrieve events for a given block.
func (ch *Chain) GetEventsForBlock(blockIndex uint32) ([]string, error) {
	viewResult, err := ch.CallView(
		blocklog.Contract.Name, blocklog.ViewGetEventsForBlock.Name,
		blocklog.ParamBlockIndex, blockIndex,
	)
	if err != nil {
		return nil, err
	}
	return eventsFromViewResult(ch.Env.T, viewResult), nil
}

// CommonAccount return the agentID of the common account (controlled by the owner)
func (ch *Chain) CommonAccount() *iscp.AgentID {
	return ch.ChainID.CommonAccount()
}

// GetLatestBlockInfo return BlockInfo for the latest block in the chain
func (ch *Chain) GetLatestBlockInfo() *blocklog.BlockInfo {
	ret, err := ch.CallView(blocklog.Contract.Name, blocklog.ViewGetLatestBlockInfo.Name)
	require.NoError(ch.Env.T, err)
	resultDecoder := kvdecoder.New(ret, ch.Log())
	blockIndex := resultDecoder.MustGetUint32(blocklog.ParamBlockIndex)
	blockInfoBin := resultDecoder.MustGetBytes(blocklog.ParamBlockInfo)

	blockInfo, err := blocklog.BlockInfoFromBytes(blockIndex, blockInfoBin)
	require.NoError(ch.Env.T, err)
	return blockInfo
}

func (ch *Chain) GetErrorMessageFormat(code iscp.VMErrorCode) (string, error) {
	ret, err := ch.CallView(errors.Contract.Name, errors.ViewGetErrorMessageFormat.Name,
		errors.ParamErrorCode, code.Bytes(),
	)
	if err != nil {
		return "", err
	}
	resultDecoder := kvdecoder.New(ret, ch.Log())
	messageFormat, err := resultDecoder.GetString(errors.ParamErrorMessageFormat)

	require.NoError(ch.Env.T, err)
	return messageFormat, nil
}

// GetBlockInfo return BlockInfo for the particular block index in the chain
func (ch *Chain) GetBlockInfo(blockIndex uint32) (*blocklog.BlockInfo, error) {
	ret, err := ch.CallView(blocklog.Contract.Name, blocklog.ViewGetBlockInfo.Name,
		blocklog.ParamBlockIndex, blockIndex)
	if err != nil {
		return nil, err
	}
	resultDecoder := kvdecoder.New(ret, ch.Log())
	blockInfoBin := resultDecoder.MustGetBytes(blocklog.ParamBlockInfo)

	blockInfo, err := blocklog.BlockInfoFromBytes(blockIndex, blockInfoBin)
	require.NoError(ch.Env.T, err)
	return blockInfo, nil
}

// IsRequestProcessed checks if the request is booked on the chain as processed
func (ch *Chain) IsRequestProcessed(reqID iscp.RequestID) bool {
	ret, err := ch.CallView(blocklog.Contract.Name, blocklog.ViewIsRequestProcessed.Name,
		blocklog.ParamRequestID, reqID)
	require.NoError(ch.Env.T, err)
	resultDecoder := kvdecoder.New(ret, ch.Log())
	isProcessed, err := resultDecoder.GetBool(blocklog.ParamRequestProcessed)
	require.NoError(ch.Env.T, err)
	return isProcessed
}

// GetRequestReceipt gets the log records for a particular request, the block index and request index in the block
func (ch *Chain) GetRequestReceipt(reqID iscp.RequestID) (*blocklog.RequestReceipt, bool) {
	ret, err := ch.CallView(blocklog.Contract.Name, blocklog.ViewGetRequestReceipt.Name,
		blocklog.ParamRequestID, reqID)
	require.NoError(ch.Env.T, err)
	resultDecoder := kvdecoder.New(ret, ch.Log())
	binRec, err := resultDecoder.GetBytes(blocklog.ParamRequestRecord)
	if err != nil || binRec == nil {
		return nil, false
	}
	ret1, err := blocklog.RequestReceiptFromBytes(binRec)

	require.NoError(ch.Env.T, err)
	ret1.BlockIndex = resultDecoder.MustGetUint32(blocklog.ParamBlockIndex)
	ret1.RequestIndex = resultDecoder.MustGetUint16(blocklog.ParamRequestIndex)

	return ret1, true
}

// GetRequestReceiptsForBlock returns all request log records for a particular block
func (ch *Chain) GetRequestReceiptsForBlock(blockIndex ...uint32) []*blocklog.RequestReceipt {
	var blockIdx uint32
	if len(blockIndex) == 0 {
		blockIdx = ch.GetLatestBlockInfo().BlockIndex
	} else {
		blockIdx = blockIndex[0]
	}

	res, err := ch.CallView(blocklog.Contract.Name, blocklog.ViewGetRequestReceiptsForBlock.Name,
		blocklog.ParamBlockIndex, blockIdx)
	if err != nil {
		return nil
	}
	recs := collections.NewArray16ReadOnly(res, blocklog.ParamRequestRecord)
	ret := make([]*blocklog.RequestReceipt, recs.MustLen())
	for i := range ret {
		data, err := recs.GetAt(uint16(i))
		require.NoError(ch.Env.T, err)
		ret[i], err = blocklog.RequestReceiptFromBytes(data)
		require.NoError(ch.Env.T, err)
		ret[i].WithBlockData(blockIdx, uint16(i))
	}
	return ret
}

// GetRequestIDsForBlock returns return the list of requestIDs settled in a particular block
func (ch *Chain) GetRequestIDsForBlock(blockIndex uint32) []iscp.RequestID {
	res, err := ch.CallView(blocklog.Contract.Name, blocklog.ViewGetRequestIDsForBlock.Name,
		blocklog.ParamBlockIndex, blockIndex)
	if err != nil {
		ch.Log().Warnf("GetRequestIDsForBlock: %v", err)
		return nil
	}
	recs := collections.NewArray16ReadOnly(res, blocklog.ParamRequestID)
	ret := make([]iscp.RequestID, recs.MustLen())
	for i := range ret {
		reqIDBin, err := recs.GetAt(uint16(i))
		require.NoError(ch.Env.T, err)
		ret[i], err = iscp.RequestIDFromBytes(reqIDBin)
		require.NoError(ch.Env.T, err)
	}
	return ret
}

// GetRequestReceiptsForBlockRange returns all request log records for range of blocks, inclusively.
// Upper bound is 'latest block' is set to 0
func (ch *Chain) GetRequestReceiptsForBlockRange(fromBlockIndex, toBlockIndex uint32) []*blocklog.RequestReceipt {
	if toBlockIndex == 0 {
		toBlockIndex = ch.GetLatestBlockInfo().BlockIndex
	}
	if fromBlockIndex > toBlockIndex {
		return nil
	}
	ret := make([]*blocklog.RequestReceipt, 0)
	for i := fromBlockIndex; i <= toBlockIndex; i++ {
		recs := ch.GetRequestReceiptsForBlock(i)
		require.True(ch.Env.T, i == 0 || len(recs) != 0)
		ret = append(ret, recs...)
	}
	return ret
}

func (ch *Chain) GetRequestReceiptsForBlockRangeAsStrings(fromBlockIndex, toBlockIndex uint32) []string {
	recs := ch.GetRequestReceiptsForBlockRange(fromBlockIndex, toBlockIndex)
	ret := make([]string, len(recs))
	for i := range ret {
		ret[i] = recs[i].String()
	}
	return ret
}

func (ch *Chain) GetControlAddresses() *blocklog.ControlAddresses {
	res, err := ch.CallView(blocklog.Contract.Name, blocklog.ViewControlAddresses.Name)
	require.NoError(ch.Env.T, err)
	par := kvdecoder.New(res, ch.Log())
	ret := &blocklog.ControlAddresses{
		StateAddress:     par.MustGetAddress(blocklog.ParamStateControllerAddress),
		GoverningAddress: par.MustGetAddress(blocklog.ParamGoverningAddress),
		SinceBlockIndex:  par.MustGetUint32(blocklog.ParamBlockIndex),
	}
	return ret
}

// AddAllowedStateController adds the address to the allowed state controlled address list
func (ch *Chain) AddAllowedStateController(addr iotago.Address, keyPair *cryptolib.KeyPair) error {
	req := NewCallParams(coreutil.CoreContractGovernance, governance.FuncAddAllowedStateControllerAddress.Name,
		governance.ParamStateControllerAddress, addr,
	)
	_, err := ch.PostRequestSync(req, keyPair)
	return err
}

// AddAllowedStateController adds the address to the allowed state controlled address list
func (ch *Chain) RemoveAllowedStateController(addr iotago.Address, keyPair *cryptolib.KeyPair) error {
	req := NewCallParams(coreutil.CoreContractGovernance, governance.FuncRemoveAllowedStateControllerAddress.Name,
		governance.ParamStateControllerAddress, addr,
	)
	_, err := ch.PostRequestSync(req, keyPair)
	return err
}

// AddAllowedStateController adds the address to the allowed state controlled address list
func (ch *Chain) GetAllowedStateControllerAddresses() []iotago.Address {
	res, err := ch.CallView(coreutil.CoreContractGovernance, governance.ViewGetAllowedStateControllerAddresses.Name)
	require.NoError(ch.Env.T, err)
	if len(res) == 0 {
		return nil
	}
	ret := make([]iotago.Address, 0)
	arr := collections.NewArray16ReadOnly(res, string(governance.ParamAllowedStateControllerAddresses))
	for i := uint16(0); i < arr.MustLen(); i++ {
		a, err := codec.DecodeAddress(arr.MustGetAt(i))
		require.NoError(ch.Env.T, err)
		ret = append(ret, a)
	}
	return ret
}

// RotateStateController rotates the chain to the new controller address.
// We assume self-governed chain here.
// Mostly use for the testing of committee rotation logic, otherwise not much needed for smart contract testing
func (ch *Chain) RotateStateController(newStateAddr iotago.Address, newStateKeyPair, ownerKeyPair *cryptolib.KeyPair) error {
	req := NewCallParams(coreutil.CoreContractGovernance, coreutil.CoreEPRotateStateController,
		coreutil.ParamStateControllerAddress, newStateAddr,
	).AddIotas(1)
	result := ch.postRequestSyncTxSpecial(req, ownerKeyPair)
	if result.Error == nil {
		ch.StateControllerAddress = newStateAddr
		ch.StateControllerKeyPair = newStateKeyPair
	}
	return result.Error
}

func (ch *Chain) postRequestSyncTxSpecial(req *CallParams, keyPair *cryptolib.KeyPair) *vm.RequestResult {
	tx, _, err := ch.RequestFromParamsToLedger(req, keyPair)
	require.NoError(ch.Env.T, err)
	reqs, err := ch.Env.RequestsForChain(tx, ch.ChainID)
	require.NoError(ch.Env.T, err)
	results := ch.runRequestsSync(reqs, "postSpecial")
	return results[0]
}

type L1L2AddressAssets struct {
	Address  iotago.Address
	AssetsL1 *iscp.FungibleTokens
	AssetsL2 *iscp.FungibleTokens
}

func (a *L1L2AddressAssets) String() string {
	return fmt.Sprintf("Address: %s\nL1 ftokens:\n  %s\nL2 ftokens:\n  %s", a.Address, a.AssetsL1, a.AssetsL2)
}

func getAddr(addrOrKeypair interface{}) iotago.Address {
	switch a := addrOrKeypair.(type) {
	case iotago.Address:
		return a
	case *cryptolib.KeyPair:
		return a.GetPublicKey().AsEd25519Address()
	case cryptolib.KeyPair:
		return a.GetPublicKey().AsEd25519Address()
	}
	panic(xerrors.Errorf("getAddr: wrong type %T", addrOrKeypair))
}

func (ch *Chain) L1L2Funds(addrOrKeypair interface{}) *L1L2AddressAssets {
	addr := getAddr(addrOrKeypair)
	return &L1L2AddressAssets{
		Address:  addr,
		AssetsL1: ch.Env.L1Assets(addr),
		AssetsL2: ch.L2Assets(iscp.NewAgentID(addr, 0)),
	}
}

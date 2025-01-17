// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

package coregovernance

import "github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"

type AddAllowedStateControllerAddressCall struct {
	Func   *wasmlib.ScFunc
	Params MutableAddAllowedStateControllerAddressParams
}

type AddCandidateNodeCall struct {
	Func   *wasmlib.ScFunc
	Params MutableAddCandidateNodeParams
}

type ChangeAccessNodesCall struct {
	Func   *wasmlib.ScFunc
	Params MutableChangeAccessNodesParams
}

type ClaimChainOwnershipCall struct {
	Func *wasmlib.ScFunc
}

type DelegateChainOwnershipCall struct {
	Func   *wasmlib.ScFunc
	Params MutableDelegateChainOwnershipParams
}

type RemoveAllowedStateControllerAddressCall struct {
	Func   *wasmlib.ScFunc
	Params MutableRemoveAllowedStateControllerAddressParams
}

type RevokeAccessNodeCall struct {
	Func   *wasmlib.ScFunc
	Params MutableRevokeAccessNodeParams
}

type RotateStateControllerCall struct {
	Func   *wasmlib.ScFunc
	Params MutableRotateStateControllerParams
}

type SetCustomMetadataCall struct {
	Func   *wasmlib.ScFunc
	Params MutableSetCustomMetadataParams
}

type SetEVMGasRatioCall struct {
	Func   *wasmlib.ScFunc
	Params MutableSetEVMGasRatioParams
}

type SetFeePolicyCall struct {
	Func   *wasmlib.ScFunc
	Params MutableSetFeePolicyParams
}

type SetGasLimitsCall struct {
	Func   *wasmlib.ScFunc
	Params MutableSetGasLimitsParams
}

type SetMaintenanceOffCall struct {
	Func *wasmlib.ScFunc
}

type SetMaintenanceOnCall struct {
	Func *wasmlib.ScFunc
}

type GetAllowedStateControllerAddressesCall struct {
	Func    *wasmlib.ScView
	Results ImmutableGetAllowedStateControllerAddressesResults
}

type GetChainInfoCall struct {
	Func    *wasmlib.ScView
	Results ImmutableGetChainInfoResults
}

type GetChainNodesCall struct {
	Func    *wasmlib.ScView
	Results ImmutableGetChainNodesResults
}

type GetChainOwnerCall struct {
	Func    *wasmlib.ScView
	Results ImmutableGetChainOwnerResults
}

type GetCustomMetadataCall struct {
	Func    *wasmlib.ScView
	Results ImmutableGetCustomMetadataResults
}

type GetEVMGasRatioCall struct {
	Func    *wasmlib.ScView
	Results ImmutableGetEVMGasRatioResults
}

type GetFeePolicyCall struct {
	Func    *wasmlib.ScView
	Results ImmutableGetFeePolicyResults
}

type GetGasLimitsCall struct {
	Func    *wasmlib.ScView
	Results ImmutableGetGasLimitsResults
}

type GetMaintenanceStatusCall struct {
	Func    *wasmlib.ScView
	Results ImmutableGetMaintenanceStatusResults
}

type Funcs struct{}

var ScFuncs Funcs

// Adds the given address to the list of identities that constitute the state controller.
func (sc Funcs) AddAllowedStateControllerAddress(ctx wasmlib.ScFuncCallContext) *AddAllowedStateControllerAddressCall {
	f := &AddAllowedStateControllerAddressCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncAddAllowedStateControllerAddress)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Adds a node to the list of candidates.
func (sc Funcs) AddCandidateNode(ctx wasmlib.ScFuncCallContext) *AddCandidateNodeCall {
	f := &AddCandidateNodeCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncAddCandidateNode)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Iterates through the given map of actions and applies them.
func (sc Funcs) ChangeAccessNodes(ctx wasmlib.ScFuncCallContext) *ChangeAccessNodesCall {
	f := &ChangeAccessNodesCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncChangeAccessNodes)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Claims the ownership of the chain if the caller matches the identity
// that was set in delegateChainOwnership().
func (sc Funcs) ClaimChainOwnership(ctx wasmlib.ScFuncCallContext) *ClaimChainOwnershipCall {
	return &ClaimChainOwnershipCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncClaimChainOwnership)}
}

// Sets the Agent ID o as the new owner for the chain.
// This change will only be effective once claimChainOwnership() is called by o.
func (sc Funcs) DelegateChainOwnership(ctx wasmlib.ScFuncCallContext) *DelegateChainOwnershipCall {
	f := &DelegateChainOwnershipCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncDelegateChainOwnership)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Removes the given address from the list of identities that constitute the state controller.
func (sc Funcs) RemoveAllowedStateControllerAddress(ctx wasmlib.ScFuncCallContext) *RemoveAllowedStateControllerAddressCall {
	f := &RemoveAllowedStateControllerAddressCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncRemoveAllowedStateControllerAddress)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Removes a node from the list of candidates.
func (sc Funcs) RevokeAccessNode(ctx wasmlib.ScFuncCallContext) *RevokeAccessNodeCall {
	f := &RevokeAccessNodeCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncRevokeAccessNode)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Called when the committee is about to be rotated to the given address.
// If it succeeds, the next state transition will become a governance transition,
// thus updating the state controller in the chain's Alias Output.
// If it fails, nothing happens.
func (sc Funcs) RotateStateController(ctx wasmlib.ScFuncCallContext) *RotateStateControllerCall {
	f := &RotateStateControllerCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncRotateStateController)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Changes optional extra metadata that is appended to the L1 AliasOutput.
func (sc Funcs) SetCustomMetadata(ctx wasmlib.ScFuncCallContext) *SetCustomMetadataCall {
	f := &SetCustomMetadataCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncSetCustomMetadata)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Sets the EVM gas ratio for the chain.
func (sc Funcs) SetEVMGasRatio(ctx wasmlib.ScFuncCallContext) *SetEVMGasRatioCall {
	f := &SetEVMGasRatioCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncSetEVMGasRatio)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Sets the fee policy for the chain.
func (sc Funcs) SetFeePolicy(ctx wasmlib.ScFuncCallContext) *SetFeePolicyCall {
	f := &SetFeePolicyCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncSetFeePolicy)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Sets the gas limits for the chain.
func (sc Funcs) SetGasLimits(ctx wasmlib.ScFuncCallContext) *SetGasLimitsCall {
	f := &SetGasLimitsCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncSetGasLimits)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Stops the maintenance mode.
func (sc Funcs) SetMaintenanceOff(ctx wasmlib.ScFuncCallContext) *SetMaintenanceOffCall {
	return &SetMaintenanceOffCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncSetMaintenanceOff)}
}

// Starts the chain maintenance mode, meaning no further requests
// will be processed except calls to the governance contract.
func (sc Funcs) SetMaintenanceOn(ctx wasmlib.ScFuncCallContext) *SetMaintenanceOnCall {
	return &SetMaintenanceOnCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncSetMaintenanceOn)}
}

// Returns the list of allowed state controllers.
func (sc Funcs) GetAllowedStateControllerAddresses(ctx wasmlib.ScViewCallContext) *GetAllowedStateControllerAddressesCall {
	f := &GetAllowedStateControllerAddressesCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetAllowedStateControllerAddresses)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns information about the chain.
func (sc Funcs) GetChainInfo(ctx wasmlib.ScViewCallContext) *GetChainInfoCall {
	f := &GetChainInfoCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetChainInfo)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the current access nodes and candidates.
func (sc Funcs) GetChainNodes(ctx wasmlib.ScViewCallContext) *GetChainNodesCall {
	f := &GetChainNodesCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetChainNodes)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the AgentID of the chain owner.
func (sc Funcs) GetChainOwner(ctx wasmlib.ScViewCallContext) *GetChainOwnerCall {
	f := &GetChainOwnerCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetChainOwner)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the extra metadata that is added to the chain AliasOutput.
func (sc Funcs) GetCustomMetadata(ctx wasmlib.ScViewCallContext) *GetCustomMetadataCall {
	f := &GetCustomMetadataCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetCustomMetadata)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the EVM gas ratio.
func (sc Funcs) GetEVMGasRatio(ctx wasmlib.ScViewCallContext) *GetEVMGasRatioCall {
	f := &GetEVMGasRatioCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetEVMGasRatio)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the fee policy.
func (sc Funcs) GetFeePolicy(ctx wasmlib.ScViewCallContext) *GetFeePolicyCall {
	f := &GetFeePolicyCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetFeePolicy)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the gas limits.
func (sc Funcs) GetGasLimits(ctx wasmlib.ScViewCallContext) *GetGasLimitsCall {
	f := &GetGasLimitsCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetGasLimits)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns whether the chain is undergoing maintenance.
func (sc Funcs) GetMaintenanceStatus(ctx wasmlib.ScViewCallContext) *GetMaintenanceStatusCall {
	f := &GetMaintenanceStatusCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetMaintenanceStatus)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

var exportMap = wasmlib.ScExportMap{
	Names: []string{
		FuncAddAllowedStateControllerAddress,
		FuncAddCandidateNode,
		FuncChangeAccessNodes,
		FuncClaimChainOwnership,
		FuncDelegateChainOwnership,
		FuncRemoveAllowedStateControllerAddress,
		FuncRevokeAccessNode,
		FuncRotateStateController,
		FuncSetCustomMetadata,
		FuncSetEVMGasRatio,
		FuncSetFeePolicy,
		FuncSetGasLimits,
		FuncSetMaintenanceOff,
		FuncSetMaintenanceOn,
		ViewGetAllowedStateControllerAddresses,
		ViewGetChainInfo,
		ViewGetChainNodes,
		ViewGetChainOwner,
		ViewGetCustomMetadata,
		ViewGetEVMGasRatio,
		ViewGetFeePolicy,
		ViewGetGasLimits,
		ViewGetMaintenanceStatus,
	},
	Funcs: []wasmlib.ScFuncContextFunction{
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
	},
	Views: []wasmlib.ScViewContextFunction{
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
	},
}

func OnDispatch(index int32) {
	if index == -1 {
		exportMap.Export()
		return
	}

	panic("Calling core contract?")
}

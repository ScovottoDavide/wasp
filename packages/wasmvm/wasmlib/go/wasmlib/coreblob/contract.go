// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

package coreblob

import "github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"

type StoreBlobCall struct {
	Func    *wasmlib.ScFunc
	Params  MutableStoreBlobParams
	Results ImmutableStoreBlobResults
}

type GetBlobFieldCall struct {
	Func    *wasmlib.ScView
	Params  MutableGetBlobFieldParams
	Results ImmutableGetBlobFieldResults
}

type GetBlobInfoCall struct {
	Func    *wasmlib.ScView
	Params  MutableGetBlobInfoParams
	Results ImmutableGetBlobInfoResults
}

type ListBlobsCall struct {
	Func    *wasmlib.ScView
	Results ImmutableListBlobsResults
}

type Funcs struct{}

var ScFuncs Funcs

// Stores a new blob in the registry.
func (sc Funcs) StoreBlob(ctx wasmlib.ScFuncCallContext) *StoreBlobCall {
	f := &StoreBlobCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncStoreBlob)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	wasmlib.NewCallResultsProxy(&f.Func.ScView, &f.Results.Proxy)
	return f
}

// Returns the chunk associated with the given blob field name.
func (sc Funcs) GetBlobField(ctx wasmlib.ScViewCallContext) *GetBlobFieldCall {
	f := &GetBlobFieldCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetBlobField)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the size of each chunk of the blob.
func (sc Funcs) GetBlobInfo(ctx wasmlib.ScViewCallContext) *GetBlobInfoCall {
	f := &GetBlobInfoCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetBlobInfo)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns a list of all blobs hashes in the registry and their sized.
func (sc Funcs) ListBlobs(ctx wasmlib.ScViewCallContext) *ListBlobsCall {
	f := &ListBlobsCall{Func: wasmlib.NewScView(ctx, HScName, HViewListBlobs)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

var exportMap = wasmlib.ScExportMap{
	Names: []string{
		FuncStoreBlob,
		ViewGetBlobField,
		ViewGetBlobInfo,
		ViewListBlobs,
	},
	Funcs: []wasmlib.ScFuncContextFunction{
		wasmlib.FuncError,
	},
	Views: []wasmlib.ScViewContextFunction{
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

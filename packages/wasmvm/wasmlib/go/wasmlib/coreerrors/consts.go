// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

package coreerrors

import "github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"

const (
	ScName        = "errors"
	ScDescription = "Errors Contract"
	HScName       = wasmtypes.ScHname(0x8f3a8bb3)
)

const (
	ParamErrorCode = "c"
	ParamTemplate  = "m"
)

const (
	ResultErrorCode = "c"
	ResultTemplate  = "m"
)

const (
	FuncRegisterError         = "registerError"
	ViewGetErrorMessageFormat = "getErrorMessageFormat"
)

const (
	HFuncRegisterError         = wasmtypes.ScHname(0x9be65f8e)
	HViewGetErrorMessageFormat = wasmtypes.ScHname(0x63fe7d56)
)

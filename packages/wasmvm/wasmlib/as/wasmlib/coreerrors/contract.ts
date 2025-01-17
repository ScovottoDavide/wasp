// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

import * as wasmlib from '../index';
import * as sc from './index';

export class RegisterErrorCall {
    func:    wasmlib.ScFunc;
    params:  sc.MutableRegisterErrorParams = new sc.MutableRegisterErrorParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableRegisterErrorResults = new sc.ImmutableRegisterErrorResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncCallContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncRegisterError);
    }
}

export class GetErrorMessageFormatCall {
    func:    wasmlib.ScView;
    params:  sc.MutableGetErrorMessageFormatParams = new sc.MutableGetErrorMessageFormatParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableGetErrorMessageFormatResults = new sc.ImmutableGetErrorMessageFormatResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetErrorMessageFormat);
    }
}

export class ScFuncs {
    // Registers an error message template.
    // note that this function must be call()ed
    static registerError(ctx: wasmlib.ScFuncCallContext): RegisterErrorCall {
        const f = new RegisterErrorCall(ctx);
        f.params = new sc.MutableRegisterErrorParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableRegisterErrorResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the message template stored for a given error code.
    static getErrorMessageFormat(ctx: wasmlib.ScViewCallContext): GetErrorMessageFormatCall {
        const f = new GetErrorMessageFormatCall(ctx);
        f.params = new sc.MutableGetErrorMessageFormatParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableGetErrorMessageFormatResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }
}

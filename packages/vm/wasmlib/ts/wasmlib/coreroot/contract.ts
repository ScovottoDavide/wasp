// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmlib from "wasmlib";
import * as sc from "./index";

export class DeployContractCall {
	func: wasmlib.ScFunc = new wasmlib.ScFunc(sc.HScName, sc.HFuncDeployContract);
	params: sc.MutableDeployContractParams = new sc.MutableDeployContractParams();
}

export class GrantDeployPermissionCall {
	func: wasmlib.ScFunc = new wasmlib.ScFunc(sc.HScName, sc.HFuncGrantDeployPermission);
	params: sc.MutableGrantDeployPermissionParams = new sc.MutableGrantDeployPermissionParams();
}

export class RevokeDeployPermissionCall {
	func: wasmlib.ScFunc = new wasmlib.ScFunc(sc.HScName, sc.HFuncRevokeDeployPermission);
	params: sc.MutableRevokeDeployPermissionParams = new sc.MutableRevokeDeployPermissionParams();
}

export class FindContractCall {
	func: wasmlib.ScView = new wasmlib.ScView(sc.HScName, sc.HViewFindContract);
	params: sc.MutableFindContractParams = new sc.MutableFindContractParams();
	results: sc.ImmutableFindContractResults = new sc.ImmutableFindContractResults();
}

export class GetContractRecordsCall {
	func: wasmlib.ScView = new wasmlib.ScView(sc.HScName, sc.HViewGetContractRecords);
	results: sc.ImmutableGetContractRecordsResults = new sc.ImmutableGetContractRecordsResults();
}

export class ScFuncs {

    static deployContract(ctx: wasmlib.ScFuncCallContext): DeployContractCall {
        let f = new DeployContractCall();
        f.func.setPtrs(f.params, null);
        return f;
    }

    static grantDeployPermission(ctx: wasmlib.ScFuncCallContext): GrantDeployPermissionCall {
        let f = new GrantDeployPermissionCall();
        f.func.setPtrs(f.params, null);
        return f;
    }

    static revokeDeployPermission(ctx: wasmlib.ScFuncCallContext): RevokeDeployPermissionCall {
        let f = new RevokeDeployPermissionCall();
        f.func.setPtrs(f.params, null);
        return f;
    }

    static findContract(ctx: wasmlib.ScViewCallContext): FindContractCall {
        let f = new FindContractCall();
        f.func.setPtrs(f.params, f.results);
        return f;
    }

    static getContractRecords(ctx: wasmlib.ScViewCallContext): GetContractRecordsCall {
        let f = new GetContractRecordsCall();
        f.func.setPtrs(null, f.results);
        return f;
    }
}
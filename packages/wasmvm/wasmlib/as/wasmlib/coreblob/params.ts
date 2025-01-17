// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

import * as wasmtypes from '../wasmtypes';
import * as sc from './index';

export class MapStringToImmutableBytes extends wasmtypes.ScProxy {

    getBytes(key: string): wasmtypes.ScImmutableBytes {
        return new wasmtypes.ScImmutableBytes(this.proxy.key(wasmtypes.stringToBytes(key)));
    }
}

export class ImmutableStoreBlobParams extends wasmtypes.ScProxy {
    // named chunks
    blobs(): sc.MapStringToImmutableBytes {
        return new sc.MapStringToImmutableBytes(this.proxy);
    }

    // data schema for external tools
    dataSchema(): wasmtypes.ScImmutableBytes {
        return new wasmtypes.ScImmutableBytes(this.proxy.root(sc.ParamDataSchema));
    }

    // smart contract program binary code
    progBinary(): wasmtypes.ScImmutableBytes {
        return new wasmtypes.ScImmutableBytes(this.proxy.root(sc.ParamProgBinary));
    }

    // smart contract program source code
    sources(): wasmtypes.ScImmutableBytes {
        return new wasmtypes.ScImmutableBytes(this.proxy.root(sc.ParamSources));
    }

    // VM type that must be used to run progBinary
    vmType(): wasmtypes.ScImmutableString {
        return new wasmtypes.ScImmutableString(this.proxy.root(sc.ParamVMType));
    }
}

export class MapStringToMutableBytes extends wasmtypes.ScProxy {

    clear(): void {
        this.proxy.clearMap();
    }

    getBytes(key: string): wasmtypes.ScMutableBytes {
        return new wasmtypes.ScMutableBytes(this.proxy.key(wasmtypes.stringToBytes(key)));
    }
}

export class MutableStoreBlobParams extends wasmtypes.ScProxy {
    // named chunks
    blobs(): sc.MapStringToMutableBytes {
        return new sc.MapStringToMutableBytes(this.proxy);
    }

    // data schema for external tools
    dataSchema(): wasmtypes.ScMutableBytes {
        return new wasmtypes.ScMutableBytes(this.proxy.root(sc.ParamDataSchema));
    }

    // smart contract program binary code
    progBinary(): wasmtypes.ScMutableBytes {
        return new wasmtypes.ScMutableBytes(this.proxy.root(sc.ParamProgBinary));
    }

    // smart contract program source code
    sources(): wasmtypes.ScMutableBytes {
        return new wasmtypes.ScMutableBytes(this.proxy.root(sc.ParamSources));
    }

    // VM type that must be used to run progBinary
    vmType(): wasmtypes.ScMutableString {
        return new wasmtypes.ScMutableString(this.proxy.root(sc.ParamVMType));
    }
}

export class ImmutableGetBlobFieldParams extends wasmtypes.ScProxy {
    // chunk name
    field(): wasmtypes.ScImmutableString {
        return new wasmtypes.ScImmutableString(this.proxy.root(sc.ParamField));
    }

    // hash of the blob
    hash(): wasmtypes.ScImmutableHash {
        return new wasmtypes.ScImmutableHash(this.proxy.root(sc.ParamHash));
    }
}

export class MutableGetBlobFieldParams extends wasmtypes.ScProxy {
    // chunk name
    field(): wasmtypes.ScMutableString {
        return new wasmtypes.ScMutableString(this.proxy.root(sc.ParamField));
    }

    // hash of the blob
    hash(): wasmtypes.ScMutableHash {
        return new wasmtypes.ScMutableHash(this.proxy.root(sc.ParamHash));
    }
}

export class ImmutableGetBlobInfoParams extends wasmtypes.ScProxy {
    // hash of the blob
    hash(): wasmtypes.ScImmutableHash {
        return new wasmtypes.ScImmutableHash(this.proxy.root(sc.ParamHash));
    }
}

export class MutableGetBlobInfoParams extends wasmtypes.ScProxy {
    // hash of the blob
    hash(): wasmtypes.ScMutableHash {
        return new wasmtypes.ScMutableHash(this.proxy.root(sc.ParamHash));
    }
}

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmtypes from "wasmlib/wasmtypes";

export const ScName        = "blocklog";
export const ScDescription = "Block log contract";
export const HScName       = new wasmtypes.ScHname(0xf538ef2b);

export const ParamBlockIndex    = "n";
export const ParamContractHname = "h";
export const ParamFromBlock     = "f";
export const ParamRequestID     = "u";
export const ParamToBlock       = "t";

export const ResultBlockIndex             = "n";
export const ResultBlockInfo              = "i";
export const ResultEvent                  = "e";
export const ResultGoverningAddress       = "g";
export const ResultRequestID              = "u";
export const ResultRequestIndex           = "r";
export const ResultRequestProcessed       = "p";
export const ResultRequestRecord          = "d";
export const ResultStateControllerAddress = "s";

export const ViewControlAddresses           = "controlAddresses";
export const ViewGetBlockInfo               = "getBlockInfo";
export const ViewGetEventsForBlock          = "getEventsForBlock";
export const ViewGetEventsForContract       = "getEventsForContract";
export const ViewGetEventsForRequest        = "getEventsForRequest";
export const ViewGetBlockInfo         = "getLatestBlockInfo";
export const ViewGetRequestIDsForBlock      = "getRequestIDsForBlock";
export const ViewGetRequestReceipt          = "getRequestReceipt";
export const ViewGetRequestReceiptsForBlock = "getRequestReceiptsForBlock";
export const ViewIsRequestProcessed         = "isRequestProcessed";

export const HViewControlAddresses           = new wasmtypes.ScHname(0x796bd223);
export const HViewGetBlockInfo               = new wasmtypes.ScHname(0xbe89f9b3);
export const HViewGetEventsForBlock          = new wasmtypes.ScHname(0x36232798);
export const HViewGetEventsForContract       = new wasmtypes.ScHname(0x682a1922);
export const HViewGetEventsForRequest        = new wasmtypes.ScHname(0x4f8d68e4);
export const HViewGetBlockInfo         = new wasmtypes.ScHname(0x084a1760);
export const HViewGetRequestIDsForBlock      = new wasmtypes.ScHname(0x5a20327a);
export const HViewGetRequestReceipt          = new wasmtypes.ScHname(0xb7f9534f);
export const HViewGetRequestReceiptsForBlock = new wasmtypes.ScHname(0x77e3beef);
export const HViewIsRequestProcessed         = new wasmtypes.ScHname(0xd57d50a9);

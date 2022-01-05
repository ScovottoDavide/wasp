// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

package dividend

import "github.com/iotaledger/wasp/packages/vm/wasmlib/go/wasmlib"

type ArrayOfImmutableAddress struct {
	objID int32
}

func (a ArrayOfImmutableAddress) Length() int32 {
	return wasmlib.GetLength(a.objID)
}

func (a ArrayOfImmutableAddress) GetAddress(index int32) wasmlib.ScImmutableAddress {
	return wasmlib.NewScImmutableAddress(a.objID, wasmlib.Key32(index))
}

type MapAddressToImmutableInt64 struct {
	objID int32
}

func (m MapAddressToImmutableInt64) GetInt64(key wasmlib.ScAddress) wasmlib.ScImmutableInt64 {
	return wasmlib.NewScImmutableInt64(m.objID, key.KeyID())
}

type ImmutableDividendState struct {
	id int32
}

func (s ImmutableDividendState) MemberList() ArrayOfImmutableAddress {
	arrID := wasmlib.GetObjectID(s.id, idxMap[IdxStateMemberList], wasmlib.TYPE_ARRAY|wasmlib.TYPE_ADDRESS)
	return ArrayOfImmutableAddress{objID: arrID}
}

func (s ImmutableDividendState) Members() MapAddressToImmutableInt64 {
	mapID := wasmlib.GetObjectID(s.id, idxMap[IdxStateMembers], wasmlib.TYPE_MAP)
	return MapAddressToImmutableInt64{objID: mapID}
}

func (s ImmutableDividendState) Owner() wasmlib.ScImmutableAgentID {
	return wasmlib.NewScImmutableAgentID(s.id, idxMap[IdxStateOwner])
}

func (s ImmutableDividendState) TotalFactor() wasmlib.ScImmutableInt64 {
	return wasmlib.NewScImmutableInt64(s.id, idxMap[IdxStateTotalFactor])
}

type ArrayOfMutableAddress struct {
	objID int32
}

func (a ArrayOfMutableAddress) Clear() {
	wasmlib.Clear(a.objID)
}

func (a ArrayOfMutableAddress) Length() int32 {
	return wasmlib.GetLength(a.objID)
}

func (a ArrayOfMutableAddress) GetAddress(index int32) wasmlib.ScMutableAddress {
	return wasmlib.NewScMutableAddress(a.objID, wasmlib.Key32(index))
}

type MapAddressToMutableInt64 struct {
	objID int32
}

func (m MapAddressToMutableInt64) Clear() {
	wasmlib.Clear(m.objID)
}

func (m MapAddressToMutableInt64) GetInt64(key wasmlib.ScAddress) wasmlib.ScMutableInt64 {
	return wasmlib.NewScMutableInt64(m.objID, key.KeyID())
}

type MutableDividendState struct {
	id int32
}

func (s MutableDividendState) AsImmutable() ImmutableDividendState {
	return ImmutableDividendState(s)
}

func (s MutableDividendState) MemberList() ArrayOfMutableAddress {
	arrID := wasmlib.GetObjectID(s.id, idxMap[IdxStateMemberList], wasmlib.TYPE_ARRAY|wasmlib.TYPE_ADDRESS)
	return ArrayOfMutableAddress{objID: arrID}
}

func (s MutableDividendState) Members() MapAddressToMutableInt64 {
	mapID := wasmlib.GetObjectID(s.id, idxMap[IdxStateMembers], wasmlib.TYPE_MAP)
	return MapAddressToMutableInt64{objID: mapID}
}

func (s MutableDividendState) Owner() wasmlib.ScMutableAgentID {
	return wasmlib.NewScMutableAgentID(s.id, idxMap[IdxStateOwner])
}

func (s MutableDividendState) TotalFactor() wasmlib.ScMutableInt64 {
	return wasmlib.NewScMutableInt64(s.id, idxMap[IdxStateTotalFactor])
}

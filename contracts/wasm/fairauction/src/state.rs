// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

#![allow(dead_code)]
#![allow(unused_imports)]

use wasmlib::*;
use wasmlib::host::*;

use crate::*;
use crate::keys::*;
use crate::structs::*;
use crate::typedefs::*;

#[derive(Clone, Copy)]
pub struct MapColorToImmutableAuction {
	pub(crate) obj_id: i32,
}

impl MapColorToImmutableAuction {
    pub fn get_auction(&self, key: &ScColor) -> ImmutableAuction {
        ImmutableAuction { obj_id: self.obj_id, key_id: key.get_key_id() }
    }
}

#[derive(Clone, Copy)]
pub struct MapColorToImmutableBidderList {
	pub(crate) obj_id: i32,
}

impl MapColorToImmutableBidderList {
    pub fn get_bidder_list(&self, key: &ScColor) -> ImmutableBidderList {
        let sub_id = get_object_id(self.obj_id, key.get_key_id(), TYPE_ARRAY | TYPE_AGENT_ID);
        ImmutableBidderList { obj_id: sub_id }
    }
}

#[derive(Clone, Copy)]
pub struct MapColorToImmutableBids {
	pub(crate) obj_id: i32,
}

impl MapColorToImmutableBids {
    pub fn get_bids(&self, key: &ScColor) -> ImmutableBids {
        let sub_id = get_object_id(self.obj_id, key.get_key_id(), TYPE_MAP);
        ImmutableBids { obj_id: sub_id }
    }
}

#[derive(Clone, Copy)]
pub struct ImmutableFairAuctionState {
    pub(crate) id: i32,
}

impl ImmutableFairAuctionState {
    pub fn auctions(&self) -> MapColorToImmutableAuction {
		let map_id = get_object_id(self.id, idx_map(IDX_STATE_AUCTIONS), TYPE_MAP);
		MapColorToImmutableAuction { obj_id: map_id }
	}

    pub fn bidder_list(&self) -> MapColorToImmutableBidderList {
		let map_id = get_object_id(self.id, idx_map(IDX_STATE_BIDDER_LIST), TYPE_MAP);
		MapColorToImmutableBidderList { obj_id: map_id }
	}

    pub fn bids(&self) -> MapColorToImmutableBids {
		let map_id = get_object_id(self.id, idx_map(IDX_STATE_BIDS), TYPE_MAP);
		MapColorToImmutableBids { obj_id: map_id }
	}

    pub fn owner_margin(&self) -> ScImmutableInt64 {
		ScImmutableInt64::new(self.id, idx_map(IDX_STATE_OWNER_MARGIN))
	}
}

#[derive(Clone, Copy)]
pub struct MapColorToMutableAuction {
	pub(crate) obj_id: i32,
}

impl MapColorToMutableAuction {
    pub fn clear(&self) {
        clear(self.obj_id);
    }

    pub fn get_auction(&self, key: &ScColor) -> MutableAuction {
        MutableAuction { obj_id: self.obj_id, key_id: key.get_key_id() }
    }
}

#[derive(Clone, Copy)]
pub struct MapColorToMutableBidderList {
	pub(crate) obj_id: i32,
}

impl MapColorToMutableBidderList {
    pub fn clear(&self) {
        clear(self.obj_id);
    }

    pub fn get_bidder_list(&self, key: &ScColor) -> MutableBidderList {
        let sub_id = get_object_id(self.obj_id, key.get_key_id(), TYPE_ARRAY | TYPE_AGENT_ID);
        MutableBidderList { obj_id: sub_id }
    }
}

#[derive(Clone, Copy)]
pub struct MapColorToMutableBids {
	pub(crate) obj_id: i32,
}

impl MapColorToMutableBids {
    pub fn clear(&self) {
        clear(self.obj_id);
    }

    pub fn get_bids(&self, key: &ScColor) -> MutableBids {
        let sub_id = get_object_id(self.obj_id, key.get_key_id(), TYPE_MAP);
        MutableBids { obj_id: sub_id }
    }
}

#[derive(Clone, Copy)]
pub struct MutableFairAuctionState {
    pub(crate) id: i32,
}

impl MutableFairAuctionState {
    pub fn as_immutable(&self) -> ImmutableFairAuctionState {
		ImmutableFairAuctionState { id: self.id }
	}

    pub fn auctions(&self) -> MapColorToMutableAuction {
		let map_id = get_object_id(self.id, idx_map(IDX_STATE_AUCTIONS), TYPE_MAP);
		MapColorToMutableAuction { obj_id: map_id }
	}

    pub fn bidder_list(&self) -> MapColorToMutableBidderList {
		let map_id = get_object_id(self.id, idx_map(IDX_STATE_BIDDER_LIST), TYPE_MAP);
		MapColorToMutableBidderList { obj_id: map_id }
	}

    pub fn bids(&self) -> MapColorToMutableBids {
		let map_id = get_object_id(self.id, idx_map(IDX_STATE_BIDS), TYPE_MAP);
		MapColorToMutableBids { obj_id: map_id }
	}

    pub fn owner_margin(&self) -> ScMutableInt64 {
		ScMutableInt64::new(self.id, idx_map(IDX_STATE_OWNER_MARGIN))
	}
}

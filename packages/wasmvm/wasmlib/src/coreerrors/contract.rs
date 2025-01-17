// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

#![allow(dead_code)]

use crate::*;
use crate::coreerrors::*;

pub struct RegisterErrorCall<'a> {
    pub func:    ScFunc<'a>,
    pub params:  MutableRegisterErrorParams,
    pub results: ImmutableRegisterErrorResults,
}

pub struct GetErrorMessageFormatCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableGetErrorMessageFormatParams,
    pub results: ImmutableGetErrorMessageFormatResults,
}

pub struct ScFuncs {
}

impl ScFuncs {
    // Registers an error message template.
    // note that this function must be call()ed
    pub fn register_error(ctx: &impl ScFuncCallContext) -> RegisterErrorCall {
        let mut f = RegisterErrorCall {
            func:    ScFunc::new(ctx, HSC_NAME, HFUNC_REGISTER_ERROR),
            params:  MutableRegisterErrorParams { proxy: Proxy::nil() },
            results: ImmutableRegisterErrorResults { proxy: Proxy::nil() },
        };
        ScFunc::link_params(&mut f.params.proxy, &f.func);
        ScFunc::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns the message template stored for a given error code.
    pub fn get_error_message_format(ctx: &impl ScViewCallContext) -> GetErrorMessageFormatCall {
        let mut f = GetErrorMessageFormatCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_GET_ERROR_MESSAGE_FORMAT),
            params:  MutableGetErrorMessageFormatParams { proxy: Proxy::nil() },
            results: ImmutableGetErrorMessageFormatResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }
}

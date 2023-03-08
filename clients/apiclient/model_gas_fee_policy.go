/*
Wasp API

REST API for the Wasp node

API version: 0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
)

// checks if the GasFeePolicy type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GasFeePolicy{}

// GasFeePolicy struct for GasFeePolicy
type GasFeePolicy struct {
	EvmGasRatio Ratio32 `json:"evmGasRatio"`
	GasPerToken Ratio32 `json:"gasPerToken"`
	// The validator fee share.
	ValidatorFeeShare int32 `json:"validatorFeeShare"`
}

// NewGasFeePolicy instantiates a new GasFeePolicy object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGasFeePolicy(evmGasRatio Ratio32, gasPerToken Ratio32, validatorFeeShare int32) *GasFeePolicy {
	this := GasFeePolicy{}
	this.EvmGasRatio = evmGasRatio
	this.GasPerToken = gasPerToken
	this.ValidatorFeeShare = validatorFeeShare
	return &this
}

// NewGasFeePolicyWithDefaults instantiates a new GasFeePolicy object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGasFeePolicyWithDefaults() *GasFeePolicy {
	this := GasFeePolicy{}
	return &this
}

// GetEvmGasRatio returns the EvmGasRatio field value
func (o *GasFeePolicy) GetEvmGasRatio() Ratio32 {
	if o == nil {
		var ret Ratio32
		return ret
	}

	return o.EvmGasRatio
}

// GetEvmGasRatioOk returns a tuple with the EvmGasRatio field value
// and a boolean to check if the value has been set.
func (o *GasFeePolicy) GetEvmGasRatioOk() (*Ratio32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EvmGasRatio, true
}

// SetEvmGasRatio sets field value
func (o *GasFeePolicy) SetEvmGasRatio(v Ratio32) {
	o.EvmGasRatio = v
}

// GetGasPerToken returns the GasPerToken field value
func (o *GasFeePolicy) GetGasPerToken() Ratio32 {
	if o == nil {
		var ret Ratio32
		return ret
	}

	return o.GasPerToken
}

// GetGasPerTokenOk returns a tuple with the GasPerToken field value
// and a boolean to check if the value has been set.
func (o *GasFeePolicy) GetGasPerTokenOk() (*Ratio32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GasPerToken, true
}

// SetGasPerToken sets field value
func (o *GasFeePolicy) SetGasPerToken(v Ratio32) {
	o.GasPerToken = v
}

// GetValidatorFeeShare returns the ValidatorFeeShare field value
func (o *GasFeePolicy) GetValidatorFeeShare() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.ValidatorFeeShare
}

// GetValidatorFeeShareOk returns a tuple with the ValidatorFeeShare field value
// and a boolean to check if the value has been set.
func (o *GasFeePolicy) GetValidatorFeeShareOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ValidatorFeeShare, true
}

// SetValidatorFeeShare sets field value
func (o *GasFeePolicy) SetValidatorFeeShare(v int32) {
	o.ValidatorFeeShare = v
}

func (o GasFeePolicy) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GasFeePolicy) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["evmGasRatio"] = o.EvmGasRatio
	toSerialize["gasPerToken"] = o.GasPerToken
	toSerialize["validatorFeeShare"] = o.ValidatorFeeShare
	return toSerialize, nil
}

type NullableGasFeePolicy struct {
	value *GasFeePolicy
	isSet bool
}

func (v NullableGasFeePolicy) Get() *GasFeePolicy {
	return v.value
}

func (v *NullableGasFeePolicy) Set(val *GasFeePolicy) {
	v.value = val
	v.isSet = true
}

func (v NullableGasFeePolicy) IsSet() bool {
	return v.isSet
}

func (v *NullableGasFeePolicy) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGasFeePolicy(val *GasFeePolicy) *NullableGasFeePolicy {
	return &NullableGasFeePolicy{value: val, isSet: true}
}

func (v NullableGasFeePolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGasFeePolicy) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


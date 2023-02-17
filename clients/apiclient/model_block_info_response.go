/*
Wasp API

REST API for the Wasp node

API version: 0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
	"time"
)

// checks if the BlockInfoResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BlockInfoResponse{}

// BlockInfoResponse struct for BlockInfoResponse
type BlockInfoResponse struct {
	AnchorTransactionId string `json:"anchorTransactionId"`
	BlockIndex uint32 `json:"blockIndex"`
	// The burned gas (uint64 as string)
	GasBurned string `json:"gasBurned"`
	// The charged gas fee (uint64 as string)
	GasFeeCharged string `json:"gasFeeCharged"`
	L1CommitmentHash string `json:"l1CommitmentHash"`
	NumOffLedgerRequests uint32 `json:"numOffLedgerRequests"`
	NumSuccessfulRequests uint32 `json:"numSuccessfulRequests"`
	PreviousL1CommitmentHash string `json:"previousL1CommitmentHash"`
	Timestamp time.Time `json:"timestamp"`
	// The total L2 base tokens (uint64 as string)
	TotalBaseTokensInL2Accounts string `json:"totalBaseTokensInL2Accounts"`
	TotalRequests uint32 `json:"totalRequests"`
	// The total storage deposit (uint64 as string)
	TotalStorageDeposit string `json:"totalStorageDeposit"`
	TransactionSubEssenceHash string `json:"transactionSubEssenceHash"`
}

// NewBlockInfoResponse instantiates a new BlockInfoResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBlockInfoResponse(anchorTransactionId string, blockIndex uint32, gasBurned string, gasFeeCharged string, l1CommitmentHash string, numOffLedgerRequests uint32, numSuccessfulRequests uint32, previousL1CommitmentHash string, timestamp time.Time, totalBaseTokensInL2Accounts string, totalRequests uint32, totalStorageDeposit string, transactionSubEssenceHash string) *BlockInfoResponse {
	this := BlockInfoResponse{}
	this.AnchorTransactionId = anchorTransactionId
	this.BlockIndex = blockIndex
	this.GasBurned = gasBurned
	this.GasFeeCharged = gasFeeCharged
	this.L1CommitmentHash = l1CommitmentHash
	this.NumOffLedgerRequests = numOffLedgerRequests
	this.NumSuccessfulRequests = numSuccessfulRequests
	this.PreviousL1CommitmentHash = previousL1CommitmentHash
	this.Timestamp = timestamp
	this.TotalBaseTokensInL2Accounts = totalBaseTokensInL2Accounts
	this.TotalRequests = totalRequests
	this.TotalStorageDeposit = totalStorageDeposit
	this.TransactionSubEssenceHash = transactionSubEssenceHash
	return &this
}

// NewBlockInfoResponseWithDefaults instantiates a new BlockInfoResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBlockInfoResponseWithDefaults() *BlockInfoResponse {
	this := BlockInfoResponse{}
	return &this
}

// GetAnchorTransactionId returns the AnchorTransactionId field value
func (o *BlockInfoResponse) GetAnchorTransactionId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AnchorTransactionId
}

// GetAnchorTransactionIdOk returns a tuple with the AnchorTransactionId field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetAnchorTransactionIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AnchorTransactionId, true
}

// SetAnchorTransactionId sets field value
func (o *BlockInfoResponse) SetAnchorTransactionId(v string) {
	o.AnchorTransactionId = v
}

// GetBlockIndex returns the BlockIndex field value
func (o *BlockInfoResponse) GetBlockIndex() uint32 {
	if o == nil {
		var ret uint32
		return ret
	}

	return o.BlockIndex
}

// GetBlockIndexOk returns a tuple with the BlockIndex field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetBlockIndexOk() (*uint32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BlockIndex, true
}

// SetBlockIndex sets field value
func (o *BlockInfoResponse) SetBlockIndex(v uint32) {
	o.BlockIndex = v
}

// GetGasBurned returns the GasBurned field value
func (o *BlockInfoResponse) GetGasBurned() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GasBurned
}

// GetGasBurnedOk returns a tuple with the GasBurned field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetGasBurnedOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GasBurned, true
}

// SetGasBurned sets field value
func (o *BlockInfoResponse) SetGasBurned(v string) {
	o.GasBurned = v
}

// GetGasFeeCharged returns the GasFeeCharged field value
func (o *BlockInfoResponse) GetGasFeeCharged() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GasFeeCharged
}

// GetGasFeeChargedOk returns a tuple with the GasFeeCharged field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetGasFeeChargedOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GasFeeCharged, true
}

// SetGasFeeCharged sets field value
func (o *BlockInfoResponse) SetGasFeeCharged(v string) {
	o.GasFeeCharged = v
}

// GetL1CommitmentHash returns the L1CommitmentHash field value
func (o *BlockInfoResponse) GetL1CommitmentHash() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.L1CommitmentHash
}

// GetL1CommitmentHashOk returns a tuple with the L1CommitmentHash field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetL1CommitmentHashOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.L1CommitmentHash, true
}

// SetL1CommitmentHash sets field value
func (o *BlockInfoResponse) SetL1CommitmentHash(v string) {
	o.L1CommitmentHash = v
}

// GetNumOffLedgerRequests returns the NumOffLedgerRequests field value
func (o *BlockInfoResponse) GetNumOffLedgerRequests() uint32 {
	if o == nil {
		var ret uint32
		return ret
	}

	return o.NumOffLedgerRequests
}

// GetNumOffLedgerRequestsOk returns a tuple with the NumOffLedgerRequests field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetNumOffLedgerRequestsOk() (*uint32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NumOffLedgerRequests, true
}

// SetNumOffLedgerRequests sets field value
func (o *BlockInfoResponse) SetNumOffLedgerRequests(v uint32) {
	o.NumOffLedgerRequests = v
}

// GetNumSuccessfulRequests returns the NumSuccessfulRequests field value
func (o *BlockInfoResponse) GetNumSuccessfulRequests() uint32 {
	if o == nil {
		var ret uint32
		return ret
	}

	return o.NumSuccessfulRequests
}

// GetNumSuccessfulRequestsOk returns a tuple with the NumSuccessfulRequests field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetNumSuccessfulRequestsOk() (*uint32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NumSuccessfulRequests, true
}

// SetNumSuccessfulRequests sets field value
func (o *BlockInfoResponse) SetNumSuccessfulRequests(v uint32) {
	o.NumSuccessfulRequests = v
}

// GetPreviousL1CommitmentHash returns the PreviousL1CommitmentHash field value
func (o *BlockInfoResponse) GetPreviousL1CommitmentHash() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.PreviousL1CommitmentHash
}

// GetPreviousL1CommitmentHashOk returns a tuple with the PreviousL1CommitmentHash field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetPreviousL1CommitmentHashOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PreviousL1CommitmentHash, true
}

// SetPreviousL1CommitmentHash sets field value
func (o *BlockInfoResponse) SetPreviousL1CommitmentHash(v string) {
	o.PreviousL1CommitmentHash = v
}

// GetTimestamp returns the Timestamp field value
func (o *BlockInfoResponse) GetTimestamp() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetTimestampOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Timestamp, true
}

// SetTimestamp sets field value
func (o *BlockInfoResponse) SetTimestamp(v time.Time) {
	o.Timestamp = v
}

// GetTotalBaseTokensInL2Accounts returns the TotalBaseTokensInL2Accounts field value
func (o *BlockInfoResponse) GetTotalBaseTokensInL2Accounts() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TotalBaseTokensInL2Accounts
}

// GetTotalBaseTokensInL2AccountsOk returns a tuple with the TotalBaseTokensInL2Accounts field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetTotalBaseTokensInL2AccountsOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TotalBaseTokensInL2Accounts, true
}

// SetTotalBaseTokensInL2Accounts sets field value
func (o *BlockInfoResponse) SetTotalBaseTokensInL2Accounts(v string) {
	o.TotalBaseTokensInL2Accounts = v
}

// GetTotalRequests returns the TotalRequests field value
func (o *BlockInfoResponse) GetTotalRequests() uint32 {
	if o == nil {
		var ret uint32
		return ret
	}

	return o.TotalRequests
}

// GetTotalRequestsOk returns a tuple with the TotalRequests field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetTotalRequestsOk() (*uint32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TotalRequests, true
}

// SetTotalRequests sets field value
func (o *BlockInfoResponse) SetTotalRequests(v uint32) {
	o.TotalRequests = v
}

// GetTotalStorageDeposit returns the TotalStorageDeposit field value
func (o *BlockInfoResponse) GetTotalStorageDeposit() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TotalStorageDeposit
}

// GetTotalStorageDepositOk returns a tuple with the TotalStorageDeposit field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetTotalStorageDepositOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TotalStorageDeposit, true
}

// SetTotalStorageDeposit sets field value
func (o *BlockInfoResponse) SetTotalStorageDeposit(v string) {
	o.TotalStorageDeposit = v
}

// GetTransactionSubEssenceHash returns the TransactionSubEssenceHash field value
func (o *BlockInfoResponse) GetTransactionSubEssenceHash() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TransactionSubEssenceHash
}

// GetTransactionSubEssenceHashOk returns a tuple with the TransactionSubEssenceHash field value
// and a boolean to check if the value has been set.
func (o *BlockInfoResponse) GetTransactionSubEssenceHashOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TransactionSubEssenceHash, true
}

// SetTransactionSubEssenceHash sets field value
func (o *BlockInfoResponse) SetTransactionSubEssenceHash(v string) {
	o.TransactionSubEssenceHash = v
}

func (o BlockInfoResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BlockInfoResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["anchorTransactionId"] = o.AnchorTransactionId
	toSerialize["blockIndex"] = o.BlockIndex
	toSerialize["gasBurned"] = o.GasBurned
	toSerialize["gasFeeCharged"] = o.GasFeeCharged
	toSerialize["l1CommitmentHash"] = o.L1CommitmentHash
	toSerialize["numOffLedgerRequests"] = o.NumOffLedgerRequests
	toSerialize["numSuccessfulRequests"] = o.NumSuccessfulRequests
	toSerialize["previousL1CommitmentHash"] = o.PreviousL1CommitmentHash
	toSerialize["timestamp"] = o.Timestamp
	toSerialize["totalBaseTokensInL2Accounts"] = o.TotalBaseTokensInL2Accounts
	toSerialize["totalRequests"] = o.TotalRequests
	toSerialize["totalStorageDeposit"] = o.TotalStorageDeposit
	toSerialize["transactionSubEssenceHash"] = o.TransactionSubEssenceHash
	return toSerialize, nil
}

type NullableBlockInfoResponse struct {
	value *BlockInfoResponse
	isSet bool
}

func (v NullableBlockInfoResponse) Get() *BlockInfoResponse {
	return v.value
}

func (v *NullableBlockInfoResponse) Set(val *BlockInfoResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableBlockInfoResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableBlockInfoResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBlockInfoResponse(val *BlockInfoResponse) *NullableBlockInfoResponse {
	return &NullableBlockInfoResponse{value: val, isSet: true}
}

func (v NullableBlockInfoResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBlockInfoResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



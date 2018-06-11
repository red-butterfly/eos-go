package eos

import (
	"encoding/hex"
	"encoding/json"
	"reflect"

	"github.com/red-butterfly/eos-go/ecc"
)

/*
{
  "server_version": "f537bc50",
  "head_block_num": 9,
  "last_irreversible_block_num": 8,
  "last_irreversible_block_id": "00000008f98f0580d7efe7abc60abaaf8a865c9428a4267df30ff7d1937a1084",
  "head_block_id": "00000009ecd0e9fb5719431f4b86f5c9ca1887f6b6f73e5a301aaff740fd6bd3",
  "head_block_time": "2018-05-19T07:47:31",
  "head_block_producer": "eosio",
  "virtual_block_cpu_limit": 100800,
  "virtual_block_net_limit": 1056996,
  "block_cpu_limit": 99900,
  "block_net_limit": 1048576
}

*/

type InfoResp struct {
	ServerVersion            string      `json:"server_version"` // "2cc40a4e"
	ChainID                  SHA256Bytes `json:"chain_id"`
	HeadBlockNum             uint32      `json:"head_block_num"`              // 2465669,
	LastIrreversibleBlockNum uint32      `json:"last_irreversible_block_num"` // 2465655
	LastIrreversibleBlockID  SHA256Bytes `json:"last_irreversible_block_id"`  // "00000008f98f0580d7efe7abc60abaaf8a865c9428a4267df30ff7d1937a1084"
	HeadBlockID              SHA256Bytes `json:"head_block_id"`               // "00259f856bfa142d1d60aff77e70f0c4f3eab30789e9539d2684f9f8758f1b88",
	HeadBlockTime            JSONTime    `json:"head_block_time"`             //  "2018-02-02T04:19:32"
	HeadBlockProducer        AccountName `json:"head_block_producer"`         // "inita"

	VirtualBlockCPULimit uint64 `json:"virtual_block_cpu_limit"`
	VirtualBlockNetLimit uint64 `json:"virtual_block_net_limit"`
	BlockCPULimit        uint64 `json:"block_cpu_limit"`
	BlockNetLimit        uint64 `json:"block_net_limit"`
}

type BlockResp struct {
	SignedBlock
	ID             SHA256Bytes `json:"id"`
	BlockNum       uint32      `json:"block_num"`
	RefBlockPrefix uint32      `json:"ref_block_prefix"`
}

// type BlockTransaction struct {
// 	Status        string            `json:"status"`
// 	CPUUsageUS    int               `json:"cpu_usage_us"`
// 	NetUsageWords int               `json:"net_usage_words"`
// 	Trx           []json.RawMessage `json:"trx"`
// }

type TransactionResp struct {
	TransactionID string `json:"transaction_id"`
	Transaction   struct {
		Signatures            []ecc.Signature `json:"signatures"`
		Compression           CompressionType `json:"compression"`
		PackedContextFreeData HexBytes        `json:"packed_context_free_data"`
		ContextFreeData       []HexBytes      `json:"context_free_data"`
		PackedTransaction     HexBytes        `json:"packed_transaction"`
		Transaction           Transaction     `json:"transaction"`
	} `json:"transaction"`
}

type SequencedTransactionResp struct {
	SeqNum int `json:"seq_num"`
	TransactionResp
}

type TransactionsResp struct {
	Transactions []SequencedTransactionResp
}

type ProducerChange struct {
}

type AccountResp struct {
	AccountName        AccountName          `json:"account_name"`
	Privileged         bool                 `json:"privileged"`
	LastCodeUpdate     JSONTime             `json:"last_code_update"`
	Created            JSONTime             `json:"created"`
	RAMQuota           int64                `json:"ram_quota"`
	RAMUsage           int64                `json:"ram_usage"`
	NetWeight          int64                `json:"net_weight"`
	CPUWeight          int64                `json:"cpu_weight"`
	NetLimit           AccountResourceLimit `json:"net_limit"`
	CPULimit           AccountResourceLimit `json:"cpu_limit"`
	Permissions        []Permission         `json:"permissions"`
	TotalResources     TotalResources       `json:"total_resources"`
	DelegatedBandwidth DelegatedBandwidth   `json:"self_delegated_bandwidth"`
	VoterInfo          VoterInfo            `json:"voter_info"`
}

type CurrencyBalanceResp struct {
	EOSBalance        Asset    `json:"eos_balance"`
	StakedBalance     Asset    `json:"staked_balance"`
	UnstakingBalance  Asset    `json:"unstaking_balance"`
	LastUnstakingTime JSONTime `json:"last_unstaking_time"`
}

type GetTableRowsRequest struct {
	JSON       bool   `json:"json"`
	Scope      string `json:"scope"`
	Code       string `json:"code"`
	Table      string `json:"table"`
	TableKey   string `json:"table_key"`
	LowerBound string `json:"lower_bound"`
	UpperBound string `json:"upper_bount"`
	Limit      uint32 `json:"limit,omitempty"` // defaults to 10 => chain_plugin.hpp:struct get_table_rows_params
}

type GetTableRowsResp struct {
	More bool            `json:"more"`
	Rows json.RawMessage `json:"rows"` // defer loading, as it depends on `JSON` being true/false.
}

func (resp *GetTableRowsResp) JSONToStructs(v interface{}) error {
	return json.Unmarshal(resp.Rows, v)
}

func (resp *GetTableRowsResp) BinaryToStructs(v interface{}) error {
	var rows []string

	err := json.Unmarshal(resp.Rows, &rows)
	if err != nil {
		return err
	}

	outSlice := reflect.ValueOf(v).Elem()
	structType := reflect.TypeOf(v).Elem().Elem()

	for _, row := range rows {
		bin, err := hex.DecodeString(row)
		if err != nil {
			return err
		}

		// access the type of the `Slice`, create a bunch of them..
		newStruct := reflect.New(structType)

		decoder := NewDecoder(bin)
		if err := decoder.Decode(newStruct.Interface()); err != nil {
			return err
		}

		outSlice = reflect.Append(outSlice, reflect.Indirect(newStruct))
	}

	reflect.ValueOf(v).Elem().Set(outSlice)

	return nil
}

type Currency struct {
	Precision uint8
	Name      CurrencyName
}

type GetRequiredKeysResp struct {
	RequiredKeys []ecc.PublicKey `json:"required_keys"`
}

// PushTransactionFullResp unwraps the responses from a successful `push_transaction`.
// FIXME: REVIEW the actual output, things have moved here.
type PushTransactionFullResp struct {
	StatusCode    string
	TransactionID string               `json:"transaction_id"`
	Processed     TransactionProcessed `json:"processed"` // WARN: is an `fc::variant` in server..
}

type TransactionProcessed struct {
	Status               string        `json:"status"`
	ID                   SHA256Bytes   `json:"id"`
	ActionTraces         []ActionTrace `json:"action_traces"`
	DeferredTransactions []string      `json:"deferred_transactions"` // that's not right... dig to find what's there..
}

type ActionTrace struct {
	Receiver AccountName `json:"receiver"`
	// Action     Action       `json:"act"` // FIXME: how do we unpack that ? what's on the other side anyway?
	Console    string       `json:"console"`
	DataAccess []DataAccess `json:"data_access"`
}

type DataAccess struct {
	Type     string      `json:"type"` // "write", "read"?
	Code     AccountName `json:"code"`
	Scope    AccountName `json:"scope"`
	Sequence int         `json:"sequence"`
}

type PushTransactionShortResp struct {
	TransactionID string `json:"transaction_id"`
	Processed     bool   `json:"processed"` // WARN: is an `fc::variant` in server..
}

//

type WalletSignTransactionResp struct {
	// Ignore the rest of the transaction, so the wallet server
	// doesn't forge some transactions on your behalf, and you send it
	// to the network..  ... although.. it's better if you can trust
	// your wallet !

	Signatures []ecc.Signature `json:"signatures"`
}

type MyStruct struct {
	Currency
	Balance uint64
}

// NetConnectionResp
type NetConnectionsResp struct {
	Peer          string           `json:"peer"`
	Connecting    bool             `json:"connecting"`
	Syncing       bool             `json:"syncing"`
	LastHandshake HandshakeMessage `json:"last_handshake"`
}

type NetStatusResp struct {
}

type NetConnectResp string

type NetDisconnectResp string

type ProducersResp struct {
	// TODO: fill this in !
}

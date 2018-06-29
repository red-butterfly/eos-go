package system

import (
	eos "github.com/red-butterfly/eos-go"
)

func NewSellRAM(account eos.AccountName, bytes uint64) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("sellram"),
		Authorization: []eos.PermissionLevel{
			{account, eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(SellRAM{
			Account: account,
			Bytes: bytes,
		}),
	}
	return a
}

// SetRAM represents the hard-coded `setram` action.
type SellRAM struct {
	Account 	eos.AccountName		`json:"account"`
	Bytes  		uint64 				`json:"bytes"`
}

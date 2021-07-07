package model

import "github.com/ONBUFF-IP-TOKEN/event-server/rest_server/controllers/context"

type AuthInfo struct {
	AuthToken  string            `json:"auth_token"`
	ExpireDate int64             `json:"expire_date"`
	WalletAuth context.LoginAuth `json:"wallet_auth"`
}

type EventItemInfo struct {
	Idx                int64
	Name               string
	Serial             string
	TokenId            int64
	TokenUri           string
	Owner              string
	PurchaseTxHash     string
	PurchaseTs         int64
	SubmitStart        int64
	SubmitEnd          int64
	MinAmountForSumbit int64
	Price              int64
	Info               string
}

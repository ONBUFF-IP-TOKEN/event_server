package constant

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
)

const (
	Result_Success                 = 0
	Result_RequireWalletAddress    = 30000
	Result_RequireEmail            = 30002
	Result_RequireTimestamp        = 30004
	Result_NotExistInfo            = 30005
	Result_ExistInfo               = 30006
	Result_NotWinning              = 30007
	Result_RequiredPurchaseNoti    = 30008
	Result_InvalidItemNumber       = 30009
	Result_AlreayPurchase          = 30010
	Result_ClosedEvent             = 30011
	Result_InProgressEvent         = 30012
	Result_InvalidWalletAddress    = 30013
	Result_NotEnoughTokenForSubmit = 30014
	Result_NotSubmitPeriod         = 30015
	Result_PurchaseStep1Err        = 30016
	Result_InvalidPurchaeStep      = 30017

	Result_MaxLengthExceed = 31000

	Result_DBError = 13000

	Result_Auth_RequireMessage    = 20000
	Result_Auth_RequireSign       = 20001
	Result_Auth_InvalidLoginInfo  = 20002
	Result_Auth_DontEncryptJwt    = 20003
	Result_Auth_InvalidJwt        = 20004
	Result_Auth_InvalidWalletType = 20005
)

var resultCodeText = map[int]string{
	Result_Success:                 "success",
	Result_RequireWalletAddress:    "Wallet address is required",
	Result_RequireEmail:            "Email is required",
	Result_RequireTimestamp:        "Timestamp is required",
	Result_NotExistInfo:            "Not exist info",
	Result_ExistInfo:               "Already submit",
	Result_NotWinning:              "You not winner",
	Result_MaxLengthExceed:         "Max length exceed",
	Result_InvalidItemNumber:       "Invalid Item number",
	Result_AlreayPurchase:          "Already purchase",
	Result_ClosedEvent:             "Closed Event",
	Result_InProgressEvent:         "In Progress Event",
	Result_InvalidWalletAddress:    "Invalid Wallet Address",
	Result_NotEnoughTokenForSubmit: "Not enough token for submit",
	Result_NotSubmitPeriod:         "There is no submit period",
	Result_PurchaseStep1Err:        "Step1 must be done first",
	Result_InvalidPurchaeStep:      "Invalid purchase step",

	Result_RequiredPurchaseNoti: "Purchase hash is required",

	Result_DBError: "Internal DB error",

	Result_Auth_RequireMessage:    "Message is required",
	Result_Auth_RequireSign:       "Sign info is required",
	Result_Auth_InvalidLoginInfo:  "Invalid login info",
	Result_Auth_DontEncryptJwt:    "Auth token create fail",
	Result_Auth_InvalidJwt:        "Invalid jwt token",
	Result_Auth_InvalidWalletType: "Invalid wallet type",
}

func ResultCodeText(code int) string {
	return resultCodeText[code]
}

func MakeResponse(code int) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Return = code
	resp.Message = resultCodeText[code]
	return resp
}

type OnbuffBaseResponse struct {
	Return  int         `json:"return"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

func (o *OnbuffBaseResponse) Success() {
	o.Return = Result_Success
	o.Message = resultCodeText[Result_Success]
}

func (o *OnbuffBaseResponse) SetResult(ret int) {
	o.Return = ret
	o.Message = resultCodeText[ret]
}

func MakeOnbuffBaseResponse(code int) *OnbuffBaseResponse {
	resp := new(OnbuffBaseResponse)
	resp.Return = code
	resp.Message = resultCodeText[code]
	return resp
}

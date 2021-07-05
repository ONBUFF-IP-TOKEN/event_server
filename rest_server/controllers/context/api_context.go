package context

import (
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/constant"
	"github.com/labstack/echo"
)

const (
	Wallet_type_metamask = "metamask"
)

// page info
type PageInfo struct {
	PageOffset int64 `query:"page_offset" validate:"required"`
	PageSize   int64 `query:"page_size" validate:"required"`
}

// page response
type PageInfoResponse struct {
	PageOffset int64 `json:"page_offset"`
	PageSize   int64 `json:"page_size"`
	TotalSize  int64 `json:"total_size"`
}

///////////// API ///////////////////////////////
// login 관련 정보
type LoginAuth struct {
	WalletAddr string `json:"wallet_address" validate:"required"`
	Message    string `json:"message" validate:"required"`
	Sign       string `json:"sign" validate:"required"`
}
type LoginParam struct {
	WalletType string    `json:"wallet_type" validate:"required"`
	WalletAuth LoginAuth `json:"wallet_auth" validate:"required"`
}

func NewLoginParam() *LoginParam {
	return new(LoginParam)
}

func (o *LoginParam) CheckValidate() *constant.OnbuffBaseResponse {
	if len(o.WalletType) == 0 && (Wallet_type_metamask != o.WalletType) {
		return constant.MakeOnbuffBaseResponse(constant.Result_Auth_InvalidWalletType)
	}
	if len(o.WalletAuth.WalletAddr) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireWalletAddress)
	}
	if len(o.WalletAuth.Message) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_Auth_RequireMessage)
	}
	if len(o.WalletAuth.Sign) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_Auth_RequireSign)
	}
	return nil
}

type LoginResponse struct {
	AuthToken  string `json:"auth_token" validate:"required"`
	ExpireDate int64  `json:"expire_date" validate:"required"`
}

/////////////////////////

// 응모하기
type Submit struct {
	Idx         int64  `json:"idx,omitempty"`
	WalletAddr  string `json:"wallet_address" validate:"required"`
	ItemNum     int64  `json:"item_number" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Ts          int64  `json:"ts" validate:"required"`
	Ret         string `json:"ret,omitempty"`
	SubmitCnt   int64  `json:"submit_cnt"`
	LastBalance int64  `json:"balance"`
}

func NewSubmit() *Submit {
	return new(Submit)
}

func (o *Submit) CheckValidate(ctx *IPBlockServerContext) *constant.OnbuffBaseResponse {
	if o.WalletAddr != ctx.WalletAddr() {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidWalletAddress)
	}
	if len(o.WalletAddr) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireWalletAddress)
	}
	if o.ItemNum <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidItemNumber)
	}
	if len(o.Email) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireEmail)
	}
	if len(o.Email) >= 28 {
		return constant.MakeOnbuffBaseResponse(constant.Result_MaxLengthExceed)
	}
	return nil
}

/////////////////////////

// 응모 여부 확인
type SubmitVerify struct {
	WalletAddr string `json:"wallet_address" validate:"required" query:"wallet_address"`
	ItemNum    int64  `json:"item_number" validate:"required" query:"item_number"`
}

func NewSubmitVerify() *SubmitVerify {
	return new(SubmitVerify)
}

func (o *SubmitVerify) CheckValidate(ctx *IPBlockServerContext) *constant.OnbuffBaseResponse {
	if o.WalletAddr != ctx.WalletAddr() {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidWalletAddress)
	}
	if len(o.WalletAddr) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireWalletAddress)
	}
	if o.ItemNum <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidItemNumber)
	}

	return nil
}

/////////////////////////

// 응모 당첨 확인
type SubmitResult struct {
	WalletAddr string `json:"wallet_address" validate:"required" query:"wallet_address"`
	ItemNum    int64  `json:"item_number" validate:"required" query:"item_number"`
}

func NewSubmitResult() *SubmitResult {
	return new(SubmitResult)
}

func (o *SubmitResult) CheckValidate(ctx *IPBlockServerContext) *constant.OnbuffBaseResponse {
	if o.WalletAddr != ctx.WalletAddr() {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidWalletAddress)
	}
	if len(o.WalletAddr) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireWalletAddress)
	}
	if o.ItemNum <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidItemNumber)
	}

	return nil
}

/////////////////////////

// 당첨자 조회
type SubmitWinner struct {
	ItemNum int64 `json:"item_number" validate:"required" query:"item_number"`
}

func NewSubmitWinner() *SubmitWinner {
	return new(SubmitWinner)
}

func (o *SubmitWinner) CheckValidate() *constant.OnbuffBaseResponse {
	if o.ItemNum <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidItemNumber)
	}
	return nil
}

type SubmitWinnerResponse struct {
	WalletAddr string `json:"wallet_address" validate:"required" query:"wallet_address"`
}

/////////////////////////

// 구매 정보 전달
type PurchaseNoti struct {
	WalletAddr          string `json:"wallet_address" validate:"required" query:"wallet_address"`
	ItemNum             int64  `json:"item_number" validate:"required" query:"item_number"`
	PurchaseTxHash      string `json:"purchase_tx_hash" validate:"required"`
	ShippingAddr        string `json:"shipping_address" validate:"required"`
	PhoneNum            string `json:"phone_number" validate:"required"`
	VerifyCheckComplete string `json:"verify_check,omitempty"`
}

func NewPurchaseNoti() *PurchaseNoti {
	return new(PurchaseNoti)
}

func (o *PurchaseNoti) CheckValidate(ctx *IPBlockServerContext) *constant.OnbuffBaseResponse {
	if o.WalletAddr != ctx.WalletAddr() {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidWalletAddress)
	}
	if len(o.WalletAddr) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireWalletAddress)
	}
	if len(o.PurchaseTxHash) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequiredPurchaseNoti)
	}
	if o.ItemNum <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidItemNumber)
	}
	return nil
}

/////////////////////////

/////////////////////////

// 응모자 리스트
type SubmitList struct {
	ItemNum int64 `json:"item_number" validate:"required" query:"item_number"`
}

func NewSubmitList() *SubmitList {
	return new(SubmitList)
}

func (o *SubmitList) CheckValidate(ctx *echo.Context) *constant.OnbuffBaseResponse {
	if o.ItemNum <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidItemNumber)
	}
	return nil
}

/////////////////////////

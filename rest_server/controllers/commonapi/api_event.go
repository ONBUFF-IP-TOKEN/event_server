package commonapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/constant"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/token"
	"github.com/labstack/echo"
)

func GetEventDuplicate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)
	params := context.NewSubmitVerify()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	// 유효한 item_number check
	if _, err := CheckExistItem(c, params.ItemNum); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	resp := new(constant.OnbuffBaseResponse)
	info, err := model.GetDB().GetEventInfo(params.WalletAddr)
	if err != nil {
		resp.SetResult(constant.Result_DBError)
	} else {
		if len(info.WalletAddr) == 0 && len(info.Email) == 0 {
			//등록한 기록이 없다.
			resp.SetResult(constant.Result_NotExistInfo)
		} else {
			//이미 등록한 기록이 있음
			resp.Success()
			resp.Value = info
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func PutEventSubmit(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)
	params := context.NewSubmit()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	// 유효한 item_number check
	itemInfo, errItem := GetExistItem(c, params.ItemNum)
	if errItem != nil {
		return c.JSON(http.StatusOK, errItem)
	}

	resp := new(constant.OnbuffBaseResponse)
	info, err := model.GetDB().GetEventInfo(params.WalletAddr)
	if err != nil {
		resp.SetResult(constant.Result_DBError)
	} else {
		curT := datetime.GetTS2MilliSec()
		if len(info.WalletAddr) == 0 && len(info.Email) == 0 {
			//등록한 기록이 없으니 등록 진행
			params.Ts = curT
			params.SubmitCnt = 1
			// 코인 balance 가져와서 50개 이상인지 check 후 저장
			balance, _ := token.GetToken().Tokens[token.Token_onit].Onit_GetBalanceOf(params.WalletAddr)
			params.LastBalance = balance
			if balance < itemInfo.MinAmountForSumbit {
				// onit 이 최소 보유량보다 적으면 에러 리턴
				resp.SetResult(constant.Result_NotEnoughTokenForSubmit)
			} else {
				// db 저장
				if _, err := model.GetDB().PutEventSubmit(params); err != nil {
					resp.SetResult(constant.Result_DBError)
				} else {
					resp.Success()
					resp.Value = params
				}
			}
		} else {
			//이미 등록한 기록이 있음 등록한지 날짜가 바뀌었다면 재응모 가능
			//날짜가 바뀌지 않았다면 재응모 불가
			lastTime := time.Unix(0, info.Ts*int64(time.Millisecond))
			curTime := time.Unix(0, curT*int64(time.Millisecond))
			if lastTime.Year() == curTime.Year() &&
				lastTime.Month() == curTime.Month() &&
				lastTime.Day() == curTime.Day() {
				//if curT < info.Ts+24*3600*1000 { // 응모하고 만 하루 지난것으로 체크할것인지 예비코드
				resp.SetResult(constant.Result_ExistInfo)
				resp.Value = info
			} else {
				//재 응모 가능
				params.Ts = curT
				params.SubmitCnt = info.SubmitCnt + 1
				// 코인 balance 가져와서 저장
				balance, _ := token.GetToken().Tokens[token.Token_onit].Onit_GetBalanceOf(info.WalletAddr)
				params.LastBalance = balance
				if balance < itemInfo.MinAmountForSumbit {
					// onit 이 최소 보유량보다 적으면 에러 리턴
					resp.SetResult(constant.Result_NotEnoughTokenForSubmit)
				} else {
					// db 저장
					if _, err := model.GetDB().UpdateEventSubmit(params); err != nil {
						resp.SetResult(constant.Result_DBError)
					} else {
						resp.Success()
						resp.Value = params
					}
				}
			}

		}
	}

	return c.JSON(http.StatusOK, resp)
}

func GetEventResult(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)
	params := context.NewSubmitResult()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	// 유효한 item_number check
	if _, err := CheckExistItem(c, params.ItemNum); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	resp := new(constant.OnbuffBaseResponse)
	info, err := model.GetDB().GetEventInfo(params.WalletAddr)
	if err != nil {
		resp.SetResult(constant.Result_DBError)
	} else {
		if len(info.WalletAddr) == 0 && len(info.Email) == 0 {
			//등록한 기록이 없다.
			resp.SetResult(constant.Result_NotExistInfo)
		} else {
			//이미 등록한 기록이 있고 당첨 확인
			if strings.ToUpper(info.Ret) == "OK" {
				//당첨
				resp.Success()
				resp.Value = info
			} else if len(info.Ret) == 0 || strings.ToUpper(info.Ret) != "OK" {
				//꽝
				resp.SetResult(constant.Result_NotWinning)
			}
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func GetEventWinner(c echo.Context) error {
	params := context.NewSubmitWinner()
	if err := c.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	resp := new(constant.OnbuffBaseResponse)

	// 유효한 item_number check
	itemInfo, err := CheckExistItem(c, params.ItemNum)
	if err != nil {
		if len(itemInfo.Owner) != 0 || len(itemInfo.PurchaseTxHash) != 0 {
			//이미 구매 완료된 이벤트
			value := &context.SubmitWinnerResponse{
				WalletAddr: itemInfo.Owner,
			}
			resp.Success()
			resp.Value = value
		} else {
			//아직 진행 중인 이벤트
			resp.SetResult(constant.Result_NotExistInfo)
		}
	} else {
		//아직 진행 중인 이벤트
		resp.SetResult(constant.Result_InProgressEvent)
	}

	return c.JSON(http.StatusOK, resp)
}

func PostEventPurchaseNoti(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)
	params := context.NewPurchaseNoti()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	// 유효한 item_number check
	itemInfo, errItem := CheckExistItem(c, params.ItemNum)
	if errItem != nil {
		return c.JSON(http.StatusOK, errItem)
	}

	resp := new(constant.OnbuffBaseResponse)
	// 응모한 정보 수집
	info, err := model.GetDB().GetEventInfo(params.WalletAddr)
	if err != nil {
		resp.SetResult(constant.Result_DBError)
	} else {
		if len(info.WalletAddr) == 0 && len(info.Email) == 0 {
			//응모한 기록이 없다.
			resp.SetResult(constant.Result_NotExistInfo)
		} else {
			//응모한 기록이 있고 당첨 확인
			if strings.ToUpper(info.Ret) == "OK" {
				// 이미 구매 완료 했는지 확인
				if itemInfo.Owner == params.WalletAddr {
					resp.SetResult(constant.Result_AlreayPurchase)
				} else {
					// 당첨자 기록 파일로그로 남기기
					file, _ := json.MarshalIndent(params, "", " ")
					_ = ioutil.WriteFile("ok.json", file, 0644)
					// 당첨자라면 메타마스크 tx hash 확인 진행
					token.GetToken().Tokens[token.Token_onit].CheckTransferResponse(params)
					resp.Success()
				}
			} else if len(info.Ret) == 0 || strings.ToUpper(info.Ret) != "OK" {
				//꽝
				resp.SetResult(constant.Result_NotWinning)
			}
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func GetLatestSubmitList(c echo.Context) error {
	params := context.NewSubmitList()
	if err := c.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(&c); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	// 유효한 item_number check
	itemInfo, errItem := GetExistItem(c, params.ItemNum)
	if errItem != nil {
		return c.JSON(http.StatusOK, errItem)
	}

	resp := new(constant.OnbuffBaseResponse)
	submits, err := model.GetDB().GetLatestSubmitList(itemInfo.Idx)
	if err != nil {
		resp.SetResult(constant.Result_DBError)
	} else {
		resp.Success()
		submitLst := context.ResSubmitList{}
		submitLst.List = submits
		resp.Value = submitLst
	}

	return c.JSON(http.StatusOK, resp)
}

func CheckExistItem(c echo.Context, ItemNum int64) (*model.EventItemInfo, *constant.OnbuffBaseResponse) {
	item, err := model.GetDB().GetEventItem(ItemNum)
	if err != nil {
		return item, constant.MakeOnbuffBaseResponse(constant.Result_DBError)
	} else {
		if len(item.Name) == 0 || len(item.Serial) == 0 {
			return item, constant.MakeOnbuffBaseResponse(constant.Result_InvalidItemNumber)
		}

		if len(item.Owner) != 0 || len(item.PurchaseTxHash) != 0 {
			//이미 구매 완료된 이벤트
			return item, constant.MakeOnbuffBaseResponse(constant.Result_ClosedEvent)
		}
	}
	return item, nil
}

func GetExistItem(c echo.Context, ItemNum int64) (*model.EventItemInfo, *constant.OnbuffBaseResponse) {
	item, err := model.GetDB().GetEventItem(ItemNum)
	if err != nil {
		return item, constant.MakeOnbuffBaseResponse(constant.Result_DBError)
	} else {
		if len(item.Name) == 0 || len(item.Serial) == 0 {
			return item, constant.MakeOnbuffBaseResponse(constant.Result_InvalidItemNumber)
		}

		if datetime.GetTS2Sec() < item.SubmitStart || datetime.GetTS2Sec() > item.SubmitEnd {
			return item, constant.MakeOnbuffBaseResponse(constant.Result_NotSubmitPeriod)
		}
	}
	return item, nil
}

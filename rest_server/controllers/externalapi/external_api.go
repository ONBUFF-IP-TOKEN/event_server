package externalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/constant"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/controllers/context"
	"github.com/labstack/echo"
)

type ExternalAPI struct {
	base.BaseController

	conf    *config.ServerConfig
	apiConf *baseconf.APIServer
	echo    *echo.Echo
}

func PreCheck(c echo.Context) base.PreCheckResponse {
	conf := config.GetInstance()
	if err := base.SetContext(c, &conf.Config, context.NewIPBlockServerContext); err != nil {
		log.Error(err)
		return base.PreCheckResponse{
			IsSucceed: false,
		}
	}

	//log.Debug(c.Request().Header["Authorization"])
	// auth token 검증
	walletAddr, isValid := auth.GetIAuth().IsValidAuthToken(c.Request().Header["Authorization"][0][7:])
	if conf.Auth.AuthEnable && !isValid {
		// auth token 오류 리턴
		res := constant.MakeOnbuffBaseResponse(constant.Result_Auth_InvalidJwt)

		return base.PreCheckResponse{
			IsSucceed: false,
			Response:  res,
		}
	}
	base.GetContext(c).(*context.IPBlockServerContext).SetWalletAddr(*walletAddr)

	return base.PreCheckResponse{
		IsSucceed: true,
	}
}

func (o *ExternalAPI) Init(e *echo.Echo) error {
	o.echo = e
	o.BaseController.PreCheck = PreCheck

	if err := o.MapRoutes(o, e, o.apiConf.Routes); err != nil {
		return err
	}

	return nil
}

func (o *ExternalAPI) GetConfig() *baseconf.APIServer {
	o.conf = config.GetInstance()
	o.apiConf = &o.conf.APIServers[1]
	return o.apiConf
}

func NewAPI() *ExternalAPI {
	return &ExternalAPI{}
}

func (o *ExternalAPI) GetHealthCheck(c echo.Context) error {
	return commonapi.GetHealthCheck(c)
}

func (o *ExternalAPI) GetVersion(c echo.Context) error {
	return commonapi.GetVersion(c, o.BaseController.MaxVersion)
}

func (o *ExternalAPI) PostEventLogin(c echo.Context) error {
	return commonapi.PostLogin(c)
}

func (o *ExternalAPI) GetEventDuplicate(c echo.Context) error {
	return commonapi.GetEventDuplicate(c)
}

func (o *ExternalAPI) PutEventSubmit(c echo.Context) error {
	return commonapi.PutEventSubmit(c)
}

func (o *ExternalAPI) GetEventResult(c echo.Context) error {
	return commonapi.GetEventResult(c)
}

func (o *ExternalAPI) GetEventWinner(c echo.Context) error {
	return commonapi.GetEventWinner(c)
}

func (o *ExternalAPI) PostEventPurchaseNoti(c echo.Context) error {
	return commonapi.PostEventPurchaseNoti(c)
}

func (o *ExternalAPI) GetEventPurchase(c echo.Context) error {
	return nil
}

func (o *ExternalAPI) GetLatestSubmitList(c echo.Context) error {
	return commonapi.GetLatestSubmitList(c)
}

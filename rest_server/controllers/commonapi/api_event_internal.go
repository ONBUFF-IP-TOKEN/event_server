package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/constant"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/model"
	"github.com/labstack/echo"
)

func PostResetWinner(c echo.Context) error {
	resp := new(constant.OnbuffBaseResponse)
	resp.Success()
	return c.JSON(http.StatusOK, resp)
}

func PostResetPurchase(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)
	params := context.NewResetPurchase()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	model.GetDB().PostResetPurchase(params)

	resp := new(constant.OnbuffBaseResponse)
	resp.Success()
	return c.JSON(http.StatusOK, resp)
}

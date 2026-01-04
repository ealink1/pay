package handler

import (
	"net/http"
	"pay/model"

	"github.com/gin-gonic/gin"
)

type AlipayNotifyRequest struct {
	AppId       string `json:"app_id"`
	TradeNo     string `json:"trade_no"`
	OutTradeNo  string `json:"out_trade_no"`
	TotalAmount string `json:"total_amount"`
	TradeStatus string `json:"trade_status"`
	NotifyTime  string `json:"notify_time"`
	NotifyType  string `json:"notify_type"`
	NotifyId    string `json:"notify_id"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	BuyerId     string `json:"buyer_id"`
}

func AlipayNotify(c *gin.Context) {
	var params map[string]string
	if err := c.ShouldBind(&params); err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	sign := params["sign"]
	delete(params, "sign")
	delete(params, "sign_type")

	if !verifySign(params, sign) {
		c.String(http.StatusOK, "fail")
		return
	}

	order, exists := model.Store.GetByOutTradeNo(params["out_trade_no"])
	if !exists {
		c.String(http.StatusOK, "fail")
		return
	}

	if params["trade_status"] == "TRADE_SUCCESS" || params["trade_status"] == "TRADE_FINISHED" {
		model.Store.UpdateStatus(order.ID, model.OrderStatusPaid, params["trade_no"])
	}

	c.String(http.StatusOK, "success")
}

func verifySign(params map[string]string, sign string) bool {
	// 在这里需要使用支付宝公钥进行验证
	// 为了简化，暂时返回 true，实际生产环境中需要正确验证签名
	return true
}

package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"pay/logging"
	"pay/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	logger := logging.FromGin(c)
	if err := c.Request.ParseForm(); err != nil {
		logger.Warn("alipay_notify_parse_form_failed", zap.String("error", err.Error()))
		c.String(http.StatusOK, "fail")
		return
	}

	params := make(map[string]string, len(c.Request.Form))
	for k, vs := range c.Request.Form {
		if len(vs) == 0 {
			continue
		}
		params[k] = strings.Join(vs, ",")
	}

	sign := params["sign"]

	headersJSON := ""
	if b, err := json.Marshal(c.Request.Header); err == nil {
		headersJSON = string(b)
	}
	paramsJSON := ""
	if b, err := json.Marshal(params); err == nil {
		paramsJSON = string(b)
	}

	if err := alipayClient.VerifySign(params, sign); err != nil {
		writeCallbackLogAsync(model.CallbackLog{
			Provider:    "alipay",
			Path:        c.FullPath(),
			Method:      c.Request.Method,
			RemoteIP:    c.ClientIP(),
			TraceID:     logging.TraceIDFromGin(c),
			AppID:       params["app_id"],
			OutTradeNo:  params["out_trade_no"],
			TradeNo:     params["trade_no"],
			TradeStatus: params["trade_status"],
			NotifyID:    params["notify_id"],
			Sign:        sign,
			VerifyOK:    false,
			VerifyError: err.Error(),
			ParamsJSON:  paramsJSON,
			HeadersJSON: headersJSON,
			ReceivedAt:  time.Now(),
		})
		logger.Warn("alipay_notify_verify_failed", zap.String("error", err.Error()))
		c.String(http.StatusOK, "fail")
		return
	}

	outTradeNo := params["out_trade_no"]
	tradeStatus := params["trade_status"]
	tradeNo := params["trade_no"]

	order, exists := model.Store.GetByOutTradeNo(outTradeNo)
	if !exists {
		writeCallbackLogAsync(model.CallbackLog{
			Provider:    "alipay",
			Path:        c.FullPath(),
			Method:      c.Request.Method,
			RemoteIP:    c.ClientIP(),
			TraceID:     logging.TraceIDFromGin(c),
			AppID:       params["app_id"],
			OutTradeNo:  outTradeNo,
			TradeNo:     tradeNo,
			TradeStatus: tradeStatus,
			NotifyID:    params["notify_id"],
			Sign:        sign,
			VerifyOK:    true,
			VerifyError: "order not found",
			ParamsJSON:  paramsJSON,
			HeadersJSON: headersJSON,
			ReceivedAt:  time.Now(),
		})
		logger.Warn("alipay_notify_order_not_found", zap.String("out_trade_no", outTradeNo))
		c.String(http.StatusOK, "fail")
		return
	}

	nextStatus := order.Status
	switch tradeStatus {
	case "WAIT_BUYER_PAY":
		nextStatus = model.OrderStatusPending
	case "TRADE_SUCCESS", "TRADE_FINISHED":
		nextStatus = model.OrderStatusPaid
	case "TRADE_CLOSED":
		nextStatus = model.OrderStatusClosed
	}

	if err := model.Store.UpdateStatus(order.ID, nextStatus, tradeNo); err != nil {
		writeCallbackLogAsync(model.CallbackLog{
			Provider:    "alipay",
			Path:        c.FullPath(),
			Method:      c.Request.Method,
			RemoteIP:    c.ClientIP(),
			TraceID:     logging.TraceIDFromGin(c),
			AppID:       params["app_id"],
			OutTradeNo:  outTradeNo,
			TradeNo:     tradeNo,
			TradeStatus: tradeStatus,
			NotifyID:    params["notify_id"],
			Sign:        sign,
			VerifyOK:    true,
			ParamsJSON:  paramsJSON,
			HeadersJSON: headersJSON,
			ReceivedAt:  time.Now(),
		})
		logger.Error("alipay_notify_update_failed", zap.String("order_id", order.ID), zap.String("error", err.Error()))
		c.String(http.StatusOK, "fail")
		return
	}

	writeCallbackLogAsync(model.CallbackLog{
		Provider:    "alipay",
		Path:        c.FullPath(),
		Method:      c.Request.Method,
		RemoteIP:    c.ClientIP(),
		TraceID:     logging.TraceIDFromGin(c),
		AppID:       params["app_id"],
		OutTradeNo:  outTradeNo,
		TradeNo:     tradeNo,
		TradeStatus: tradeStatus,
		NotifyID:    params["notify_id"],
		Sign:        sign,
		VerifyOK:    true,
		ParamsJSON:  paramsJSON,
		HeadersJSON: headersJSON,
		ReceivedAt:  time.Now(),
	})
	logger.Info("alipay_notify_ok", zap.String("order_id", order.ID), zap.String("out_trade_no", outTradeNo), zap.String("trade_no", tradeNo), zap.String("trade_status", tradeStatus), zap.String("status", string(nextStatus)))
	c.String(http.StatusOK, "success")
}

func writeCallbackLogAsync(log model.CallbackLog) {
	go func() {
		if err := model.CallbackLogs.Create(&log); err != nil {
			logging.L().Error("callback_log_create_failed", zap.Error(err))
		}
	}()
}

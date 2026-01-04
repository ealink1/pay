package ealipay

import (
	"net/url"
	"sort"
	"strings"
)

type AppPayRequest struct {
	OutTradeNo     string `json:"out_trade_no"`
	TotalAmount    string `json:"total_amount"`
	Subject        string `json:"subject"`
	Body           string `json:"body,omitempty"`
	TimeoutExpress string `json:"timeout_express,omitempty"`
	ProductCode    string `json:"product_code"`
}

type AppPayResponse struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	SubCode  string `json:"sub_code,omitempty"`
	SubMsg   string `json:"sub_msg,omitempty"`
	OrderStr string `json:"order_str"`
}

func (c *AlipayClient) AppPay(req *AppPayRequest) (string, error) {
	if req.ProductCode == "" {
		req.ProductCode = "QUICK_WAP_WAY"
	}

	if req.TimeoutExpress == "" {
		req.TimeoutExpress = "30m"
	}

	bizContent, err := c.buildBizContent(req)
	if err != nil {
		return "", err
	}

	params := c.buildCommonParams("alipay.trade.app.pay", bizContent)
	params["timestamp"] = getCurrentTimestamp()

	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteString("&")
		}
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(params[k])
	}

	sign, err := c.sign(builder.String())
	if err != nil {
		return "", err
	}
	params["sign"] = sign

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	return values.Encode(), nil
}

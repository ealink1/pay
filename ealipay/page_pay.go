package ealipay

type PagePayRequest struct {
	OutTradeNo     string  `json:"out_trade_no"`
	TotalAmount    string  `json:"total_amount"`
	Subject        string  `json:"subject"`
	Body           string  `json:"body,omitempty"`
	TimeoutExpress string  `json:"timeout_express,omitempty"`
	ProductCode    string  `json:"product_code"`
}

type PagePayResponse struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code,omitempty"`
	SubMsg  string `json:"sub_msg,omitempty"`
}

func (c *AlipayClient) PagePay(req *PagePayRequest) (string, error) {
	if req.ProductCode == "" {
		req.ProductCode = "FAST_INSTANT_TRADE_PAY"
	}

	if req.TimeoutExpress == "" {
		req.TimeoutExpress = "30m"
	}

	url, err := c.buildUrl("alipay.trade.page.pay", req)
	if err != nil {
		return "", err
	}

	return url, nil
}

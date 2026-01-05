package ealipay

type WapPayRequest struct {
	OutTradeNo     string `json:"out_trade_no"`
	TotalAmount    string `json:"total_amount"`
	Subject        string `json:"subject"`
	Body           string `json:"body,omitempty"`
	TimeoutExpress string `json:"timeout_express,omitempty"`
	ProductCode    string `json:"product_code"`
}

func (c *AlipayClient) WapPay(req *WapPayRequest) (string, error) {
	if req.ProductCode == "" {
		req.ProductCode = "QUICK_WAP_WAY"
	}

	if req.TimeoutExpress == "" {
		req.TimeoutExpress = "30m"
	}

	url, err := c.buildUrl("alipay.trade.wap.pay", req)
	if err != nil {
		return "", err
	}

	return url, nil
}

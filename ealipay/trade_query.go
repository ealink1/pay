package ealipay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type TradeQueryRequest struct {
	OutTradeNo string `json:"out_trade_no,omitempty"`
	TradeNo    string `json:"trade_no,omitempty"`
}

type TradeQueryResponse struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	SubCode     string `json:"sub_code,omitempty"`
	SubMsg      string `json:"sub_msg,omitempty"`
	TradeStatus string `json:"trade_status,omitempty"`
	OutTradeNo  string `json:"out_trade_no,omitempty"`
	TradeNo     string `json:"trade_no,omitempty"`
	TotalAmount string `json:"total_amount,omitempty"`
}

func (c *AlipayClient) TradeQuery(req *TradeQueryRequest) (*TradeQueryResponse, error) {
	if req == nil || (req.OutTradeNo == "" && req.TradeNo == "") {
		return nil, fmt.Errorf("out_trade_no or trade_no is required")
	}

	bizContent, err := c.buildBizContent(req)
	if err != nil {
		return nil, err
	}

	params := c.buildCommonParams("alipay.trade.query", bizContent)
	params["timestamp"] = getCurrentTimestamp()

	sign, err := c.generateSign(params)
	if err != nil {
		return nil, err
	}
	params["sign"] = sign

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.GatewayUrl, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("alipay http status %d: %s", resp.StatusCode, string(body))
	}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	raw := envelope["alipay_trade_query_response"]
	if len(raw) == 0 {
		return nil, fmt.Errorf("missing alipay_trade_query_response: %s", string(body))
	}

	var out TradeQueryResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.Code != "10000" {
		if out.SubCode != "" || out.SubMsg != "" {
			return &out, fmt.Errorf("alipay error: %s %s (%s %s)", out.Code, out.Msg, out.SubCode, out.SubMsg)
		}
		return &out, fmt.Errorf("alipay error: %s %s", out.Code, out.Msg)
	}

	return &out, nil
}

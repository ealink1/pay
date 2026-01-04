package handler

import (
	"net/http"
	"pay/ealipay"
	"pay/model"

	"github.com/gin-gonic/gin"
)

type CreateAppOrderRequest struct {
	TotalAmount string `json:"total_amount" binding:"required"`
	Subject     string `json:"subject" binding:"required"`
	Body        string `json:"body"`
}

type CreateAppOrderResponse struct {
	OrderID  string `json:"order_id"`
	OrderStr string `json:"order_str"`
}

func CreateAppOrder(c *gin.Context) {
	var req CreateAppOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := &model.Order{
		OutTradeNo:  generateOutTradeNo(),
		TotalAmount: req.TotalAmount,
		Subject:     req.Subject,
		Body:        req.Body,
	}

	if err := model.Store.Create(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	payReq := &ealipay.AppPayRequest{
		OutTradeNo:  order.OutTradeNo,
		TotalAmount: order.TotalAmount,
		Subject:     order.Subject,
		Body:        order.Body,
	}

	orderStr, err := alipayClient.AppPay(payReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成支付参数失败"})
		return
	}

	model.Store.UpdateStatus(order.ID, model.OrderStatusPending, "")

	c.JSON(http.StatusOK, CreateAppOrderResponse{
		OrderID:  order.ID,
		OrderStr: orderStr,
	})
}

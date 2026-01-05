package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"pay/ealipay"
	"pay/logging"
	"pay/model"
	"time"

	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

type CreateOrderRequest struct {
	TotalAmount string `json:"total_amount" binding:"required"`
	Subject     string `json:"subject" binding:"required"`
	Body        string `json:"body"`
}

type CreateOrderResponse struct {
	OrderID   string `json:"order_id"`
	QrCode    string `json:"qr_code"`
	QrCodeURL string `json:"qr_code_url"`
}

type UpdateOrderStatusRequest struct {
	Status  string `json:"status" binding:"required"`
	TradeNo string `json:"trade_no"`
}

var alipayClient *ealipay.AlipayClient

func InitAlipayClient(config *ealipay.Config) error {
	client, err := ealipay.NewClient(config)
	if err != nil {
		return err
	}
	alipayClient = client
	return nil
}

func CreateOrder(c *gin.Context) {
	logger := logging.FromGin(c)
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("create_order_bad_request", zap.String("error", err.Error()))
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
		logger.Error("create_order_store_failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	payReq := &ealipay.PagePayRequest{
		OutTradeNo:  order.OutTradeNo,
		TotalAmount: order.TotalAmount,
		Subject:     order.Subject,
		Body:        order.Body,
	}

	payUrl, err := alipayClient.PagePay(payReq)
	if err != nil {
		logger.Error("create_order_page_pay_failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成支付链接失败"})
		return
	}

	qrCodeData, err := qrcode.Encode(payUrl, qrcode.Medium, 256)
	if err != nil {
		logger.Error("create_order_qrcode_failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成二维码失败"})
		return
	}

	qrCodeBase64 := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(qrCodeData))

	order.QrCode = qrCodeBase64
	model.Store.UpdateStatus(order.ID, model.OrderStatusPending, "")

	c.JSON(http.StatusOK, CreateOrderResponse{
		OrderID:   order.ID,
		QrCode:    qrCodeBase64,
		QrCodeURL: payUrl,
	})
}

func GetOrder(c *gin.Context) {
	orderID := c.Param("id")

	order, exists := model.Store.GetByID(orderID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func ListOrders(c *gin.Context) {
	orders := model.Store.List()
	c.JSON(http.StatusOK, orders)
}

func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, exists := model.Store.GetByID(orderID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	var status model.OrderStatus
	switch req.Status {
	case "pending":
		status = model.OrderStatusPending
	case "paid":
		status = model.OrderStatusPaid
	case "failed":
		status = model.OrderStatusFailed
	case "closed":
		status = model.OrderStatusClosed
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单状态"})
		return
	}

	if err := model.Store.UpdateStatus(order.ID, status, req.TradeNo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新订单状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "订单状态更新成功"})
}

func SyncOrderStatus(c *gin.Context) {
	logger := logging.FromGin(c)
	orderID := c.Param("id")

	order, exists := model.Store.GetByID(orderID)
	if !exists {
		logger.Warn("sync_order_not_found", zap.String("order_id", orderID))
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}
	if order.OutTradeNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单缺少 out_trade_no"})
		return
	}

	resp, err := alipayClient.TradeQuery(&ealipay.TradeQueryRequest{OutTradeNo: order.OutTradeNo})
	if err != nil {
		logger.Error("sync_order_trade_query_failed", zap.String("order_id", order.ID), zap.String("error", err.Error()))
		c.JSON(http.StatusBadGateway, gin.H{"error": "查询支付宝订单失败", "detail": err.Error()})
		return
	}

	nextStatus := order.Status
	switch resp.TradeStatus {
	case "WAIT_BUYER_PAY":
		nextStatus = model.OrderStatusPending
	case "TRADE_SUCCESS", "TRADE_FINISHED":
		nextStatus = model.OrderStatusPaid
	case "TRADE_CLOSED":
		nextStatus = model.OrderStatusClosed
	case "TRADE_FAIL":
		nextStatus = model.OrderStatusFailed
	}

	if err := model.Store.UpdateStatus(order.ID, nextStatus, resp.TradeNo); err != nil {
		logger.Error("sync_order_update_failed", zap.String("order_id", order.ID), zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新订单状态失败"})
		return
	}

	updated, _ := model.Store.GetByID(order.ID)
	logger.Info("sync_order_ok", zap.String("order_id", order.ID), zap.String("alipay_trade_status", resp.TradeStatus), zap.String("status", string(nextStatus)))
	c.JSON(http.StatusOK, gin.H{
		"order":               updated,
		"alipay_trade_status": resp.TradeStatus,
	})
}

func generateOutTradeNo() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000_000_000))
	if err != nil {
		return time.Now().Format("20060102150405") + "000000000000"
	}
	return fmt.Sprintf("%s%012d", time.Now().Format("20060102150405"), n.Int64())
}

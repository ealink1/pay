package main

import (
	"log"
	"pay/ealipay"
	"pay/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	config := &ealipay.Config{
		AppId:           "9021000150609532",
		PrivateKey:      "MIIEowIBAAKCAQEAs/p0goI6iFufCYFTKDMQ8wIVmgXVVmnOuuX9F8qWtbb0Pc/5K8nFCUe+l09wM5dGosfvtY4jAc5BjhOQaw1vAGK2P6GcHlBVtvgxAUwmV0q2jkPiRvjsZ0oEY8t69k/1mxtrd8pg8FTIgZW85r+Gz8kcGOapi4peu3cu+xLk45T7rQWHxRAfh32xp1QHWBouIjaZOyikRH4Gku4ZnCueTqeDmrqXOLvuGvZ92bLntJd+iDueWSbbTIb+fn4rvZgCyA3CDLHtttNvRRzRXOZtPvmWIJQoWDZ1/gEv51fjdfezN6WQBppG6axbEKbiS867nm0vGzJkti8UNl8raGhOoQIDAQABAoIBAFudxctdoZAiG54KEBupixo42Gg0SfoYGF05kBGZVgigXkpM4QkyR7PGqrV5gaMxgYqBfnuMJDPaG7LIML7d8sBef2l6ye8Ac/GU+9UuP2I2LSHUWo5ITobxvbRTM3/JCjxvw9AR3DDa58pXP/ayTlzdggkG+g2HXVvOesLiRlO2e1bMKodhElETektkTZ27FmeZfBCYcXKXkKVYMMfATZ/Ap6eIGzye1gpWZONQHqBgpHtCZzIClsgQBqCvZY1ESAX6BGKnTIgtA+dJnVBgWhv9qIgC/RvE8CVM9IgdayDlAUcI9JyzgXWYcD9RkebByi5hRUXjMOQUV0VDDWsW3sECgYEA40fz0fDPIps2ICiM6WG8XrnTHdvt6y3IGAB1RPBWvud4aTkQ685QoBGciQzUO2e44I6KGT97hi/sMunoLgblPPIgdIpZiHodDWBbI/1eDNblCHGtunx9IVtkol3dYxToVr5RQqZodCwrHXNhpF67Sc1Vnli5mYV0hBrEGm4cdEkCgYEAyrhc0ThEcYXMv8C8WdL4YwMwig0QXT4lNo0vpXid+P6Ddr1pGQcnCvXB3PfdeAAranaJYCvtLnZbYAOAGKx4fUWV7ZNxzyEu07y8DjbnRcypDK8wDoPqCyCuS6x4fTcXhk9t80JfvbDSNNRz9Q/8VQQt+5kxG70GI0iQJDrSV5kCgYEApCSb53xV7DVSUslWc1q9s1/bI85pNpc60nLKPr6gt4DuSngHS3YWbnQprCUSxdB0CeGHxRI/ALtth5u8rjkWp/xqCiC85r7ian2zdPuQSA+PG5kWEf/EUynxNP47XEqGPdd3Un5iI7yeasegthgghP2BnzmO2VwzuRCnnjr129kCgYAb4EZDLu2afr+tDp/X6j7lvqaKFUnOyKDtY3TN2ExA1R7W0S0GmAkyZKEH9b2qprtRpIM3ilLPNM9T4KdYvT7EWzFGviPES9fYnfduLPaYjpAggmalWFZyuUe+eDUJYu4FNh70eIgZ2ZrOUPixFkWomy6HjoVGPzP83hmUIdKS4QKBgGc6E5duy3wijKsBMYjGMJHnf1DHIgwM/imfwy5N79HfQOVIgM/ebg7kjIdTl9MW9v597O0nuxMi3IZzT9fe059tSdlZWLiGQNkLsCz97TJEbLNJoX3VDGn8v97mbiDppP3S9kmNy32b4kHK/TBAyw5emSULA94EgbRfadmLwt55",
		AlipayPublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3sLSLP88Gnl168l+I1ukFwYfnxJyW27BF1vnuWg88hynfyAe9/8aUUUKVfyx0gM4G7dQGRSxD4WR7du70pJNp7q3JwON0q7GDgEo8usEbL4xyyrTzN+rSIqEF4Zc6XxeOOGANijjX8EoOrHZK99Szg/n0d9m/Z6bVgDFjkqfX0HGc/EsqyGS9uFYd+G3Ignw2Ywkp+ZW2CthQ1ESxJKaDoALA+IHLVsd1I9mkkVLtmZJq4u/CTo1ogLHRx3nbNN9qhVAW6WaJyOo5+JZfo7daJyxh5jblZd+PbV83Fy1k4KA6dbyeSmbyEvhX8xsZ1gNw23INmZA0ZY4SR3K2dqO8wIDAQAB",
		IsSandbox:       true,
	}

	if err := handler.InitAlipayClient(config); err != nil {
		log.Fatalf("初始化支付宝客户端失败: %v", err)
	}

	r := gin.Default()

	r.Static("/static", "./front/static")
	r.LoadHTMLGlob("front/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	api := r.Group("/api")
	{
		api.POST("/orders", handler.CreateOrder)
		api.POST("/app-orders", handler.CreateAppOrder)
		api.GET("/orders", handler.ListOrders)
		api.GET("/orders/:id", handler.GetOrder)
		api.PUT("/orders/:id/status", handler.UpdateOrderStatus) // 调试用的订单状态更新接口
		api.POST("/alipay/notify", handler.AlipayNotify)
	}

	log.Println("服务器启动在 http://localhost:3423")
	if err := r.Run(":3423"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

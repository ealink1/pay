package main

import (
	"log"
	"pay/config"
	"pay/ealipay"
	"pay/handler"
	"pay/logging"
	"pay/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	appCfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	logger, err := logging.Init(appCfg.Log)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	mysqlCfg, err := config.LoadMySQLConfig("config.yaml")
	if err != nil {
		logger.Fatal("mysql_config_load_failed", zap.Error(err))
	}

	db, err := gorm.Open(mysql.Open(mysqlCfg.DSN), &gorm.Config{})
	if err != nil {
		logger.Fatal("mysql_connect_failed", zap.Error(err))
	}
	if err := model.InitGormOrderStore(db); err != nil {
		logger.Fatal("order_store_init_failed", zap.Error(err))
	}

	sandbox := appCfg.Pay.AlipaySandbox
	if sandbox.AppId == "" && sandbox.PrivateKey == "" && sandbox.AlipayPublicKey == "" {
		sandbox.AppId = "9021000150609532"
		sandbox.PrivateKey = "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCGKuP0LKRmr4TZqt2dlOf0sehYfbAtZY002A6Bq8kUtSxc4frHlqlZmjsjdw0gJH9ZZi+SIvAG9XmzeKEL3AOk/q3ezuz6KP0l22NlZMpyEwjnh7VUWBfWXxO+mpzFBWnJtjIF4Esax97E3FlRQFXtYDhFPPJyJeYXnFikgG7TUqTYvPCxZNmx0moXIjjNdmDFSorfpflymLaNE3v5EVHbi9aiyk/FOOO1OhzBSme6viuym4Uc8AkCBCbYIMlzB/glrKZIjzGoNb7E8rN1vz6JWQX6fuP6Fq8M0AkvSAYjsuhbWrKnSyIrC1i6U4LDQ9uTjFToqkJcea0MJYHZWXAlAgMBAAECggEAZ7nDEGRvGYA62jlyKkET7yaX9cn+KaqoN6GN3Yxc4iiLSqfexO1isgY+EFYbDK2K0yfgQT/Hh+nCFBF/mHaZTrci3u3lYiXMSLdLKfl5ViYHLVDKzJFqpG5PCn3oE53ywmKcW9Si2+qH/HRKjTmK9QD9n/HVkpBgSgKyuUMd6zuA6wfSE2hm3BIzhCO8lvu9xqHYoBfJYirxD5v4JRMNv2ZusVdgE+OcklXedt9UzOXt7MsgVsmGZ/Rt6Er0PrHY12fcccu1udlFoT3ubyFRFAjhkAtWEpi5ql6cf8Fa5zgRkRjMqYQ6+xpVVxDACdlDiJlTu3aIfDi0A1EH0ERigQKBgQD0TcXOLLhO2CYFDWYfkuuHnStb/7/wKBVFjF4h7e9Av81qhyCY5Mw37cWTklZEbufnxQdwDpgU0/OrCXAQhFM14WCqs+Koo3bFuigulg/JhtsNkd8v7WkJjUF/WZ1uqm1Iy/58jV3tBN1k0nymauNNIQvNig1MbjAWW1kcszAgxQKBgQCMl0UPfntuBdfJ2c31+mXhLvRv03jfXBrd27P5Pt7kuBC1XJ0fBPyWPsZwg0w3Z4gY9QP/I2dK4+SgtTnWS3yKetImcuBY1D8A1KZQUkZGrPY+kxI8KreO8By0O/OweeN2ZUkNl8686SheDbu7mKehIJ7P+K+LXNZc6f+oen5H4QKBgQDSU1qmo+2RQ5mH4/8106EewfsgW1B9i6S0maI5B8VhMz/AJNG1j9UZmYTuBaBrjifta715hbb8x3USnS9zqNiSnJRColfS49hPZnNNmDfDQmy4hAtoEbbKWGg5IYfeTK+FasqPpI1mjzejo2tZQtCqCHdG30GPuZWAyegwQzx+GQKBgExEz+k843bnYo4VQ19azKQhleeIYH1DeSu8QWFIkyCfHilVKcOnL+POAFcPU2yHFNT9LoLd0O5WvTPVvJ+dad2yDYlgLobh9Z/cvLC8QXWb5SZDINRVFClN5zR7hZLKPPSAs+XU4gmnrwd/CcYWZXHKwXzvW0QORBg5tUDP2uvhAoGALkMuBw2LlvqC2HoLBil7EZtbeS3UmCxxNNhJGLJ+0++so2IHFPU++iBUO+TJPgv0pba61cibOd1L2yEBtS+CXwDgr3kJm0P+8T7S7qIX+hzErt7v9vQr626fCG0kdlYQFeIidad/Qcq7YhyQuIiSht9TwZTmVvi9e/f8iOPAWlc="
		sandbox.AlipayPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3sLSLP88Gnl168l+I1ukFwYfnxJyW27BF1vnuWg88hynfyAe9/8aUUUKVfyx0gM4G7dQGRSxD4WR7du70pJNp7q3JwON0q7GDgEo8usEbL4xyyrTzN+rSIqEF4Zc6XxeOOGANijjX8EoOrHZK99Szg/n0d9m/Z6bVgDFjkqfX0HGc/EsqyGS9uFYd+G3Ignw2Ywkp+ZW2CthQ1ESxJKaDoALA+IHLVsd1I9mkkVLtmZJq4u/CTo1ogLHRx3nbNN9qhVAW6WaJyOo5+JZfo7daJyxh5jblZd+PbV83Fy1k4KA6dbyeSmbyEvhX8xsZ1gNw23INmZA0ZY4SR3K2dqO8wIDAQAB"
	}

	alipayCfg := &ealipay.Config{
		AppId:           sandbox.AppId,
		PrivateKey:      sandbox.PrivateKey,
		AlipayPublicKey: sandbox.AlipayPublicKey,
		IsSandbox:       true,
		NotifyURL:       sandbox.NotifyURL,
		ReturnURL:       sandbox.ReturnURL,
	}

	if err := handler.InitAlipayClient(alipayCfg); err != nil {
		logger.Fatal("alipay_client_init_failed", zap.Error(err))
	}

	r := gin.New()
	r.Use(logging.Middleware(logger, appCfg.Trace))
	r.Use(gin.Recovery())

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
		api.POST("/orders/:id/sync", handler.SyncOrderStatus)
		api.POST("/alipay/notify", handler.AlipayNotify)
		api.POST("/alipay/sandbox/notify", handler.AlipayNotify)
	}

	logger.Info("server_start", zap.String("addr", "http://localhost:3423"))
	if err := r.Run(":3423"); err != nil {
		logger.Fatal("server_run_failed", zap.Error(err))
	}
}

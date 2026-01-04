package ealipay

import "fmt"

func ExampleUsage() {
	config := &Config{
		AppId:           "9021000150609532",
		PrivateKey:      "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCGKuP0LKRmr4TZqt2dlOf0sehYfbAtZY002A6Bq8kUtSxc4frHlqlZmjsjdw0gJH9ZZi+SIvAG9XmzeKEL3AOk/q3ezuz6KP0l22NlZMpyEwjnh7VUWBfWXxO+mpzFBWnJtjIF4Esax97E3FlRQFXtYDhFPPJyJeYXnFikgG7TUqTYvPCxZNmx0moXIjjNdmDFSorfpflymLaNE3v5EVHbi9aiyk/FOOO1OhzBSme6viuym4Uc8AkCBCbYIMlzB/glrKZIjzGoNb7E8rN1vz6JWQX6fuP6Fq8M0AkvSAYjsuhbWrKnSyIrC1i6U4LDQ9uTjFToqkJcea0MJYHZWXAlAgMBAAECggEAZ7nDEGRvGYA62jlyKkET7yaX9cn+KaqoN6GN3Yxc4iiLSqfexO1isgY+EFYbDK2K0yfgQT/Hh+nCFBF/mHaZTrci3u3lYiXMSLdLKfl5ViYHLVDKzJFqpG5PCn3oE53ywmKcW9Si2+qH/HRKjTmK9QD9n/HVkpBgSgKyuUMd6zuA6wfSE2hm3BIzhCO8lvu9xqHYoBfJYirxD5v4JRMNv2ZusVdgE+OcklXedt9UzOXt7MsgVsmGZ/Rt6Er0PrHY12fcccu1udlFoT3ubyFRFAjhkAtWEpi5ql6cf8Fa5zgRkRjMqYQ6+xpVVxDACdlDiJlTu3aIfDi0A1EH0ERigQKBgQD0TcXOLLhO2CYFDWYfkuuHnStb/7/wKBVFjF4h7e9Av81qhyCY5Mw37cWTklZEbufnxQdwDpgU0/OrCXAQhFM14WCqs+Koo3bFuigulg/JhtsNkd8v7WkJjUF/WZ1uqm1Iy/58jV3tBN1k0nymauNNIQvNig1MbjAWW1kcszAgxQKBgQCMl0UPfntuBdfJ2c31+mXhLvRv03jfXBrd27P5Pt7kuBC1XJ0fBPyWPsZwg0w3Z4gY9QP/I2dK4+SgtTnWS3yKetImcuBY1D8A1KZQUkZGrPY+kxI8KreO8By0O/OweeN2ZUkNl8686SheDbu7mKehIJ7P+K+LXNZc6f+oen5H4QKBgQDSU1qmo+2RQ5mH4/8106EewfsgW1B9i6S0maI5B8VhMz/AJNG1j9UZmYTuBaBrjifta715hbb8x3USnS9zqNiSnJRColfS49hPZnNNmDfDQmy4hAtoEbbKWGg5IYfeTK+FasqPpI1mjzejo2tZQtCqCHdG30GPuZWAyegwQzx+GQKBgExEz+k843bnYo4VQ19azKQhleeIYH1DeSu8QWFIkyCfHilVKcOnL+POAFcPU2yHFNT9LoLd0O5WvTPVvJ+dad2yDYlgLobh9Z/cvLC8QXWb5SZDINRVFClN5zR7hZLKPPSAs+XU4gmnrwd/CcYWZXHKwXzvW0QORBg5tUDP2uvhAoGALkMuBw2LlvqC2HoLBil7EZtbeS3UmCxxNNhJGLJ+0++so2IHFPU++iBUO+TJPgv0pba61cibOd1L2yEBtS+CXwDgr3kJm0P+8T7S7qIX+hzErt7v9vQr626fCG0kdlYQFeIidad/Qcq7YhyQuIiSht9TwZTmVvi9e/f8iOPAWlc=",
		AlipayPublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3sLSLP88Gnl168l+I1ukFwYfnxJyW27BF1vnuWg88hynfyAe9/8aUUUKVfyx0gM4G7dQGRSxD4WR7du70pJNp7q3JwON0q7GDgEo8usEbL4xyyrTzN+rSIqEF4Zc6XxeOOGANijjX8EoOrHZK99Szg/n0d9m/Z6bVgDFjkqfX0HGc/EsqyGS9uFYd+G3Ignw2Ywkp+ZW2CthQ1ESxJKaDoALA+IHLVsd1I9mkkVLtmZJq4u/CTo1ogLHRx3nbNN9qhVAW6WaJyOo5+JZfo7daJyxh5jblZd+PbV83Fy1k4KA6dbyeSmbyEvhX8xsZ1gNw23INmZA0ZY4SR3K2dqO8wIDAQAB",
		IsSandbox:       true,
	}

	client, err := NewClient(config)
	if err != nil {
		panic(err)
	}

	req := &PagePayRequest{
		OutTradeNo:  "ORDER_20240104_001",
		TotalAmount: "100.00",
		Subject:     "商品名称",
		Body:        "商品描述",
		ProductCode: "FAST_INSTANT_TRADE_PAY",
	}

	payUrl, err := client.PagePay(req)
	if err != nil {
		panic(err)
	}

	fmt.Println("支付链接:", payUrl)
}

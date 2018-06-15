// Create by Yale 2018/6/13 16:52
package wechatgz

import (
	"paysdk"
	"paysdk/wechat"
	"testing"
)

func TestWeChatGZ_Order(t *testing.T) {

	gz := WeChat()
	gz.AddInfo(&wechat.WeChatPayInfo{
		AppID:     "wx195270080d9f5d3a",
		ApiKey:    "0smxoPugcsWjbL3uVvUNZIiOzZIQ7pRz",
		MchID:     "1243432302",
		NotifyUrl: "https://950103cc.ngrok.io",
	})
	gz.Order(&paysdk.OrderParams{
		BusTradeNo: "aaa1231231145411",
		Title:      "test111",
		Detail:     "撒啊",
		Price:      0.01,
	})

}

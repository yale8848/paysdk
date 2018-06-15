// Create by Yale 2018/6/13 10:35
package wechatapp

import (
	"fmt"
	"github.com/yale8848/paysdk"
	"github.com/yale8848/paysdk/util"
	"github.com/yale8848/paysdk/wechat"
	"strings"
	"time"
)

var weChatGZ *WeChatApp

func init() {
	weChatGZ = &WeChatApp{}
}
func WeChat() *WeChatApp {
	return weChatGZ
}

type WeChatApp struct {
	wechat.WeChat
}

func (w *WeChatApp) Order(params *paysdk.OrderParams) error {

	wechatInfo := w.GetInfo(params.AppName)

	err := w.CheckOrderParams(params)
	if err != nil {
		return err
	}

	return w.UnifyOrder(wechatInfo, params, func(m *map[string]string) {
		(*m)["trade_type"] = "APP"
	}, func(xmlRes *wechat.WeChatResResult) (*wechat.WeChatRetResult, error) {
		c := make(map[string]string)
		c["appid"] = wechatInfo.AppID
		c["noncestr"] = util.GetRandomStr()
		c["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
		c["package"] = "Sign=WXPay"
		c["partnerid"] = wechatInfo.MchID
		c["prepayid"] = xmlRes.PrepayID

		sign, err := w.GenSign(wechatInfo.ApiKey, c)
		if err != nil {
			return nil, err
		}
		sign = strings.ToUpper(sign)

		wrr := &wechat.WeChatRetResult{}
		wrr.Package = c["package"]
		wrr.AppID = wechatInfo.AppID
		wrr.Sign = sign
		wrr.NonceStr = c["noncestr"]
		wrr.PartnerId = c["partnerid"]
		wrr.TimeStamp = c["timestamp"]
		wrr.PrepayId = c["prepayid"]

		return wrr, nil
	})
}

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

var weChatApp *WeChatApp

func init() {
	weChatApp = &WeChatApp{}
}
func WeChat() *WeChatApp {
	return weChatApp
}

type WeChatApp struct {
	wechat.WeChat
}
type Pay struct {
	wechat.WeChatPay
}

func (p *Pay) Order(params *paysdk.OrderParams) error {
	w := weChatApp
	wechatInfo := w.GetInfo(params.AppName)

	err := w.CheckOrderParams(params)
	if err != nil {
		return err
	}

	res, ret, err := w.UnifyOrder(wechatInfo, params, func(m *map[string]string) {
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
		wrr.MchID = c["partnerid"]
		wrr.TimeStamp = c["timestamp"]
		wrr.PrepayId = c["prepayid"]
		wrr.SignType = "MD5"

		return wrr, nil
	})
	if err != nil {
		return err
	}
	p.WeChatRetResult = ret
	p.WeChatResResult = res
	return nil
}
func (w *WeChatApp) Pay() *Pay {
	return &Pay{}
}

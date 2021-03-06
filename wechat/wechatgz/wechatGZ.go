// Create by Yale 2018/6/13 10:35
package wechatgz

import (
	"errors"
	"fmt"
	"github.com/yale8848/paysdk"
	"github.com/yale8848/paysdk/util"
	"github.com/yale8848/paysdk/wechat"
	"strings"
	"time"
)

var weChatGZ *WeChatGZ

func init() {
	weChatGZ = &WeChatGZ{}
}
func WeChat() *WeChatGZ {
	return weChatGZ
}

type WeChatGZ struct {
	wechat.WeChat
}

func (w *WeChatGZ) CheckOrderParams(params *paysdk.OrderParams) error {
	err := w.WeChat.CheckOrderParams(params)
	if err != nil {
		return err
	}

	if len(params.Openid) == 0 {
		return errors.New("参数有误")
	}
	return nil

}

type Pay struct {
	wechat.WeChatPay
}

func (p *Pay) Order(params *paysdk.OrderParams) error {
	w := weChatGZ
	wechatInfo := w.GetInfo(params.AppName)

	err := w.CheckOrderParams(params)
	if err != nil {
		return err
	}
	res, ret, err := w.UnifyOrder(wechatInfo, params, func(m *map[string]string) {

		(*m)["trade_type"] = "JSAPI"
		(*m)["openid"] = params.Openid
	}, func(xmlRes *wechat.WeChatResResult) (*wechat.WeChatRetResult, error) {
		c := make(map[string]string)
		c["appId"] = wechatInfo.AppID
		c["nonceStr"] = util.GetRandomStr()
		c["timeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
		c["package"] = "prepay_id=" + xmlRes.PrepayID
		c["signType"] = "MD5"

		sign, err := w.GenSign(wechatInfo.ApiKey, c)
		if err != nil {
			return nil, err
		}
		sign = strings.ToUpper(sign)

		wrr := &wechat.WeChatRetResult{}

		wrr.AppID = wechatInfo.AppID
		wrr.NonceStr = c["nonceStr"]
		wrr.TimeStamp = c["timeStamp"]
		wrr.Package = c["package"]
		wrr.SignType = c["signType"]
		wrr.Sign = sign
		wrr.MchID = wechatInfo.MchID
		wrr.PrepayId = xmlRes.PrepayID
		return wrr, nil
	})
	if err != nil {
		return err
	}
	p.WeChatRetResult = ret
	p.WeChatResResult = res
	return nil

}
func (w *WeChatGZ) Pay() *Pay {
	return &Pay{}
}

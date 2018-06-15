// Create by Yale 2018/6/13 14:10
package wechat

import (
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"paysdk"
	"paysdk/common"
	"paysdk/util"
	"sort"
	"strings"
)

const ORDER_URL = "https://api.mch.weixin.qq.com/pay/unifiedorder"

type AddMapFun func(m *map[string]string)
type AddMapFunSign func(xmlRes *WeChatResResult) (*WeChatRetResult, error)

type WeChat struct {
	infoMap map[string]*WeChatPayInfo

	WeChatResResult *WeChatResResult
	WeChatRetResult *WeChatRetResult
}
type WeChatPayInfo struct {
	AppID  string
	MchID  string
	ApiKey string

	NotifyUrl string
	PayUrl    string
}
type WeChatRetResult struct {
	AppID     string
	PartnerId string
	Package   string
	PrepayId  string
	NonceStr  string
	TimeStamp string
	Sign      string
	SignType  string
	CodeUrl   string
}

type WeChatResResult struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`

	AppID      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	TradeType string `xml:"trade_type"`
	PrepayID  string `xml:"prepay_id"`
	CodeURL   string `xml:"code_url"`
}

// WechatBaseResult 基本信息
type WechatBaseResult struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// WechatReturnData 返回通用数据
type WechatReturnData struct {
	AppID      string `xml:"appid,emitempty"`
	MchID      string `xml:"mch_id,emitempty"`
	DeviceInfo string `xml:"device_info,emitempty"`
	NonceStr   string `xml:"nonce_str,emitempty"`
	Sign       string `xml:"sign,emitempty"`
	ResultCode string `xml:"result_code,emitempty"`
	ErrCode    string `xml:"err_code,emitempty"`
	ErrCodeDes string `xml:"err_code_des,emitempty"`
}

// WechatResultData 结果通用数据
type WechatResultData struct {
	OpenID        string `xml:"openid,emitempty"`
	IsSubscribe   string `xml:"is_subscribe,emitempty"`
	TradeType     string `xml:"trade_type,emitempty"`
	BankType      string `xml:"bank_type,emitempty"`
	FeeType       string `xml:"fee_type,emitempty"`
	TotalFee      int64  `xml:"total_fee,emitempty"`
	CashFeeType   string `xml:"cash_fee_type,emitempty"`
	CashFee       int64  `xml:"cash_fee,emitempty"`
	CouponFee     int64  `xml:"coupon_fee,emitempty"`
	CouponCount   int64  `xml:"coupon_count,emitempty"`
	TransactionID string `xml:"transaction_id,emitempty"`
	OutTradeNO    string `xml:"out_trade_no,emitempty"`
	Attach        string `xml:"attach,emitempty"`
	TimeEnd       string `xml:"time_end,emitempty"`
}

type WeChatPayResult struct {
	WechatBaseResult
	WechatReturnData
	WechatResultData
}

func (w *WeChat) init() {
	if w.infoMap == nil {
		w.infoMap = make(map[string]*WeChatPayInfo)
	}
}

func (w *WeChat) AddInfo(payInfo *WeChatPayInfo) {
	w.init()
	w.infoMap[""] = payInfo
}
func (w *WeChat) AddInfoGyAppName(appName string, payInfo *WeChatPayInfo) {
	w.init()
	w.infoMap[appName] = payInfo
}
func (w *WeChat) GetInfo(appName string) *WeChatPayInfo {
	if w.infoMap == nil {
		return nil
	}
	return w.infoMap[appName]
}
func getDefaultOrderUrl(url string) string {

	if len(url) == 0 {
		return ORDER_URL
	}
	return url
}
func (w *WeChat) CheckOrderParams(params *paysdk.OrderParams) error {

	if params != nil {
		bLen := len(params.BusTradeNo)
		sLen := len(params.Title)
		dLen := len(params.Detail)
		if bLen > 0 && bLen <= 32 && sLen > 0 && sLen <= 128 && dLen >= 0 && dLen <= 6000 && params.Price > 0 {
			return nil
		}
	}
	return errors.New("参数有误")

}
func (w *WeChat) UnifyOrder(wechatInfo *WeChatPayInfo, params *paysdk.OrderParams, request AddMapFun, response AddMapFunSign,
) error {
	var m = make(map[string]string)
	m["appid"] = wechatInfo.AppID
	m["mch_id"] = wechatInfo.MchID
	m["nonce_str"] = util.GetRandomStr()
	m["body"] = params.Title

	if len(params.Detail) > 0 {
		m["detail"] = params.Detail
	}
	m["out_trade_no"] = params.BusTradeNo
	m["total_fee"] = util.GetPriceFenStr(params.Price)
	m["spbill_create_ip"] = util.LocalIP()
	m["notify_url"] = wechatInfo.NotifyUrl
	request(&m)

	sign, err := w.GenSign(wechatInfo.ApiKey, m)
	if err != nil {
		return err
	}
	m["sign"] = sign

	xmlStr := util.Map2XML(m)
	re, err := common.HttpPost(getDefaultOrderUrl(wechatInfo.PayUrl), "text/xml:charset=UTF-8", xmlStr)
	if err != nil {
		return err
	}

	var xmlRes WeChatResResult

	err = xml.Unmarshal(re, &xmlRes)
	if err != nil {
		return err
	}
	if xmlRes.ReturnCode != "SUCCESS" {
		return errors.New(xmlRes.ReturnMsg)
	}

	if xmlRes.ResultCode != "SUCCESS" {
		return errors.New(xmlRes.ErrCodeDes)
	}

	c, err := response(&xmlRes)
	if err != nil {
		return err
	}
	w.WeChatRetResult = c
	w.WeChatResResult = &xmlRes
	return nil
}
func CheckSign(data string, sign string, key string) error {
	signData := data + "&key=" + key
	c := md5.New()
	_, err := c.Write([]byte(signData))
	if err != nil {
		return err
	}
	signOut := fmt.Sprintf("%x", c.Sum(nil))
	if strings.ToUpper(sign) == strings.ToUpper(signOut) {
		return nil
	}
	return errors.New("签名交易错误")
}
func (w *WeChat) GenSign(key string, m map[string]string) (string, error) {
	delete(m, "sign")
	delete(m, "key")
	var signData []string
	for k, v := range m {
		if v != "" {
			signData = append(signData, fmt.Sprintf("%s=%s", k, v))
		}
	}
	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + key
	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return "", err
	}
	signByte := c.Sum(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", signByte), nil
}

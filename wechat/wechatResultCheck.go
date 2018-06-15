// Create by Yale 2018/6/14 15:03
package wechat

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"paysdk/util"
	"sort"
	"strings"
)

type WeChatResultCheck struct {
	returnCode string
	returnMsg  string
	result     []byte
}

func (w *WeChatResultCheck) ReturnFail(msg string) string {
	formatStr := `<xml><return_code><![CDATA[%s]]></return_code>
                  <return_msg>![CDATA[%s]]</return_msg></xml>`
	returnBody := fmt.Sprintf(formatStr, "FAIL", msg)

	return returnBody
}
func (w *WeChatResultCheck) ReturnBody() string {
	formatStr := `<xml><return_code><![CDATA[%s]]></return_code>
                  <return_msg>![CDATA[%s]]</return_msg></xml>`
	returnBody := fmt.Sprintf(formatStr, w.returnCode, w.returnMsg)

	return returnBody
}
func (w *WeChatResultCheck) CheckSign(key string) error {
	m, err := util.XmlToMap(w.result)
	if err != nil {
		w.returnMsg = "参数错误"
		w.returnCode = "FAIL"
		return err
	}
	var signData []string
	for k, v := range m {
		if k == "sign" {
			continue
		}
		signData = append(signData, fmt.Sprintf("%v=%v", k, v))
	}
	sort.Strings(signData)
	signData2 := strings.Join(signData, "&")
	err = CheckSign(signData2, m["sign"], key)
	if err != nil {
		w.returnMsg = err.Error()
		w.returnCode = "FAIL"
		return err
	}
	return nil
}
func (w *WeChatResultCheck) Result(r *http.Request) (*WeChatPayResult, error) {
	w.returnCode = "FAIL"
	if r == nil {
		w.returnMsg = "request nil"
		return nil, errors.New("request nil")
	}
	var reXML WeChatPayResult
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.returnCode = "FAIL"
		w.returnMsg = "BodyError"
		return nil, err
	}
	err = xml.Unmarshal(body, &reXML)
	if err != nil {
		w.returnMsg = "参数错误"
		w.returnCode = "FAIL"
		return nil, err
	}

	if reXML.ReturnCode != "SUCCESS" {
		w.returnCode = "FAIL"
		return &reXML, errors.New(reXML.ReturnCode)
	}
	w.result = body
	w.returnCode = "SUCCESS"
	return &reXML, nil

}

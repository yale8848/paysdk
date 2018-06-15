// Create by Yale 2018/6/15 15:05
package test

import (
	"fmt"
	"github.com/yale8848/paysdk"
	"github.com/yale8848/paysdk/wechat"
	"github.com/yale8848/paysdk/wechat/wechatapp"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func initPay() {
	wapp := wechatapp.WeChat()
	wapp.AddInfo(&wechat.WeChatPayInfo{
		AppID:     "",
		ApiKey:    "",
		MchID:     "",
		NotifyUrl: "",
	})
	wapp.AddInfoGyAppName("app1", &wechat.WeChatPayInfo{
		AppID:     "",
		ApiKey:    "",
		MchID:     "",
		NotifyUrl: "",
	})
}

func callback(w http.ResponseWriter, r *http.Request) {

	result := wechat.WeChatResultCheck{}

	_, err := result.Result(r)
	if err != nil {
		w.Write([]byte(result.ReturnBody()))
		return
	}

	err = result.CheckSign("")
	if err != nil {
		w.Write([]byte(result.ReturnBody()))
		return
	}

	w.Write([]byte(result.ReturnBody()))

}

func start() {

	initPay()

	http.HandleFunc("/wechat/getParam", func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		busTradeNo := r.PostFormValue("busTradeNo")
		price := r.PostFormValue("price")
		subject := r.PostFormValue("subject")
		body := r.PostFormValue("body")

		res := NewResponse()
		result, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("%s\n", result)
		if len(busTradeNo) == 0 || len(price) == 0 || len(subject) == 0 || len(body) == 0 {
			w.Write(res.ToJsonBytes())
			return
		}
		var s paysdk.OrderParams
		p, _ := strconv.ParseFloat(price, 64)
		s.BusTradeNo = busTradeNo
		s.Price = p
		s.Detail = body
		s.Title = subject
		s.AppName = "app1"

		app := wechatapp.WeChat()
		err := app.Order(&s)
		if err != nil {
			res.Message = err.Error()
			w.Write(res.ToJsonBytes())
			return
		}
		type appResult struct {
			AppID     string `json:"appid"`
			PartnerId string `json:"partnerId"`
			Package   string `json:"packageValue"`
			PrepayId  string `json:"prepayId"`
			NonceStr  string `json:"nonceStr"`
			TimeStamp string `json:"timeStamp"`
			Sign      string `json:"sign"`
		}
		ar := &appResult{
			AppID:     app.WeChatRetResult.AppID,
			PartnerId: app.WeChatRetResult.PartnerId,
			Package:   app.WeChatRetResult.Package,
			PrepayId:  app.WeChatRetResult.PrepayId,
			NonceStr:  app.WeChatRetResult.NonceStr,
			TimeStamp: app.WeChatRetResult.TimeStamp,
			Sign:      app.WeChatRetResult.Sign,
		}
		res.SetOk(ar)
		w.Write(res.ToJsonBytes())
	})
	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {

		callback(w, r)

	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("hello"))

	})
	http.ListenAndServe(":888", nil)
}

func TestPay(t *testing.T) {

	start()
}

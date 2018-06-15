// Create by Yale 2018/6/13 10:32
package paysdk

type OrderParams struct {
	BusTradeNo string  `json:"busTradeNo"`
	Price      float64 `json:"price"`
	Title      string  `json:"title"`
	Detail     string  `json:"detail"`

	Openid  string `json:"openid"`
	PayType int    `json:"payType"`
	AppName string `json:"appName"`
}

type PaySDK interface {
	Order(params *OrderParams) error
}

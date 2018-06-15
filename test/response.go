// Create by Yale 2018/6/15 15:07
package test

import "encoding/json"

type Response struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Results interface{} `json:"results"`
}

func NewResponse() *Response {
	return &Response{Success: false}
}

func (r *Response) SetOk(obj interface{}) *Response {
	r.Success = true
	r.Results = obj
	return r
}

func (r *Response) ToJsonBytes() []byte {
	btyes, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return btyes
}

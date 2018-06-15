// Create by Yale 2018/6/13 17:33
package common

import (
	"net/http"
	"strings"
	"io/ioutil"
	"crypto/tls"
	"time"
)
func getHttpsClient() *http.Client {

	config := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{TLSClientConfig: config}
	return &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
}
func getHttpClient() *http.Client {
	return &http.Client{Timeout:   30 * time.Second,}
}
func HttpPost(url string, contentType string, data string)([]byte, error)   {
	client := getHttpClient()
	resp, err := client.Post(url,contentType,strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
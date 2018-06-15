// Create by Yale 2018/6/13 15:45
package util

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

func GetRandomStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func GetPriceFen(fee float64) int64 {
	return int64(fee * 100)
}
func GetPriceFenStr(fee float64) string {
	return strconv.FormatInt(int64(fee*100), 10)
}
func GetPriceYunStr(fee int64) string {
	f := float64(fee) / float64(100)
	return strconv.FormatFloat(f, 'f', -1, 64)
}
func GetPriceYun(fee int64) float64 {
	f := float64(fee) / float64(100)
	return f
}

func LocalIP() string {
	info, _ := net.InterfaceAddrs()
	for _, addr := range info {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return ""
}

func Map2XML(m map[string]string) string {
	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	xml := fmt.Sprintf("<xml>%s</xml>", buf.String())
	return xml

}

func XmlToMap(xmlData []byte) (map[string]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	m := make(map[string]string)
	var token xml.Token
	var err error
	var k string
	for token, err = decoder.Token(); err == nil; token, err = decoder.Token() {
		if v, ok := token.(xml.StartElement); ok {
			k = v.Name.Local
			continue
		}
		if v, ok := token.(xml.CharData); ok {
			data := string(v.Copy())
			if strings.TrimSpace(data) == "" {
				continue
			}
			m[k] = data
		}
	}

	if err != nil && err != io.EOF {
		return nil, err
	}
	return m, nil
}

// Create by Yale 2018/6/13 17:06
package util

import (
	"fmt"
	"testing"
)

func TestLocalIP(t *testing.T) {

	fmt.Println(LocalIP())
}

func TestGetFee(t *testing.T) {

	f := GetPriceFen(0.01)
	fmt.Println(f)
	f = GetPriceFen(100.01)
	fmt.Println(f)

	f = GetPriceFen(100.87)
	fmt.Println(f)
	f = GetPriceFen(100.10)
	fmt.Println(f)
	f = GetPriceFen(100.102)
	fmt.Println(f)
}

func TestGetPriceYun(t *testing.T) {
	f := GetPriceYunStr(1)
	fmt.Println(f)

	f = GetPriceYunStr(10001)
	fmt.Println(f)

	f = GetPriceYunStr(10087)
	fmt.Println(f)

	f = GetPriceYunStr(10010)
	fmt.Println(f)

	f = GetPriceYunStr(100102)
	fmt.Println(f)

}

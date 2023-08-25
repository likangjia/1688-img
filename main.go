package main

import (
	"1688/colly"
	"flag"
)

var goodsUrl = flag.String("g", "aa", "need -g url")
var storeDir = flag.String("d", "ccc", "need -g url")

func main() {
	flag.Parse()
	//goodsUrl := "https://detail.1688.com/offer/692549361042.html"
	//// 遍历商品页面
	//err := CollPage(goodsUrl)
	//fmt.Println("err is ", err)
	c := colly.GetColly(colly.DetailColl, colly.Options{
		Url:      *goodsUrl,
		StoreDir: "./tmp/" + *storeDir,
	})
	c.CollyPage()
}

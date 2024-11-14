package colly

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/likangjia/1688-img/file"
	header2 "github.com/likangjia/1688-img/header"
	"github.com/spf13/cast"
	"os"
	"regexp"
	"strings"
	"time"
)

type DetailPage struct {
	CollyBase
	DownImages map[interface{}]string
}

var i = 0
var title string

func (d *DetailPage) CollyPage() {
	//selector := `#\31 081181308831 > div > div > div.layout-left > div.detail-gallery-wrapper > div > div.detail-gallery-turn > div > div`
	d.DownImages = make(map[interface{}]string)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
	)
	headers := header2.GetRequestHeaderWithTxt()
	// 设置请求头部信息
	c.OnRequest(func(r *colly.Request) {
		for k, h := range headers {
			r.Headers.Set(k, h)
		}
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})
	c.OnResponse(func(r *colly.Response) {
		// 获取详情数据
		body := string(r.Body)
		var re *regexp.Regexp
		if r.Request.URL.String() == d.Url {
			d.matchMainAndProp(body)
		}
		// 匹配detailUrl后面的链接
		re = regexp.MustCompile(`"detailUrl":"(.*?)"`)
		match := re.FindStringSubmatch(body)
		if len(match) > 1 {
			c.Visit(match[1])
		} else {
			d.matchDetail(body)
		}
	})
	//<img class=\"dynamic-backup-img\" style=\"display: block;width: 100.0%;height: auto;\" title=\"预览状态下无法点击，发布后，可点击跳转到对应的商品页面\" src=\"https://cbu01.alicdn.com/img/ibank/O1CN01IfKe2U1Bs2qr6tBEv_!!0-0-cib.jpg\" usemap=\"#_sdmap_12\"/>

	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("header", e.Text)
		title = e.Text

	})

	//#\31 081181309884 > div > div.detail-affix-sku-wrapper.hide > div.pc-sku-wrapper > div.sku-module-wrapper.sku-prop-module > div.prop-item-wrapper

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	err := c.Visit(d.Url)
	if err != nil {
		return
	}
	if d.StoreDir != "" {
		err = os.MkdirAll(d.genTileDir(), os.ModePerm)
		if err != nil {
			return
		}
		for k, s := range d.DownImages {
			split := strings.Split(s, ".")
			ex := strings.Split(split[len(split)-1], "?")

			err = file.Download(s, d.genTileDir(), fmt.Sprintf("%s.%s", k, ex[0]))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	return
}

func formatSortName(name interface{}) string {
	defer func() {
		i++
	}()
	// 过滤敏感词

	name = file.HandleFileName(cast.ToString(name))
	return fmt.Sprintf("%d--%v", i, name)
}

func (d *DetailPage) matchDetail(body string) {
	re := regexp.MustCompile(`<img(.*?)src=\\"(.*?)\\"(.*?)/>`)
	matches := re.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		trimStrUrl := strings.Replace(match[2], "\\", "", -1)
		d.DownImages[formatSortName("detail")] = trimStrUrl
	}
}

func (d *DetailPage) matchMainAndProp(body string) {

	reMainList := regexp.MustCompile(`"offerImgList":\[(.*?)]`)
	matches := reMainList.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		trimStr := strings.Replace(match[1], "\\", "", -1)
		trimStr = fmt.Sprintf("[%s]", trimStr)
		var mJsonMain = make([]string, 0)
		err := json.Unmarshal([]byte(trimStr), &mJsonMain)
		if err != nil {
			fmt.Println("err", err)
			return
		}
		for _, m := range mJsonMain {
			d.DownImages[formatSortName(i)] = m
		}
	}

	re := regexp.MustCompile(`value":\[(.*?)]`)
	matches = re.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		var mJsonMainSku = make([]map[string]string, 0)
		err := json.Unmarshal([]byte(fmt.Sprintf("[%s]", match[1])), &mJsonMainSku)
		if err != nil {
			fmt.Println("err", err)
			return
		}
		for _, m := range mJsonMainSku {
			if _, ok := m["imageUrl"]; !ok {
				continue
			}
			d.DownImages[formatSortName(m["name"])] = m["imageUrl"]
		}
	}
}

func (d *DetailPage) genTileDir() string {
	return fmt.Sprintf("%s/%s/%s", d.StoreDir, time.Now().Format("2006-01-02"), strings.Replace(file.HandleFileName(title), "/", "or", -1))
}

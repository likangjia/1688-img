package colly

import (
	"1688/file"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
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
	// 设置请求头部信息
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Pragma", "no-cache")
		r.Headers.Set("Cache-Control", "no-cache")
		r.Headers.Set("sec-ch-ua", `"Not/A)Brand";v="99", "Google Chrome";v="115", "Chromium";v="115"`)
		r.Headers.Set("sec-ch-ua-mobile", "?0")
		r.Headers.Set("sec-ch-ua-platform", `"macOS"`)
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36`)
		r.Headers.Set("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`)
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
		r.Headers.Set("Cookie", `cna=jrWxGv0u2HMCAXkhke9fObLc; _bl_uid=LjlUIf389wp6IC9X6f0sp3hw5C63; taklid=94ff91e5576e4f8a88b7f92c1f322fee; ali_apache_id=11.186.201.1.1688632544400.352396.1; xlly_s=1; _csrf_token=1691640098552; __cn_logon__=false; cookie2=18f6829662bed8434a5567b63195c88c; t=a5e2eef6946124d32f1a406e3f7b72c5; _tb_token_=e735e8b8ee651; _m_h5_tk=33b5576d3149bd256cd150dc3cc90181_1691655169471; _m_h5_tk_enc=80848228d031e4538fc508f9d2766617; JSESSIONID=97569E04F81E036D8B73E923A2CC4A7D; tfstk=dlq22Fj35cV5vdzxU2iaUweZbAnxrDCIilGsIR2ihjcc6tHibSNoSxDMjL4zs73jim697hyZwEDiIxDoExHBiFDDozJap7o__RBxIfVg_fafPMwYHcnGO5SCA-eqdcfI1zdg6-ntj69WjwsUHWLh57Wyfl1GsF0cxRhu0wEtga-Srb-MYpHqoXsKahxa4xyrtOlp6YDNllU2sFumeYlCUTkx_1Ud.; l=fBxA9mxgNLyNmh_DKOfZFurza779sIRAguPzaNbMi9fPO8WD5orRW19Rbf-kCnGVF6-vR3S7jpCMBeYBcCVBSYbht4pwuEkmnmOk-Wf..; isg=BO3tp224dZhlqREx2KM83WkX_I9nSiEcc7A-mC_yKwTzpg1Y95sP7BB8kHpAJjnU`)
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

	re := regexp.MustCompile(`skuProps\\":(.*?),\\"value\\":(.*?),\\"prop`)
	matches = re.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		trimStrSku := strings.Replace(match[2], "\\", "", -1)
		var mJsonMainSku = make([]map[string]string, 0)
		err := json.Unmarshal([]byte(trimStrSku), &mJsonMainSku)
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
	return fmt.Sprintf("%s/%s/%s", d.StoreDir, time.Now().Format("2006-01-02"), strings.Replace(title, "/", "or", -1))
}

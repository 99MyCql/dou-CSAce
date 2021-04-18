// 从 dblp 中爬取论文的引用数据
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	paperModel "douCSAce/app/paper/model"
	"douCSAce/pkg"
)

const (
	citationUrl  = "https://opencitations.net/index/api/v1/citations/%s?format=json&exclude=citing&sort=desc(creation)&mailto=ajax@dblp.org"
	doiSearchUrl = "https://dblp.org/search/publ/doi"
	logFileName  = ""
	// 由于 ArangoDB 文档 Update 之后会放到集合结尾，所以当处理完最后一个文档（最初的最后一个文档），就算处理完毕
	endKey = "journals-tocs-BugnionDRSW12"
)

func getPaper(dblpKey string) *paperModel.Paper {
	var p *paperModel.Paper
	key := paperModel.GenKey(dblpKey)
	if exist, _ := paperModel.IsExist(key); exist == true {
		p, _ = paperModel.FindByKey(key)
	}
	return p
}

// 爬取论文引用数据
func crawlPaperCit(p *paperModel.Paper) {
	pkg.Log.Info(p.ID)

	// 删除论文关联的 CitBy 边
	if err := p.DeleteCitBy(); err != nil {
		pkg.Log.Fatal(err)
	}

	// 请求获取论文引用数据的接口
	if strings.Index(p.DoiUrl, "https://doi.org/") == -1 {
		pkg.Log.Warn("no doi url")
		return
	}
	url := fmt.Sprintf(citationUrl, strings.ReplaceAll(p.DoiUrl, "https://doi.org/", ""))
	pkg.Log.Info(url)
	rsp, err := http.Get(url)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	var citJson []map[string]interface{}
	if err := json.Unmarshal(body, &citJson); err != nil {
		pkg.Log.Fatal(err)
	}

	// 根据返回的引用论文的doi，请求 dblp 接口获取对应的论文
	citDois := "doi="
	for i := 0; i < len(citJson); i++ {
		citDois += strings.ReplaceAll(citJson[i]["citing"].(string), "coci => ", "") + "+"
	}
	pkg.Log.Info(citDois)
	rsp, err = http.Post(doiSearchUrl,
		"application/x-www-form-urlencoded; charset=UTF-8",
		strings.NewReader(citDois))
	if err != nil {
		pkg.Log.Fatal(err)
	}

	// 解析返回的 html ，获取论文ID
	html, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	citCount := 0
	html.Find("#main>ul>li").Each(func(i int, s *goquery.Selection) {
		dblpKey, _ := s.Attr("id")
		pkg.Log.Infof("%d: %s", i, dblpKey)
		citingP := getPaper(dblpKey)
		if citingP != nil {
			cb := &paperModel.CitBy{
				From: p.ID,
				To:   citingP.ID,
			}
			if err := cb.Create(); err != nil {
				pkg.Log.Fatal(err)
			}
		}
		citCount++
	})
	if err := p.Update(map[string]interface{}{
		"CitationCount": citCount,
	}); err != nil {
		pkg.Log.Fatal(err)
	}
}

// 初始化：读取配置、启动日志、连接数据库
func init() {
	pkg.Conf = pkg.NewConfig("../conf.yaml")
	pkg.Log = pkg.NewLog(logFileName, pkg.DebugLevel)
	pkg.DB = pkg.NewDB(pkg.Conf.ArangoDB.Url, pkg.Conf.ArangoDB.Username,
		pkg.Conf.ArangoDB.Passwd, pkg.Conf.ArangoDB.Database, pkg.Conf.ArangoDB.ModelColNameMap)
}

func main() {
	for i := int64(0); ; i++ {
		count := 1000
		papers, err := paperModel.List(i, count)
		if err != nil {
			pkg.Log.Fatal(err)
		}
		for j := 0; j < count; j++ {
			pkg.Log.Infof("====== %d ======", i+int64(j))
			crawlPaperCit(papers[j])
			if papers[j].Key == endKey {
				pkg.Log.Info("crawler finished. congratulation!")
				return
			}
		}
		i += int64(count)
	}
}

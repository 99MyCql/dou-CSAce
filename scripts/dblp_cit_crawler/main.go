// 从 dblp 中爬取论文的引用数据
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/beevik/etree"

	paperModel "douCSAce/app/paper/model"
	"douCSAce/pkg"
)

const (
	logFilepath  = "log.log"
	confFilepath = "../../conf.yaml"
	// 由于 ArangoDB 文档 Update 之后会放到集合结尾，所以当处理完最后一个文档（最初的最后一个文档），就算处理完毕
	endKey    = "journals-ese-Briand97"
	goroutine = 1
	logLevel  = pkg.DebugLevel
)

var (
	start = 0
)

// getPaper 根据 dblpKey 从数据库中读取对应数据
// func getPaper(dblpKey string) *paperModel.Paper {
// 	var p *paperModel.Paper
// 	key := paperModel.GenKey(dblpKey)
// 	if exist, _ := paperModel.IsExist(key); exist == true {
// 		p, _ = paperModel.FindByKey(key)
// 	}
// 	return p
// }

// requestDBLPDoi 根据 doi 接口返回的引用论文的doi，请求 dblp 接口获取对应的论文
// func requestDBLPDoi(citDoiStr string, p *paperModel.Paper) {
// 	if citDoiStr == "" {
// 		pkg.Log.Warn("no cit")
// 		return
// 	}
//
// 	// 请求 dblp 接口
// 	rsp, err := http.Post("https://dblp.org/search/publ/doi",
// 		"application/x-www-form-urlencoded; charset=UTF-8",
// 		strings.NewReader("doi="+citDoiStr))
// 	if err != nil {
// 		pkg.Log.Fatal(err)
// 	}
// 	if rsp.StatusCode != http.StatusOK {
// 		pkg.Log.Fatal(rsp.Status)
// 	}
//
// 	// 解析返回的 html ，获取论文ID
// 	html, err := goquery.NewDocumentFromReader(rsp.Body)
// 	if err != nil {
// 		pkg.Log.Fatal(err)
// 	}
// 	html.Find("#main>ul>li").Each(func(i int, s *goquery.Selection) {
// 		dblpKey, _ := s.Attr("id")
// 		pkg.Log.Infof("cit %d: %s", i, dblpKey)
// 		citingP := getPaper(dblpKey)
// 		if citingP != nil {
// 			cb := &paperModel.CitBy{
// 				From: p.ID,
// 				To:   citingP.ID,
// 			}
// 			if err := cb.Create(); err != nil {
// 				pkg.Log.Fatal(err)
// 			}
// 		}
// 	})
// }

// requestOpen 请求 https://opencitations.net/ 接口获取论文引用数据
func requestOpen(doi string) ([]string, error) {
	url := fmt.Sprintf("https://opencitations.net/index/api/v1/citations/%s?format=json&exclude=citing&sort=desc(creation)&mailto=ajax@dblp.org", doi)
	pkg.Log.Info(url)
	rsp, err := http.Get(url)
	if err != nil {
		pkg.Log.Error(err)
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		pkg.Log.Error(rsp.Status + ":" + url)
		return nil, errors.New(rsp.Status)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
	}

	var dataJson []map[string]interface{}
	if err := json.Unmarshal(body, &dataJson); err != nil {
		pkg.Log.Fatal(err)
	}
	var citDois []string
	for i := 0; i < len(dataJson); i++ {
		citDois = append(citDois, strings.ReplaceAll(dataJson[i]["citing"].(string), "coci => ", ""))
	}
	return citDois, nil
}

// requestSem 请求 https://api.semanticscholar.org 接口获取论文引用数据
func requestSem(doi string) ([]string, error) {
	url := fmt.Sprintf(
		"https://api.semanticscholar.org/v1/paper/%s?include_unknown_references=true&mailto=ajax@dblp.org", doi)
	pkg.Log.Info(url)
	rsp, err := http.Get(url)
	if err != nil {
		pkg.Log.Error(err)
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		pkg.Log.Error(rsp.Status + ":" + url)
		return nil, errors.New(rsp.Status)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
	}

	var dataJson map[string]interface{}
	if err := json.Unmarshal(body, &dataJson); err != nil {
		pkg.Log.Fatal(err)
	}
	citJson := dataJson["citations"].([]interface{})
	var citDois []string
	for i := 0; i < len(citJson); i++ {
		if citJson[i].(map[string]interface{})["doi"] == nil {
			citDois = append(citDois, "")
		} else {
			citDois = append(citDois, citJson[i].(map[string]interface{})["doi"].(string))
		}
	}
	return citDois, nil
}

// 爬取论文引用数据
func crawlPaperCit(p *paperModel.Paper) {
	pkg.Log.Info(p.Key)

	// 删除论文关联的 CitBy 边
	if err := p.DeleteCitBy(); err != nil {
		pkg.Log.Fatal(err)
	}

	// 根据论文的 DoiUrl 请求对应接口，获取论文引用数据
	var citDois []string
	isMatch := false
	if p.DoiUrl == "" {
		isMatch = true
		pkg.Log.Warn("doi url is empty!")
	}
	for i := range handledUrls {
		if strings.HasPrefix(p.DoiUrl, handledUrls[i][0]) {
			var err error
			if handledUrls[i][2] == "requestOpen" {
				citDois, err = requestOpen(strings.ReplaceAll(p.DoiUrl, handledUrls[i][0], handledUrls[i][1]))
				if err != nil {
					citDois, err = requestSem(strings.ReplaceAll(p.DoiUrl, handledUrls[i][0], handledUrls[i][1]))
				}
			} else if handledUrls[i][2] == "requestSem" {
				citDois, err = requestSem(strings.ReplaceAll(p.DoiUrl, handledUrls[i][0], handledUrls[i][1]))
				if err != nil {
					citDois, err = requestOpen(strings.ReplaceAll(p.DoiUrl, handledUrls[i][0], handledUrls[i][1]))
				}
			}
			// 如果爬取遇到错误，则 start++ ，并不进行论文引用数更新
			if err != nil {
				start++
				return
			}
			isMatch = true
			break
		}
	}
	for i := 0; !isMatch && i < len(unhandledUrls); i++ {
		if strings.Contains(p.DoiUrl, unhandledUrls[i]) {
			pkg.Log.Warn("doi url can't be handled: " + p.DoiUrl)
			isMatch = true
			break
		}
	}
	// 如果遇到未匹配的 doiUrl 也 start++ ，不进行论文引用数更新
	if !isMatch {
		pkg.Log.Error(p.Key + "'s doi url don't match: " + p.DoiUrl)
		start++
		return
	}
	pkg.Log.Infof("%s: %d-%+v", p.Key, len(citDois), citDois)

	// 请求 dblp 接口，获取引用数据对应的论文 dblpKey
	// requestDBLPDoi(citDoiStr, p)

	// 更新论文的引用数
	if err := p.Update(map[string]interface{}{"citationCount": uint64(len(citDois))}); err != nil {
		pkg.Log.Fatal(err)
	}
}

// changePaperDoiUrl 修改对应论文的 doiUrl
func changePaperDoiUrl(p *paperModel.Paper) {
	pkg.Log.Infof("paper(%s)'s old doiUrl: %s", p.Key, p.DoiUrl)
	// 爬取 dblp 论文数据
	url := "https://dblp.org/rec/" + p.GetDblpKey() + ".xml"
	pkg.Log.Info(url)
	rsp, err := http.Get(url)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	if rsp.StatusCode == http.StatusNotFound {
		pkg.Log.Error(rsp.Status + ":" + url)
		return
	} else if rsp.StatusCode != http.StatusOK {
		pkg.Log.Fatal(rsp.Status + ":" + url)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
	}

	// 解析 xml 数据
	paperXml := etree.NewDocument()
	if err := paperXml.ReadFromBytes(body); err != nil {
		pkg.Log.Fatal(err)
	}
	urls := paperXml.FindElements("/dblp/inproceedings/ee")
	for _, url := range urls {
		pkg.Log.Debug(url.Text())
		if strings.HasPrefix(url.Text(), "https://doi.org/") {
			p.DoiUrl = url.Text()
			// 更新论文的引用数
			if err := p.Update(map[string]interface{}{"doiUrl": p.DoiUrl}); err != nil {
				pkg.Log.Fatal(err)
			}
		}
	}
	pkg.Log.Infof("paper(%s)'s new doiUrl: %s", p.Key, p.DoiUrl)
}

// 初始化：读取配置、启动日志、连接数据库
func init() {
	pkg.Conf = pkg.NewConfig(confFilepath)
	pkg.Log = pkg.NewLog(logFilepath, logLevel)
	pkg.DB = pkg.NewDB(pkg.Conf.ArangoDB.Url, pkg.Conf.ArangoDB.Username,
		pkg.Conf.ArangoDB.Passwd, pkg.Conf.ArangoDB.Database, pkg.Conf.ArangoDB.ModelColNameMap)
}

func main() {
	for i := uint64(0); ; {
		pkg.Log.Infof("****** %d - start:%d ******", i, start)
		// 由于 ArangoDB 文档 Update 之后会放到集合结尾，所以每次都从 start 开始（遇到错误不进行更新 start++）
		papers, err := paperModel.List(uint64(start), uint(goroutine))
		if err != nil {
			pkg.Log.Fatal(err)
		}

		var end = false
		var wg sync.WaitGroup
		wg.Add(len(papers))
		for j := 0; j < len(papers); j++ {
			go func(index int) {
				defer wg.Done()
				pkg.Log.Infof("====== %d ======", i+uint64(index))
				// 更新这些 paper 对应的 doiUrl （当初爬取 paper 时未能获取 paper xml 中全部 ee 属性所遗留的 bug ）
				if strings.HasPrefix(papers[index].DoiUrl, "http://openaccess.thecvf.com/") ||
					strings.HasPrefix(papers[index].DoiUrl, "https://www.aclweb.org/") ||
					strings.HasPrefix(papers[index].DoiUrl, "http://dl.ifip.org/") {
					changePaperDoiUrl(papers[index])
				}
				crawlPaperCit(papers[index])
				if papers[index].Key == endKey {
					end = true
				}
			}(j)
		}
		wg.Wait()
		if end {
			pkg.Log.Info("Crawler finished. Congratulation!")
			break
		}
		i += uint64(len(papers))
	}
	// p, _ := paperModel.FindByKey("conf-acl-McCallumVVMR18")
	// changePaperDoiUrl(p)
	// crawlPaperCit(p)
}

// 从 dblp 中爬取论文的引用数据
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	paperModel "douCSAce/app/paper/model"
	"douCSAce/pkg"
)

const (
	logFilepath  = ""
	confFilepath = "../../conf.yaml"
	// 由于 ArangoDB 文档 Update 之后会放到集合结尾，所以当处理完最后一个文档（最初的最后一个文档），就算处理完毕
	endKey    = "journals-ese-Briand97"
	goroutine = 10
	logLevel  = pkg.DebugLevel
)

var (
	start = 0
	// 可爬取的URL
	handledUrls = [][2]string{
		{"https://doi.org/", "requestOpen"},
		{"http://doi.acm.org/", "requestSem"},
		{"http://doi.ieeecomputersociety.org/", "requestOpen"},
		{"http://link.springer.com/article/", "requestSem"},
		{"https://onlinelibrary.wiley.com/doi/abs/", "requestSem"},
		{"http://link.springer.com/chapter/", "requestSem"},
		{"https://www.tandfonline.com/doi/full/", "requestOpen"},
	}
	// 不可爬取的URL
	unhandledUrls = []string{
		"springerlink.com/",
		"usenix.org/",
		"portal.acm.org/",
		"ieeexplore.ieee.org/",
		"ceur-ws.org/",
		"sigda.org/",
		"net.doit.wisc.edu/",
		"storageconference.org/",
		"kluweronline.com/",
		"iptps.org/",
		"journals.cambridge.org/",
		"ejournals.wspc.com.sg/",
		"journals.riverpublishers.com/",
		"www.rintonpress.com/",
		"static.cambridge.org",
		"www.cs.bham.ac.uk/",
		"dl.acm.org/",
		"www.se-hci.org/",
		"se.uwaterloo.ca/",
		"www.acm.org/",
		"netfiles.uiuc.edu/",
		"www.disi.unige.it/",
		"www.info.uni-karlsruhe.de/",
		"easychair.org/",
		"portalparts.acm.org/",
		"www.computer.org/",
		"ewic.bcs.org/",
		"scg.unibe.ch/",
		"ftp.cs.man.ac.uk/",
		"www.vldb.org/",
		"www.worldscinet.com/",
		"www.informaworld.com/",
		"gemo.futurs.inria.fr/",
		"proceedings.mlr.press/",
		"drive.google.com/",
		"www.cs.rpi.edu/",
		"www.aaai.org/",
		"www.sdsc.edu/",
		"www.cs.cmu.edu/",
		"p2pir.is.informatik.uni-duisburg.de/",
		"link.springer.com/content/pdf/",
		"www.semanticweb.org/",
		"openproceedings.org/",
		"wwwiti.cs.uni-magdeburg.de/",
		"www.edbt2000.uni-konstanz.de/",
		"epubs.siam.org/",
		"www.booksonline.iospress.nl/",
		"www.cs.ust.hk/",
		"crpit.scem.westernsydney.edu.au/",
		"www-rocq.inria.fr/",
		"www.research.att.com/",
		"www.db.ucsd.edu/",
		"www.cse.ogi.edu/",
		"webdb2008.como.polimi.it/",
		"webdb09.cse.buffalo.edu/",
		"webdb2011.rutgers.edu/",
		"db.disi.unitn.eu/",
		"cidrdb.org/",
	}
)

// getPaper 根据 dblpKey 从数据库中读取对应数据
func getPaper(dblpKey string) *paperModel.Paper {
	var p *paperModel.Paper
	key := paperModel.GenKey(dblpKey)
	if exist, _ := paperModel.IsExist(key); exist == true {
		p, _ = paperModel.FindByKey(key)
	}
	return p
}

// requestOpen 请求 https://opencitations.net/ 接口获取论文引用数据
func requestOpen(doi string) []string {
	url := fmt.Sprintf("https://opencitations.net/index/api/v1/citations/%s?format=json&exclude=citing&sort=desc(creation)&mailto=ajax@dblp.org", doi)
	pkg.Log.Info(url)
	rsp, err := http.Get(url)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	if rsp.StatusCode != http.StatusOK {
		pkg.Log.Error(rsp.Status + ":" + url)
		start++
		return nil
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
	return citDois
}

// requestSem 请求 https://api.semanticscholar.org 接口获取论文引用数据
func requestSem(doi string) []string {
	url := fmt.Sprintf(
		"https://api.semanticscholar.org/v1/paper/%s?include_unknown_references=true&mailto=ajax@dblp.org", doi)
	pkg.Log.Info(url)
	rsp, err := http.Get(url)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	if rsp.StatusCode == http.StatusNotFound {
		pkg.Log.Error(rsp.Status + ":" + url)
		start++
		return nil
	} else if rsp.StatusCode != http.StatusOK {
		pkg.Log.Fatal(rsp.Status)
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
	return citDois
}

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
			if handledUrls[i][1] == "requestOpen" {
				citDois = requestOpen(strings.ReplaceAll(p.DoiUrl, handledUrls[i][0], ""))
			} else if handledUrls[i][1] == "requestSem" {
				citDois = requestSem(strings.ReplaceAll(p.DoiUrl, handledUrls[i][0], ""))
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

// 初始化：读取配置、启动日志、连接数据库
func init() {
	pkg.Conf = pkg.NewConfig(confFilepath)
	pkg.Log = pkg.NewLog(logFilepath, logLevel)
	pkg.DB = pkg.NewDB(pkg.Conf.ArangoDB.Url, pkg.Conf.ArangoDB.Username,
		pkg.Conf.ArangoDB.Passwd, pkg.Conf.ArangoDB.Database, pkg.Conf.ArangoDB.ModelColNameMap)
}

func main() {
	for i := uint64(0); ; {
		pkg.Log.Infof("****** %d ******", i)
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
				pkg.Log.Infof("====== %d ======", i+uint64(index))
				crawlPaperCit(papers[index])
				if papers[index].Key == endKey {
					end = true
				}
				wg.Done()
			}(j)
		}
		wg.Wait()
		if end {
			pkg.Log.Info("Crawler finished. Congratulation!")
			break
		}
		i += uint64(len(papers))
	}
	// p, _ := paperModel.FindByKey("journals-ton-MaloneDL07")
	// crawlPaperCit(p)
}

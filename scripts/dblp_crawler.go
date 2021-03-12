package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/beevik/etree"

	"douCSAce/app/field"
	"douCSAce/app/journal"
	"douCSAce/pkg"
)

const (
	jsonFile      = "ccf_field.json"
	CCFFieldType  = 1
	ConfPaperType = 1
	JouPaperType  = 2
)

func handlePaper(paperXml *etree.Element) {
}

func handleJournal(journalMap map[string]interface{}, f *field.Field, category string) {
	j := &journal.Journal{
		Name:          journalMap["name"].(string),
		ShortName:     journalMap["shortName"].(string),
		Publisher:     journalMap["publisher"].(string),
		DblpUrl:       journalMap["url"].(string),
		PaperCount:    0,
		CitationCount: 0,
	}
	// if err := j.Create(); err != nil {
	// 	pkg.Log.Fatal(err)
	// }
	// j2f := &journal.JouBelongToField{
	// 	From: j.ID,
	// 	To:   f.ID,
	// 	Note: category,
	// }
	// if err := j2f.Create(); err != nil {
	// 	pkg.Log.Fatal(err)
	// }

	// 请求接口获取数据
	journalUrl := j.DblpUrl + "index.xml"
	pkg.Log.Info(journalUrl)
	rsp, err := http.Get(journalUrl)
	if err != nil {
		pkg.Log.Fatal(err)
		return
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
		return
	}

	// 解析 xml 数据
	jouXml := etree.NewDocument()
	if err := jouXml.ReadFromBytes(body); err != nil {
		pkg.Log.Fatal(err)
		return
	}
	vols := jouXml.FindElements("/bht/ul/li/ref")
	// 爬取期刊种每个卷的数据
	for i := 0; i < len(vols); i++ {
		// 请求接口获取数据
		volUrl := "https://dblp.org/" +
			strings.ReplaceAll(vols[i].SelectAttr("href").Value, ".html", ".xml")
		pkg.Log.Info(volUrl)
		rsp, err := http.Get(volUrl)
		if err != nil {
			pkg.Log.Fatal(err)
		}
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			pkg.Log.Fatal(err)
		}

		// 解析 xml 数据
		volXml := etree.NewDocument()
		if err := volXml.ReadFromBytes(body); err != nil {
			pkg.Log.Fatal(err)
		}
		papers := volXml.FindElements("/bht/dblpcites/r/article")
		for j := 0; j < len(papers); j++ {
			handlePaper(papers[j])
		}
	}
}

func handleJournals(journalsMap map[string]interface{}, f *field.Field) {
	journalsAMap := journalsMap["A"].([]interface{})
	for i := 0; i < len(journalsAMap); i++ {
		journalMap := journalsAMap[i].(map[string]interface{})
		handleJournal(journalMap, f, "A")
	}

	journalBMap := journalsMap["B"].([]interface{})
	for i := 0; i < len(journalBMap); i++ {
		journalMap := journalBMap[i].(map[string]interface{})
		handleJournal(journalMap, f, "B")
	}

	journalCMap := journalsMap["C"].([]interface{})
	for i := 0; i < len(journalCMap); i++ {
		journalMap := journalCMap[i].(map[string]interface{})
		handleJournal(journalMap, f, "C")
	}
}

func handleConferences(confMap map[string]interface{}, f *field.Field) {

}

func handleField(fieldMap map[string]interface{}, typ uint) {
	pkg.Log.Info(fieldMap["zhName"].(string))
	f := &field.Field{
		Name:          fieldMap["name"].(string),
		ZhName:        fieldMap["zhName"].(string),
		Type:          typ,
		PaperCount:    0,
		CitationCount: 0,
	}
	// if err := f.Create(); err != nil {
	// 	pkg.Log.Fatal(err)
	// }
	handleJournals(fieldMap["journal"].(map[string]interface{}), f)
	handleConferences(fieldMap["conference"].(map[string]interface{}), f)
}

func init() {
	pkg.Conf = pkg.NewConfig("../conf.yaml")
	pkg.Log = pkg.NewLog(pkg.Conf.LogPath, pkg.DebugLevel)
	pkg.DB = pkg.NewDB(pkg.Conf.ArangoDB.Url, pkg.Conf.ArangoDB.Username,
		pkg.Conf.ArangoDB.Passwd, pkg.Conf.ArangoDB.Database, pkg.Conf.ArangoDB.Cols)
}

func main() {
	// 打开 json 文件
	f, err := os.Open(jsonFile)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	defer f.Close()

	// 解析 json 文件
	var jsonData []map[string]interface{}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&jsonData); err != nil {
		pkg.Log.Fatal(err)
	}

	// 处理各个研究方向
	for i := 0; i < len(jsonData); i++ {
		handleField(jsonData[i], CCFFieldType)
	}
}

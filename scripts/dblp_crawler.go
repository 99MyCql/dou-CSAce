// 从 dblp 爬取数据，处理，并存储到数据库
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/beevik/etree"

	authorModel "douCSAce/app/author/model"
	confInsModel "douCSAce/app/confInstance/model"
	confSerModel "douCSAce/app/confSeries/model"
	fieldModel "douCSAce/app/field/model"
	jouModel "douCSAce/app/journal/model"
	paperModel "douCSAce/app/paper/model"
	"douCSAce/pkg"
)

const (
	jsonFileName  = "ccf_field.json"
	CCFFieldType  = 1
	ConfPaperType = 1
	JouPaperType  = 2
)

func handleAuthor(authorXml *etree.Element) *authorModel.Author {
	key := authorModel.GenKey(authorXml.SelectAttr("pid").Value)
	pkg.Log.Info("handle author:" + key)

	var a *authorModel.Author
	// 如果已存在则直接获取
	if exist, _ := authorModel.IsExist(key); exist {
		a, _ = authorModel.FindByKey(key)
	} else {
		a = &authorModel.Author{
			Key:  key,
			Name: authorXml.Text(),
		}
		if err := a.Create(); err != nil {
			pkg.Log.Fatal(err)
		}
	}
	return a
}

func handlePaper(paperXml *etree.Element, typ uint, venue interface{}) {
	key := paperModel.GenKey(paperXml.SelectAttr("key").Value)
	pkg.Log.Info("handle paper:" + key)

	// 处理可能出现的空值
	var (
		doiUrl string
		pages  string
		number string
	)
	if ee := paperXml.SelectElement("ee"); ee != nil {
		doiUrl = ee.Text()
	}
	if pa := paperXml.SelectElement("pages"); pa != nil {
		pages = pa.Text()
	}
	if nu := paperXml.SelectElement("number"); nu != nil {
		number = nu.Text()
	}

	// 处理不同类型的 paper
	var p *paperModel.Paper
	if typ == ConfPaperType {
		// 如果已存在则直接获取（此处不删除，因为一篇 paper 可能发表在多个不同期刊会议）
		if exist, _ := paperModel.IsExist(key); exist == true {
			p, _ = paperModel.FindByKey(key)
		} else {
			p = &paperModel.Paper{
				Key:       key,
				Title:     paperXml.SelectElement("title").Text(),
				Type:      typ,
				Pages:     pages,
				Year:      paperXml.SelectElement("year").Text(),
				BookTitle: paperXml.SelectElement("booktitle").Text(),
				DoiUrl:    doiUrl,
				DblpUrl:   "https://dblp.org/" + paperXml.SelectElement("url").Text(),
			}
			if err := p.Create(); err != nil {
				pkg.Log.Fatal(err)
			}
		}
		poCI := &paperModel.PublishOnConfIns{
			From: p.ID,
			To:   venue.(*confInsModel.ConfInstance).ID,
		}
		if err := poCI.Create(); err != nil {
			pkg.Log.Fatal(err)
		}
	} else if typ == JouPaperType {
		if exist, _ := paperModel.IsExist(key); exist == true {
			p, _ = paperModel.FindByKey(key)
		} else {
			p = &paperModel.Paper{
				Key:     paperModel.GenKey(paperXml.SelectAttr("key").Value),
				Title:   paperXml.SelectElement("title").Text(),
				Type:    typ,
				Pages:   pages,
				Year:    paperXml.SelectElement("year").Text(),
				Volume:  paperXml.SelectElement("volume").Text(),
				Number:  number,
				DoiUrl:  doiUrl,
				DblpUrl: "https://dblp.org/" + paperXml.SelectElement("url").Text(),
			}
			if err := p.Create(); err != nil {
				pkg.Log.Fatal(err)
			}
		}
		poJ := &paperModel.PublishOnJou{
			From: p.ID,
			To:   venue.(*jouModel.Journal).ID,
		}
		if err := poJ.Create(); err != nil {
			pkg.Log.Fatal(err)
		}
	} else {
		pkg.Log.Fatal(typ)
		return
	}

	authors := paperXml.SelectElements("author")
	for i := 0; i < len(authors); i++ {
		a := handleAuthor(authors[i])
		wb := &paperModel.WriteBy{
			From: p.ID,
			To:   a.ID,
		}
		if err := wb.Create(); err != nil {
			pkg.Log.Fatal(err)
		}
	}
}

func handleVolume(volumeXml *etree.Element, jou *jouModel.Journal) {
	// 请求接口获取数据
	volUrl := "https://dblp.org/" +
		strings.ReplaceAll(volumeXml.SelectAttr("href").Value, ".html", ".xml")
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
	for i := 0; i < len(papers); i++ {
		handlePaper(papers[i], JouPaperType, jou)
	}
}

func handleJournal(journalMap map[string]interface{}, f *fieldModel.Field, category string) {
	// 如果当前 journal 已爬取过，则跳过
	if journalMap["isHandle"].(bool) {
		return
	}

	key := jouModel.GenKey(journalMap["shortName"].(string))
	pkg.Log.Info("handle journal:" + key)

	// 如果 journal 已存在，则删除点及其关联的边（因为可能是上次爬取未正常结束的脏数据）
	if exist, _ := jouModel.IsExist(key); exist == true {
		jouOld, _ := jouModel.FindByKey(key)
		if err := jouOld.DeleteJouBelongToField(); err != nil {
			pkg.Log.Fatal(err)
		}
		if err := jouOld.Delete(); err != nil {
			pkg.Log.Fatal(err)
		}
	}

	jou := &jouModel.Journal{
		Key:       jouModel.GenKey(journalMap["shortName"].(string)),
		Name:      journalMap["name"].(string),
		ShortName: journalMap["shortName"].(string),
		Publisher: journalMap["publisher"].(string),
		Url:       journalMap["url"].(string),
	}
	if err := jou.Create(); err != nil {
		pkg.Log.Fatal(err)
	}
	j2f := &jouModel.JouBelongToField{
		From: jou.ID,
		To:   f.ID,
		Note: category,
	}
	if err := j2f.Create(); err != nil {
		pkg.Log.Fatal(err)
	}

	// 非 dblp 接口则终止
	if strings.Count(jou.Url, "https://dblp.org") == 0 {
		pkg.Log.Warn("isn't dblp url:" + jou.Url)
		return
	}

	// 请求接口获取数据
	journalUrl := jou.Url + "index.xml"
	pkg.Log.Info(journalUrl)
	rsp, err := http.Get(journalUrl)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
	}

	// 解析 xml 数据
	jouXml := etree.NewDocument()
	if err := jouXml.ReadFromBytes(body); err != nil {
		pkg.Log.Fatal(err)
	}
	vols := jouXml.FindElements("/bht/ul/li/ref")
	// 爬取期刊种每个卷的数据
	for i := 0; i < len(vols); i++ {
		handleVolume(vols[i], jou)
	}
}

func handleJournals(journalsMap map[string]interface{}, f *fieldModel.Field) {
	journalsAMap := journalsMap["A"].([]interface{})
	for i := 0; i < len(journalsAMap); i++ {
		journalMap := journalsAMap[i].(map[string]interface{})
		handleJournal(journalMap, f, "A")
	}

	journalsBMap := journalsMap["B"].([]interface{})
	for i := 0; i < len(journalsBMap); i++ {
		journalMap := journalsBMap[i].(map[string]interface{})
		handleJournal(journalMap, f, "B")
	}

	journalsCMap := journalsMap["C"].([]interface{})
	for i := 0; i < len(journalsCMap); i++ {
		journalMap := journalsCMap[i].(map[string]interface{})
		handleJournal(journalMap, f, "C")
	}
}

func handleConfInstances(confInstanceXml *etree.Element, cs *confSerModel.ConfSeries) {
	key := confInsModel.GenKey(confInstanceXml.SelectAttr("key").Value)
	pkg.Log.Info("handle conference instance:" + key)

	// 如果 confInstance 已存在，则删除点及其关联的边（因为可能是上次爬取未正常结束的脏数据）
	if exist, _ := confInsModel.IsExist(key); exist == true {
		confInsOld, _ := confInsModel.FindByKey(key)
		if err := confInsOld.DeleteConfInsBelongToConfSer(); err != nil {
			pkg.Log.Fatal(err)
		}
		if err := confInsOld.Delete(); err != nil {
			pkg.Log.Fatal(err)
		}
	}

	// 处理可能为空的值
	var (
		isbn      string
		doiUrl    string
		bookTitle string
		publisher string
	)
	if i := confInstanceXml.SelectElement("isbn"); i != nil {
		isbn = i.Text()
	}
	if d := confInstanceXml.SelectElement("ee"); d != nil {
		doiUrl = d.Text()
	}
	if b := confInstanceXml.SelectElement("booktitle"); b != nil {
		bookTitle = b.Text()
	}
	if p := confInstanceXml.SelectElement("publisher"); p != nil {
		publisher = p.Text()
	}

	ci := &confInsModel.ConfInstance{
		Key:       key,
		Title:     confInstanceXml.SelectElement("title").Text(),
		Publisher: publisher,
		BookTitle: bookTitle,
		Year:      confInstanceXml.SelectElement("year").Text(),
		Isbn:      isbn,
		DoiUrl:    doiUrl,
		DblpUrl:   "https://dblp.org/" + confInstanceXml.SelectElement("url").Text(),
	}
	if err := ci.Create(); err != nil {
		pkg.Log.Fatal(err)
	}
	ci2cs := &confInsModel.ConfInsBelongToConfSer{
		From: ci.ID,
		To:   cs.ID,
	}
	if err := ci2cs.Create(); err != nil {
		pkg.Log.Fatal(err)
	}

	// 请求接口获取数据
	confInsUrl := strings.ReplaceAll(ci.DblpUrl, ".html", ".xml")
	pkg.Log.Info(confInsUrl)
	rsp, err := http.Get(confInsUrl)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
	}

	// 解析 xml 数据
	confInsXml := etree.NewDocument()
	if err := confInsXml.ReadFromBytes(body); err != nil {
		pkg.Log.Fatal(err)
	}
	inproceedings := confInsXml.FindElements("/bht/dblpcites/r/inproceedings")
	for j := 0; j < len(inproceedings); j++ {
		handlePaper(inproceedings[j], ConfPaperType, ci)
	}
}

func handleConfSeries(confSeriesMap map[string]interface{}, f *fieldModel.Field, category string) {
	// isHandle 为 true 表示已爬取过
	if confSeriesMap["isHandle"].(bool) {
		return
	}

	key := confSerModel.GenKey(confSeriesMap["shortName"].(string))
	pkg.Log.Info("handle conference series:" + key)

	// 如果 confSeries 已存在，则删除点及其关联的边（因为可能是上次爬取未正常结束的脏数据）
	if exist, _ := confSerModel.IsExist(key); exist == true {
		confSerOld, _ := confSerModel.FindByKey(key)
		if err := confSerOld.DeleteConfSerBelongToField(); err != nil {
			pkg.Log.Fatal(err)
		}
		if err := confSerOld.Delete(); err != nil {
			pkg.Log.Fatal(err)
		}
	}

	cs := &confSerModel.ConfSeries{
		Key:       key,
		Name:      confSeriesMap["name"].(string),
		ShortName: confSeriesMap["shortName"].(string),
		Publisher: confSeriesMap["publisher"].(string),
		Url:       confSeriesMap["url"].(string),
	}
	if err := cs.Create(); err != nil {
		pkg.Log.Fatal(err)
	}
	cs2f := &confSerModel.ConfSerBelongToField{
		From: cs.ID,
		To:   f.ID,
		Note: category,
	}
	if err := cs2f.Create(); err != nil {
		pkg.Log.Fatal(err)
	}

	// 非 dblp 接口则终止
	if strings.Count(cs.Url, "https://dblp.org") == 0 {
		pkg.Log.Warn("isn't dblp url:" + cs.Url)
		return
	}

	// 请求接口获取数据
	confUrl := cs.Url + "index.xml"
	pkg.Log.Info(confUrl)
	rsp, err := http.Get(confUrl)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
	}

	// 解析 xml 数据
	confXml := etree.NewDocument()
	if err := confXml.ReadFromBytes(body); err != nil {
		pkg.Log.Fatal(err)
	}
	proceedings := confXml.FindElements("/bht/dblpcites/r/proceedings")
	// 爬取期刊种每个会议实例的数据
	for i := 0; i < len(proceedings); i++ {
		handleConfInstances(proceedings[i], cs)
	}
}

func handleConferences(confMap map[string]interface{}, f *fieldModel.Field) {
	confsAMap := confMap["A"].([]interface{})
	for i := 0; i < len(confsAMap); i++ {
		confMap := confsAMap[i].(map[string]interface{})
		handleConfSeries(confMap, f, "A")
	}

	confsBMap := confMap["B"].([]interface{})
	for i := 0; i < len(confsBMap); i++ {
		confMap := confsBMap[i].(map[string]interface{})
		handleConfSeries(confMap, f, "B")
	}

	confsCMap := confMap["C"].([]interface{})
	for i := 0; i < len(confsCMap); i++ {
		confMap := confsCMap[i].(map[string]interface{})
		handleConfSeries(confMap, f, "C")
	}
}

func handleField(fieldMap map[string]interface{}, typ uint) {
	key := fieldModel.GenKey(typ, fieldMap["name"].(string))
	pkg.Log.Info("handle field:" + key)

	var f *fieldModel.Field
	// 若已存在则直接获取
	if exist, _ := fieldModel.IsExist(key); exist == true {
		f, _ = fieldModel.FindByKey(key)
	} else {
		f = &fieldModel.Field{
			Key:           fieldModel.GenKey(typ, fieldMap["name"].(string)),
			Name:          fieldMap["name"].(string),
			ZhName:        fieldMap["zhName"].(string),
			Type:          typ,
			PaperCount:    0,
			CitationCount: 0,
		}
		if err := f.Create(); err != nil {
			pkg.Log.Fatal(err)
		}
	}

	handleJournals(fieldMap["journal"].(map[string]interface{}), f)
	handleConferences(fieldMap["conference"].(map[string]interface{}), f)
}

// 初始化：读取配置、启动日志、连接数据库
func init() {
	pkg.Conf = pkg.NewConfig("../conf.yaml")
	pkg.Log = pkg.NewLog("", pkg.DebugLevel)
	pkg.DB = pkg.NewDB(pkg.Conf.ArangoDB.Url, pkg.Conf.ArangoDB.Username,
		pkg.Conf.ArangoDB.Passwd, pkg.Conf.ArangoDB.Database, pkg.Conf.ArangoDB.ColNames)
}

func main() {
	// 打开 json 文件
	f, err := os.Open(jsonFileName)
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

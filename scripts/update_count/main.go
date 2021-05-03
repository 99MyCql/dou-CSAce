// 用于更新部分模型的论文数和引用数
package main

import (
	authorModel "douCSAce/app/author/model"
	confSerModel "douCSAce/app/confSeries/model"
	fieldModel "douCSAce/app/field/model"
	jouModel "douCSAce/app/journal/model"
	"douCSAce/pkg"
)

const (
	logFilepath  = "log.log"
	confFilepath = "../../conf.yaml"
	logLevel     = pkg.DebugLevel

	updateAuthor = ``

	updateJou = `for j in journals
	let pcits = (
		for p in 1 inbound j._id publish_on_jou
			return p.citationCount
	)
	let c = sum(pcits)
	let p = count(pcits)
	update j._key with { citationCount: c, paperCount: p } in journals`

	updateConfIns = `for ci in confInstances
	let pcits = (
		for p in 1 inbound ci._id publish_on_confIns
			return p.citationCount
	)
	let c = sum(pcits)
	let p = count(pcits)
	update ci._key with { citationCount: c, paperCount: p } in confInstances`

	updateConfSer = `for cs in confSeries
	let pcits = (
		for p in 2 inbound cs._id publish_on_confIns, confIns_belong_to_confSer
			return p.citationCount
	)
	let c = sum(pcits)
	let p = count(pcits)
	update cs._key with { citationCount: c, paperCount: p } in confSeries`

	updateField = `for f in fields
	let c = sum(
		for v in 1 inbound f._id jou_belong_to_field, confSer_belong_to_field
			return v.citationCount
	)
	let p = sum(
		for v in 1 inbound f._id jou_belong_to_field, confSer_belong_to_field
			return v.paperCount
	)
	update f._key with { citationCount: c, paperCount: p } in fields`
)

func updateJouCountPYear() {
	jous, err := jouModel.List(0, 1000)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	for _, j := range jous {
		if j.CitCountPYear != nil {
			continue
		}
		err := j.UpdCountPYear()
		if err != nil {
			pkg.Log.Fatal(err)
		}
	}
}

func updateConfSerCountPYear() {
	confSers, err := confSerModel.List(0, 1000)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	for _, c := range confSers {
		if c.CitCountPYear != nil {
			continue
		}
		if err := c.UpdCountPYear(); err != nil {
			pkg.Log.Fatal(err)
		}
	}
}

func updateFieldCountPYear() {
	fields, err := fieldModel.List(0, 1000)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	for _, field := range fields {
		if field.CitCountPYear != nil {
			continue
		}
		if err := field.UpdCountPYear(); err != nil {
			pkg.Log.Fatal(err)
		}
	}
}

func updateAuthorCount() {
	for {
		authors, err := authorModel.List(0, 10000)
		if err != nil {
			pkg.Log.Fatal(err)
		}
		for _, author := range authors {
			if author.Key == "182-6438" {
				return
			}
			if err := author.UpdCount(); err != nil {
				pkg.Log.Fatal(err)
			}
		}
	}
}

func updateAuthorCountPYear() {
	for i := 0; ; i++ {
		pkg.Log.Infof("====== %d ======", i)
		authors, err := authorModel.List(0, 10000)
		if err != nil {
			pkg.Log.Fatal(err)
		}
		for _, author := range authors {
			if author.Key == "182-6438" {
				return
			}
			if err := author.UpdCountPYear(); err != nil {
				pkg.Log.Fatal(err)
			}
		}
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
	updateAuthorCountPYear()
}

// 清除指定集合中的全部数据
package main

import (
	"context"
	"fmt"

	"douCSAce/pkg"
)

const (
	targetColName = "cit_by" // 指定集合名。若为空，则清除所有集合中的数据
)

func truncateCol(colName string) {
	query := fmt.Sprintf("FOR doc IN %s REMOVE doc IN %s", colName, colName)
	ctx := context.Background()
	cursor, err := pkg.DB.Database.Query(ctx, query, nil)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	cursor.Close()
	pkg.Log.Infof("collection %s truncate successfully", colName)
}

// 初始化：读取配置、启动日志、连接数据库
func init() {
	pkg.Conf = pkg.NewConfig("../conf.yaml")
	pkg.Log = pkg.NewLog("", pkg.DebugLevel)
	pkg.DB = pkg.NewDB(pkg.Conf.ArangoDB.Url, pkg.Conf.ArangoDB.Username,
		pkg.Conf.ArangoDB.Passwd, pkg.Conf.ArangoDB.Database, pkg.Conf.ArangoDB.ModelColNameMap)
}

func main() {
	if targetColName == "" {
		for _, v := range pkg.Conf.ArangoDB.ModelColNameMap {
			truncateCol(v)
		}
	} else {
		truncateCol(targetColName)
	}
}

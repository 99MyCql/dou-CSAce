// 清除所有集合中的数据
package main

import (
	"context"
	"fmt"

	"douCSAce/pkg"
)

// 初始化：读取配置、启动日志、连接数据库
func init() {
	pkg.Conf = pkg.NewConfig("../conf.yaml")
	pkg.Log = pkg.NewLog("", pkg.DebugLevel)
	pkg.DB = pkg.NewDB(pkg.Conf.ArangoDB.Url, pkg.Conf.ArangoDB.Username,
		pkg.Conf.ArangoDB.Passwd, pkg.Conf.ArangoDB.Database, pkg.Conf.ArangoDB.ModelColNameMap)
}

func main() {
	collections := pkg.Conf.ArangoDB.ModelColNameMap
	for _, v := range collections {
		query := fmt.Sprintf("FOR doc IN %s REMOVE doc IN %s", v, v)
		ctx := context.Background()
		cursor, err := pkg.DB.Database.Query(ctx, query, nil)
		if err != nil {
			pkg.Log.Fatal(err)
		}
		cursor.Close()
		pkg.Log.Info(fmt.Sprintf("collection %s truncate successfully", v))
	}
}

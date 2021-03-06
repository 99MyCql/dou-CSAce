// 将当前数据库的指定集合中的数据迁移到另一个数据库的源集合中
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/arangodb/go-driver"

	"douCSAce/pkg"
)

const (
	targetDBUrl          = "http://x.x.x.x:8529/"
	targetDBUsername     = ""
	targetDBPasswd       = ""
	targetDBDatabaseName = ""
	targetDBColName      = ""
	sourceDBColName      = ""
	start                = 0
)

// 初始化：读取配置、启动日志、连接数据库
func init() {
	pkg.Conf = pkg.NewConfig("../conf.yaml")
}

func main() {
	ctx := context.Background()

	// 连接目标数据库，打开指定数据库
	sourceDatabase := pkg.OpenDB(
		pkg.ConnectDB(pkg.Conf.ArangoDB.Url, pkg.Conf.ArangoDB.Username, pkg.Conf.ArangoDB.Passwd),
		targetDBDatabaseName)
	// 打开目标集合
	sourceCol, err := sourceDatabase.Collection(ctx, sourceDBColName)
	if err != nil {
		log.Fatal(err)
	}

	// 连接目标数据库
	targetClient := pkg.ConnectDB(targetDBUrl, targetDBUsername, targetDBPasswd)
	// 打开目标指定数据库
	targetDatabase := pkg.OpenDB(targetClient, targetDBDatabaseName)
	// 打开目标集合，如果不存在，则创建
	exist, err := targetDatabase.CollectionExists(ctx, targetDBColName)
	if err != nil {
		log.Fatal(err)
	}
	var targetCol driver.Collection
	if !exist {
		options := &driver.CreateCollectionOptions{}
		// 带下划线的是边集合
		if strings.Index(targetDBColName, "_") != -1 {
			options.Type = driver.CollectionTypeEdge
		}
		targetCol, err = targetDatabase.CreateCollection(ctx, targetDBColName, options)
	} else {
		targetCol, err = targetDatabase.Collection(ctx, targetDBColName)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("open collection %s successfully", targetDBColName)

	// 将源集合中的数据复制到目标集合，查询的返回上限是 1000 条记录，需以此为基本单位
	sum, _ := sourceCol.Count(ctx)
	for i := 0; i <= int(sum/1000); i++ {
		query := fmt.Sprintf("FOR d IN %s LIMIT %d, %d RETURN d", sourceDBColName, start+i*1000, 1000)
		cursor, err := sourceDatabase.Query(ctx, query, nil)
		if err != nil {
			log.Fatal(err)
		}
		for {
			var doc map[string]interface{}
			meta, err := cursor.ReadDocument(ctx, &doc)
			if driver.IsNoMoreDocuments(err) {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			log.Printf("Got doc with key '%s' from query\n", meta.Key)
			if _, err := targetCol.CreateDocument(ctx, doc); err != nil {
				log.Fatal(err)
			}
		}
		cursor.Close()
	}
}

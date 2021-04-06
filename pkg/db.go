package pkg

import (
	"context"
	"fmt"
	"log"
	"strings"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type DBInfo struct {
	Conn     driver.Connection
	Client   driver.Client
	Database driver.Database
	Cols     map[string]driver.Collection
}

// NewDB 初始化数据库连接
func NewDB(DBUrl string, username string, passwd string, database string, cols map[string]string) *DBInfo {
	db := new(DBInfo)

	// 连接数据库
	var err error
	db.Conn, err = http.NewConnection(http.ConnectionConfig{
		// The driver has a built-in connection pooling and the connection limit (ConnLimit) defaults to 32.
		Endpoints: []string{DBUrl},
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	db.Client, err = driver.NewClient(driver.ClientConfig{
		Connection:     db.Conn,
		Authentication: driver.BasicAuthentication(username, passwd),
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Print("connect to ArangoDB successfully")

	// 打开指定数据库
	ctx := context.Background()
	db.Database, err = db.Client.Database(ctx, database)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Printf("open database %s successfully", database)

	// 打开集合
	db.Cols = make(map[string]driver.Collection)
	for k, v := range cols {
		ctx := context.Background()
		exist, err := db.Database.CollectionExists(ctx, v)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		// 如果不存在，则创建
		var col driver.Collection
		if !exist {
			options := &driver.CreateCollectionOptions{}
			if strings.Index(v, "_") != -1 {
				options.Type = driver.CollectionTypeEdge
			}
			col, err = db.Database.CreateCollection(ctx, v, options)
		} else {
			col, err = db.Database.Collection(ctx, v)
		}
		if err != nil {
			log.Fatal(err)
			return nil
		}
		db.Cols[k] = col
		log.Printf("open collection %s successfully", v)
	}

	return db
}

// ComDocCreate 通用数据库文档创建操作
func ComDocCreate(data interface{}, modelName string) (string, error) {
	Log.Info(fmt.Sprintf("ready to create document: %+v", data))
	ctx := context.Background()
	docMeta, err := DB.Cols[modelName].CreateDocument(ctx, data)
	if err != nil {
		return "", err
	}
	Log.Info(fmt.Sprintf("create document successfully: %+v", docMeta))
	return docMeta.ID.String(), nil
}

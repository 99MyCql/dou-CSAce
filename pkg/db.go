package pkg

import (
	"context"
	"fmt"
	"log"
	"strings"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

const (
	SortDesc = "DESC"
	SortAsc  = "ASC"
)

type DBInfo struct {
	Client   driver.Client
	Database driver.Database
	Cols     map[string]driver.Collection // key: modelName -> value: driver.Collection
}

// ConnectDB 连接数据库
func ConnectDB(DBUrl string, username string, passwd string) driver.Client {
	var (
		err    error
		conn   driver.Connection
		client driver.Client
	)

	// 连接数据库
	conn, err = http.NewConnection(http.ConnectionConfig{
		// The driver has a built-in connection pooling and the connection limit (ConnLimit) defaults to 32.
		Endpoints: []string{DBUrl},
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	client, err = driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(username, passwd),
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Print("connect to ArangoDB successfully")
	return client
}

// OpenDB 打开指定数据库
func OpenDB(DBClient driver.Client, databaseName string) driver.Database {
	ctx := context.Background()
	database, err := DBClient.Database(ctx, databaseName)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Printf("open database %s successfully", databaseName)
	return database
}

// 打开集合，若不存在则创建
func OpenCols(database driver.Database, colNameMap map[string]string) map[string]driver.Collection {
	cols := make(map[string]driver.Collection)
	for modelName, colName := range colNameMap {
		ctx := context.Background()
		exist, err := database.CollectionExists(ctx, colName)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		var col driver.Collection
		if !exist {
			// 如果不存在，则创建
			options := &driver.CreateCollectionOptions{}
			// 带下划线的是边集合
			if strings.Index(colName, "_") != -1 {
				options.Type = driver.CollectionTypeEdge
			}
			col, err = database.CreateCollection(ctx, colName, options)
		} else {
			col, err = database.Collection(ctx, colName)
		}
		if err != nil {
			log.Fatal(err)
			return nil
		}
		cols[modelName] = col
		log.Printf("open collection %s successfully", colName)
	}
	return cols
}

// NewDB 初始化数据库连接
func NewDB(DBUrl string, username string, passwd string, databaseName string, colNameMap map[string]string) *DBInfo {
	db := new(DBInfo)
	db.Client = ConnectDB(DBUrl, username, passwd)
	db.Database = OpenDB(db.Client, databaseName)
	db.Cols = OpenCols(db.Database, colNameMap)
	return db
}

// ComCreate 通用数据库文档创建操作
func ComCreate(data interface{}, modelName string) (string, error) {
	Log.Info(fmt.Sprintf("ready to create document: %+v", data))
	ctx := context.Background()
	docMeta, err := DB.Cols[modelName].CreateDocument(ctx, data)
	if err != nil {
		Log.Error(err)
		return "", err
	}
	Log.Info(fmt.Sprintf("create document successfully: %+v", docMeta))
	return docMeta.ID.String(), nil
}

// ComUpdate
func ComUpdate(modelName string, key string, data interface{}) error {
	Log.Info(fmt.Sprintf("ready to update document %s/%s: %+v", modelName, key, data))
	ctx := context.Background()
	docMeta, err := DB.Cols[modelName].UpdateDocument(ctx, key, data)
	if err != nil {
		Log.Error(err)
		return err
	}
	Log.Info(fmt.Sprintf("update document successfully: %s", docMeta.ID.String()))
	return nil
}

// ComList
func ComList(query string, count uint64) ([]map[string]interface{}, error) {
	Log.Info(fmt.Sprintf("ready to query: %s", query))
	ctx := context.Background()
	// 默认返回 1000 条数据，若大于则需设置 BatchSize
	if count == 0 {
		ctx = driver.WithQueryFullCount(ctx)
	} else if count > 1000 {
		ctx = driver.WithQueryBatchSize(ctx, int(count))
	}
	cursor, err := DB.Database.Query(ctx, query, nil)
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	defer cursor.Close()

	var list []map[string]interface{}
	for cursor.HasMore() {
		var tmp map[string]interface{}
		if _, err := cursor.ReadDocument(ctx, &tmp); err != nil {
			Log.Error(err)
			return nil, err
		}
		list = append(list, tmp)
	}
	Log.Info(fmt.Sprintf("query successfully: %s", query))
	return list, nil
}

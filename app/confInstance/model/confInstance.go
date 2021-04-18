package model

import (
	"context"
	"fmt"
	"strings"

	"douCSAce/pkg"
)

// ConfInstance 会议实例（哪一年哪一场的会议）
type ConfInstance struct {
	ID            string `json:"-"`    // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key           string `json:"_key"` // 唯一标识，等同于 dblp 中会议实例的 key ，_key = dblpKey.replaceAll("/", "-")，比如：conf-ppopp-2020
	Title         string `json:"title"`
	Publisher     string `json:"publisher"`
	BookTitle     string `json:"bookTitle"`
	Year          string `json:"year"`
	Isbn          string `json:"isbn"` // dblp 中会议实例的属性，用处未知，暂且保存下来
	DoiUrl        string `json:"doiUrl"`
	DblpUrl       string `json:"dblpUrl"`
	PaperCount    uint64 `json:"paperCount"`
	CitationCount uint64 `json:"citationCount"`
}

// Create 在数据库中创建数据
func (c *ConfInstance) Create() error {
	var err error
	if c.ID, err = pkg.ComDocCreate(c, pkg.ConfInstanceName); err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (c *ConfInstance) Delete() error {
	_, err := pkg.DB.Cols[pkg.ConfInstanceName].RemoveDocument(nil, c.Key)
	return err
}

// DeleteConfInsBelongToConfSer 删除所有与 confInstance 关联的 ConfInsBelongToConfSer 边
func (c *ConfInstance) DeleteConfInsBelongToConfSer() error {
	ctx := context.Background()
	query := fmt.Sprintf("for ci, ci2cs in outbound '%s' %s remove ci2cs in %s", c.ID,
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.ConfInsBelongToConfSerName],
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.ConfInsBelongToConfSerName])
	_, err := pkg.DB.Database.Query(ctx, query, nil)
	return err
}

// GenKey 返回 Key，dblpKey 为 dblp 中 conference instance 的 key 属性
func GenKey(dblpKey string) string {
	return strings.ReplaceAll(dblpKey, "/", "-")
}

// IsExist 判断 key 是否已存在
func IsExist(key string) (bool, error) {
	return pkg.DB.Cols[pkg.ConfInstanceName].DocumentExists(nil, key)
}

// FindByKey 通过 key 查找
func FindByKey(key string) (*ConfInstance, error) {
	c := new(ConfInstance)
	docMeta, err := pkg.DB.Cols[pkg.ConfInstanceName].ReadDocument(nil, key, c)
	if err != nil {
		return nil, err
	}
	c.ID = docMeta.ID.String()
	return c, nil
}

// Count 返回该集合中文档（记录）数量
func Count() (int64, error) {
	ctx := context.Background()
	return pkg.DB.Cols[pkg.ConfInstanceName].Count(ctx)
}

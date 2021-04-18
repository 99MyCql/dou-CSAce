package model

import (
	"context"
	"fmt"
	"strings"

	"douCSAce/pkg"
)

// ConfSeries 会议
type ConfSeries struct {
	ID            string `json:"-"`    // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key           string `json:"_key"` // 唯一标识，等同于 dblp 中期刊的 key ，_key = shortName ，比如：ppopp
	Name          string `json:"name"`
	ShortName     string `json:"shortName"`
	Publisher     string `json:"publisher"`
	Url           string `json:"url"`
	PaperCount    uint64 `json:"paperCount"`
	CitationCount uint64 `json:"citationCount"`
}

// Create 在数据库中创建数据
func (c *ConfSeries) Create() error {
	var err error
	if c.ID, err = pkg.ComDocCreate(c, pkg.ConfSeriesName); err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (c *ConfSeries) Delete() error {
	_, err := pkg.DB.Cols[pkg.ConfSeriesName].RemoveDocument(nil, c.Key)
	return err
}

// DeleteConfSerBelongToField 删除所有与 confSeries 关联的 ConfSerBelongToField 边
func (c *ConfSeries) DeleteConfSerBelongToField() error {
	ctx := context.Background()
	query := fmt.Sprintf("for cs, cs2f in outbound '%s' %s remove cs2f in %s", c.ID,
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.ConfSerBelongToFieldName],
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.ConfSerBelongToFieldName])
	_, err := pkg.DB.Database.Query(ctx, query, nil)
	return err
}

// GenKey 需传入会议实体的 ShortName 属性
func GenKey(shortName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(shortName, "/", "-"), " ", "-")
}

// IsExist 判断 key 是否已存在
func IsExist(key string) (bool, error) {
	return pkg.DB.Cols[pkg.ConfSeriesName].DocumentExists(nil, key)
}

// FindByKey 通过 key 查找
func FindByKey(key string) (*ConfSeries, error) {
	c := new(ConfSeries)
	docMeta, err := pkg.DB.Cols[pkg.ConfSeriesName].ReadDocument(nil, key, c)
	if err != nil {
		return nil, err
	}
	c.ID = docMeta.ID.String()
	return c, nil
}

// Count 返回该集合中文档（记录）数量
func Count() (int64, error) {
	ctx := context.Background()
	return pkg.DB.Cols[pkg.ConfSeriesName].Count(ctx)
}

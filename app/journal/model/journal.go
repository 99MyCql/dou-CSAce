package model

import (
	"context"
	"fmt"
	"strings"

	"douCSAce/pkg"
)

// Journal 期刊
type Journal struct {
	ID            string `json:"-"`    // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key           string `json:"_key"` // 唯一标识，等同于 dblp 中期刊的 key ，_key = shortName ，比如：tocs
	Name          string `json:"name"`
	ShortName     string `json:"shortName"`
	Publisher     string `json:"publisher"`
	Url           string `json:"url"`
	PaperCount    uint64 `json:"paperCount"`
	CitationCount uint64 `json:"citationCount"`
}

// Create 在数据库中创建数据
func (j *Journal) Create() error {
	var err error
	if j.ID, err = pkg.ComDocCreate(j, pkg.JournalName); err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (j *Journal) Delete() error {
	_, err := pkg.DB.Cols[pkg.JournalName].RemoveDocument(nil, j.Key)
	return err
}

// DeleteJouBelongToField 删除所有与 journal 关联的 JouBelongToField 边
func (j *Journal) DeleteJouBelongToField() error {
	ctx := context.Background()
	query := fmt.Sprintf("for j, j2f in outbound '%s' %s remove j2f in %s", j.ID,
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.JouBelongToFieldName], pkg.Conf.ArangoDB.ModelColNameMap[pkg.JouBelongToFieldName])
	_, err := pkg.DB.Database.Query(ctx, query, nil)
	return err
}

// GenKey 需传入 Journal 的 ShortName 属性
func GenKey(shortName string) string {
	return strings.ReplaceAll(shortName, " ", "-")
}

// IsExist 判断 key 是否已存在
func IsExist(key string) (bool, error) {
	return pkg.DB.Cols[pkg.JournalName].DocumentExists(nil, key)
}

// FindByKey 通过 key 查找
func FindByKey(key string) (*Journal, error) {
	j := new(Journal)
	docMeta, err := pkg.DB.Cols[pkg.JournalName].ReadDocument(nil, key, j)
	if err != nil {
		return nil, err
	}
	j.ID = docMeta.ID.String()
	return j, nil
}

// Count 返回该集合中文档（记录）数量
func Count() (int64, error) {
	ctx := context.Background()
	return pkg.DB.Cols[pkg.JournalName].Count(ctx)
}

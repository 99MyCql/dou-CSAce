package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	authorModel "douCSAce/app/author/model"
	paperModel "douCSAce/app/paper/model"
	"douCSAce/pkg"
)

// ConfInstance 会议实例（哪一年哪一场的会议）
type ConfInstance struct {
	ID            string `json:"_id,omitempty"` // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key           string `json:"_key"`          // 唯一标识，等同于 dblp 中会议实例的 key ，_key = dblpKey.replaceAll("/", "-")，比如：conf-ppopp-2020
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
	if c.ID, err = pkg.ComCreate(c, pkg.ConfInstanceName); err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (c *ConfInstance) Delete() error {
	_, err := pkg.DB.Cols[pkg.ConfInstanceName].RemoveDocument(nil, c.Key)
	return err
}

// Update 更新数据
func (c *ConfInstance) Update(updateData map[string]interface{}) error {
	if updateData == nil {
		return pkg.ComUpdate(pkg.ConfInstanceName, c.Key, c)
	}
	return pkg.ComUpdate(pkg.ConfInstanceName, c.Key, updateData)
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

// ListPaper
func (c *ConfInstance) ListPaper(offset uint64, count uint64, sortAttr string, sortType string) (
	[]*paperModel.Paper, error) {
	limitQuery := ""
	if count != 0 {
		limitQuery = fmt.Sprintf("limit %d, %d", offset, count)
	}
	sortQuery := ""
	if sortAttr != "" {
		sortQuery = fmt.Sprintf("sort p.%s %s", sortAttr, sortType)
	}
	query := fmt.Sprintf(`for p in 1 inbound '%s' publish_on_confIns
	%s %s return p`, c.ID, sortQuery, limitQuery)
	data, err := pkg.ComList(query, count)
	if err != nil {
		return nil, err
	}
	var papers []*paperModel.Paper
	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &papers)
	return papers, err
}

// ListAuthor
func (c *ConfInstance) ListAuthor(offset uint64, count uint64, sortAttr string, sortType string) (
	[]*authorModel.Author, error) {
	limitQuery := ""
	if count != 0 {
		limitQuery = fmt.Sprintf("limit %d, %d", offset, count)
	}
	sortQuery := ""
	if sortAttr != "" {
		sortQuery = fmt.Sprintf("sort author.%s %s", sortAttr, sortType)
	}
	query := fmt.Sprintf(`for p in 1 inbound '%s' publish_on_confIns
	for a in outbound p._id write_by
		COLLECT author = a
		%s %s return author`, c.ID, sortQuery, limitQuery)
	data, err := pkg.ComList(query, count)
	if err != nil {
		return nil, err
	}
	var authors []*authorModel.Author
	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &authors)
	return authors, err
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

// List
func List(offset uint64, count uint64) ([]*ConfInstance, error) {
	query := fmt.Sprintf("FOR d IN %s LIMIT %d, %d RETURN d",
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.ConfInstanceName], offset, count)
	data, err := pkg.ComList(query, count)
	var confIns []*ConfInstance
	b, _ := json.Marshal(&data)
	err = json.Unmarshal(b, &confIns)
	return confIns, err
}

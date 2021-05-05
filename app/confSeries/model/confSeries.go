package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	authorModel "douCSAce/app/author/model"
	confInsModel "douCSAce/app/confInstance/model"
	paperModel "douCSAce/app/paper/model"
	"douCSAce/pkg"
)

// ConfSeries 会议
type ConfSeries struct {
	ID              string            `json:"_id,omitempty"` // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key             string            `json:"_key"`          // 唯一标识，等同于 dblp 中期刊的 key ，_key = shortName ，比如：ppopp
	Name            string            `json:"name"`
	ShortName       string            `json:"shortName"`
	Publisher       string            `json:"publisher"`
	Url             string            `json:"url"`
	Category        string            `json:"category"`
	PaperCount      uint64            `json:"paperCount"`
	CitationCount   uint64            `json:"citationCount"`
	PaperCountPYear map[string]uint64 `json:"paperCountPYear"`
	CitCountPYear   map[string]uint64 `json:"citCountPYear"`
}

// Create 在数据库中创建数据
func (c *ConfSeries) Create() error {
	var err error
	if c.ID, err = pkg.ComCreate(c, pkg.ConfSeriesName); err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (c *ConfSeries) Delete() error {
	_, err := pkg.DB.Cols[pkg.ConfSeriesName].RemoveDocument(nil, c.Key)
	return err
}

// Update 更新数据
func (c *ConfSeries) Update(updateData map[string]interface{}) error {
	if updateData == nil {
		return pkg.ComUpdate(pkg.ConfSeriesName, c.Key, c)
	}
	return pkg.ComUpdate(pkg.ConfSeriesName, c.Key, updateData)
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

// UpdCountPYear 更新每年的引用数和论文数
func (c *ConfSeries) UpdCountPYear() error {
	query := fmt.Sprintf(`for p in 2 inbound '%s' publish_on_confIns, confIns_belong_to_confSer
    return {'citationCount':p.citationCount, 'year':p.year}`, c.ID)
	data, err := pkg.ComList(query, 0)
	if err != nil {
		return err
	}
	c.CitCountPYear = make(map[string]uint64)
	c.PaperCountPYear = make(map[string]uint64)
	for _, tmp := range data {
		if _, ok := c.CitCountPYear[tmp["year"].(string)]; !ok {
			c.CitCountPYear[tmp["year"].(string)] = 0
		}
		c.CitCountPYear[tmp["year"].(string)] += uint64(tmp["citationCount"].(float64))
		if _, ok := c.PaperCountPYear[tmp["year"].(string)]; !ok {
			c.PaperCountPYear[tmp["year"].(string)] = 0
		}
		c.PaperCountPYear[tmp["year"].(string)]++
	}
	if err := c.Update(map[string]interface{}{
		"paperCountPYear": c.PaperCountPYear,
		"citCountPYear":   c.CitCountPYear}); err != nil {
		return err
	}
	return nil
}

// ListPaper
func (c *ConfSeries) ListPaper(offset uint64, count uint64, sortAttr string, sortType string) (
	[]*paperModel.Paper, error) {
	limitQuery := ""
	if count != 0 {
		limitQuery = fmt.Sprintf("limit %d, %d", offset, count)
	}
	sortQuery := ""
	if sortAttr != "" {
		sortQuery = fmt.Sprintf("sort p.%s %s", sortAttr, sortType)
	}
	query := fmt.Sprintf(`for p in 2 inbound '%s' publish_on_confIns, confIns_belong_to_confSer
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
func (c *ConfSeries) ListAuthor(offset uint64, count uint64, sortAttr string, sortType string) (
	[]*authorModel.Author, error) {
	limitQuery := ""
	if count != 0 {
		limitQuery = fmt.Sprintf("limit %d, %d", offset, count)
	}
	sortQuery := ""
	if sortAttr != "" {
		sortQuery = fmt.Sprintf("sort author.%s %s", sortAttr, sortType)
	}
	query := fmt.Sprintf(`for p in 2 inbound '%s' publish_on_confIns, confIns_belong_to_confSer
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

// ListConfIns
func (c *ConfSeries) ListConfIns(offset uint64, count uint64, sortAttr string, sortType string) (
	[]*confInsModel.ConfInstance, error) {
	limitQuery := ""
	if count != 0 {
		limitQuery = fmt.Sprintf("limit %d, %d", offset, count)
	}
	sortQuery := ""
	if sortAttr != "" {
		sortQuery = fmt.Sprintf("sort ci.%s %s", sortAttr, sortType)
	}
	query := fmt.Sprintf(`for ci in 1 inbound '%s' confIns_belong_to_confSer
	%s %s return ci`, c.ID, sortQuery, limitQuery)
	data, err := pkg.ComList(query, count)
	if err != nil {
		return nil, err
	}
	var confInsList []*confInsModel.ConfInstance
	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &confInsList)
	return confInsList, err
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

// List
func List(offset uint64, count uint64) ([]*ConfSeries, error) {
	query := fmt.Sprintf("FOR d IN %s LIMIT %d, %d RETURN d",
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.ConfSeriesName], offset, count)
	data, err := pkg.ComList(query, count)
	var confSeries []*ConfSeries
	b, _ := json.Marshal(&data)
	err = json.Unmarshal(b, &confSeries)
	return confSeries, err
}

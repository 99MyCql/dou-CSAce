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

// Journal 期刊
type Journal struct {
	ID              string            `json:"_id,omitempty"` // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key             string            `json:"_key"`          // 唯一标识，等同于 dblp 中期刊的 key ，_key = shortName ，比如：tocs
	Name            string            `json:"name"`
	ShortName       string            `json:"shortName"`
	Publisher       string            `json:"publisher"`
	Url             string            `json:"url"`
	PaperCount      uint64            `json:"paperCount"`
	CitationCount   uint64            `json:"citationCount"`
	PaperCountPYear map[string]uint64 `json:"paperCountPYear"`
	CitCountPYear   map[string]uint64 `json:"citCountPYear"`
}

// Create 在数据库中创建数据
func (j *Journal) Create() error {
	var err error
	if j.ID, err = pkg.ComCreate(j, pkg.JournalName); err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (j *Journal) Delete() error {
	_, err := pkg.DB.Cols[pkg.JournalName].RemoveDocument(nil, j.Key)
	return err
}

// Update 更新数据
func (j *Journal) Update(updateData map[string]interface{}) error {
	if updateData == nil {
		return pkg.ComUpdate(pkg.JournalName, j.Key, j)
	}
	return pkg.ComUpdate(pkg.JournalName, j.Key, updateData)
}

// DeleteJouBelongToField 删除所有与 journal 关联的 JouBelongToField 边
func (j *Journal) DeleteJouBelongToField() error {
	ctx := context.Background()
	query := fmt.Sprintf("for j, j2f in outbound '%s' %s remove j2f in %s", j.ID,
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.JouBelongToFieldName], pkg.Conf.ArangoDB.ModelColNameMap[pkg.JouBelongToFieldName])
	_, err := pkg.DB.Database.Query(ctx, query, nil)
	return err
}

// ListPaper
func (j *Journal) ListPaper(offset uint64, count uint64, sortAttr string, sortType string) (
	[]*paperModel.Paper, error) {
	limitQuery := ""
	if count != 0 {
		limitQuery = fmt.Sprintf("limit %d, %d", offset, count)
	}
	sortQuery := ""
	if sortAttr != "" {
		sortQuery = fmt.Sprintf("sort p.%s %s", sortAttr, sortType)
	}
	query := fmt.Sprintf(`for p in 1 inbound '%s' publish_on_jou
	%s %s return p`, j.ID, sortQuery, limitQuery)
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
func (j *Journal) ListAuthor(offset uint64, count uint64, sortAttr string, sortType string) (
	[]*authorModel.Author, error) {
	limitQuery := ""
	if count != 0 {
		limitQuery = fmt.Sprintf("limit %d, %d", offset, count)
	}
	sortQuery := ""
	if sortAttr != "" {
		sortQuery = fmt.Sprintf("sort author.%s %s", sortAttr, sortType)
	}
	query := fmt.Sprintf(`for paper in 1 inbound '%s' publish_on_jou
	for author, wb in outbound paper._id write_by
		%s %s return author`, j.ID, sortQuery, limitQuery)
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

// UpdCountPYear 更新每年的引用数和论文数
func (j *Journal) UpdCountPYear() error {
	query := fmt.Sprintf(`for p in 1 inbound '%s' publish_on_jou
	return {'citationCount':p.citationCount, 'year':p.year}`, j.ID)
	data, err := pkg.ComList(query, 0)
	if err != nil {
		return err
	}
	j.CitCountPYear = make(map[string]uint64)
	j.PaperCountPYear = make(map[string]uint64)
	for _, tmp := range data {
		if _, ok := j.CitCountPYear[tmp["year"].(string)]; !ok {
			j.CitCountPYear[tmp["year"].(string)] = 0
		}
		j.CitCountPYear[tmp["year"].(string)] += uint64(tmp["citationCount"].(float64))
		if _, ok := j.PaperCountPYear[tmp["year"].(string)]; !ok {
			j.PaperCountPYear[tmp["year"].(string)] = 0
		}
		j.PaperCountPYear[tmp["year"].(string)]++
	}
	if err := j.Update(map[string]interface{}{
		"paperCountPYear": j.PaperCountPYear,
		"citCountPYear":   j.CitCountPYear}); err != nil {
		return err
	}
	return nil
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

// List
func List(offset uint64, count uint64) ([]*Journal, error) {
	query := fmt.Sprintf("FOR d IN %s LIMIT %d, %d RETURN d",
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.JournalName], offset, count)
	data, err := pkg.ComList(query, count)
	var jous []*Journal
	b, _ := json.Marshal(&data)
	err = json.Unmarshal(b, &jous)
	return jous, err
}

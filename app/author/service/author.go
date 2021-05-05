package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"douCSAce/app/author/model"
	paperModel "douCSAce/app/paper/model"
	"douCSAce/pkg"
)

type Author model.Author

// Create 在数据库中创建数据
func (a *Author) Create() error {
	var err error
	if a.ID, err = pkg.ComCreate(a, pkg.AuthorName); err != nil {
		return err
	}
	return nil
}

// Update 更新数据
func (a *Author) Update(updateData map[string]interface{}) error {
	if updateData == nil {
		return pkg.ComUpdate(pkg.AuthorName, a.Key, a)
	}
	return pkg.ComUpdate(pkg.AuthorName, a.Key, updateData)
}

// UpdCount 更新引用数和论文数
func (a *Author) UpdCount() error {
	query := fmt.Sprintf(`for p in 1 inbound '%s' write_by
	return {'citationCount':p.citationCount}`, a.ID)
	data, err := pkg.ComList(query, 0)
	if err != nil {
		return err
	}
	a.CitationCount = 0
	a.PaperCount = uint64(len(data))
	for _, tmp := range data {
		a.CitationCount += uint64(tmp["citationCount"].(float64))
	}
	if err := a.Update(map[string]interface{}{
		"paperCount":    a.PaperCount,
		"citationCount": a.CitationCount}); err != nil {
		return err
	}
	return nil
}

// UpdCountPYear 更新每年的引用数和论文数
func (a *Author) UpdCountPYear() error {
	query := fmt.Sprintf(`for p in 1 inbound '%s' write_by
	return {'citationCount':p.citationCount, 'year':p.year}`, a.ID)
	data, err := pkg.ComList(query, 0)
	if err != nil {
		return err
	}
	a.CitCountPYear = make(map[string]uint64)
	a.PaperCountPYear = make(map[string]uint64)
	for _, tmp := range data {
		if _, ok := a.CitCountPYear[tmp["year"].(string)]; !ok {
			a.CitCountPYear[tmp["year"].(string)] = 0
		}
		a.CitCountPYear[tmp["year"].(string)] += uint64(tmp["citationCount"].(float64))
		if _, ok := a.PaperCountPYear[tmp["year"].(string)]; !ok {
			a.PaperCountPYear[tmp["year"].(string)] = 0
		}
		a.PaperCountPYear[tmp["year"].(string)]++
	}
	if err := a.Update(map[string]interface{}{
		"paperCountPYear": a.PaperCountPYear,
		"citCountPYear":   a.CitCountPYear}); err != nil {
		return err
	}
	return nil
}

// ListPaper
func (a *Author) ListPaper(offset uint64, count uint64, sortAttr string, sortType string) (
	[]*paperModel.Paper, error) {
	limitQuery := ""
	if count != 0 {
		limitQuery = fmt.Sprintf("limit %d, %d", offset, count)
	}
	sortQuery := ""
	if sortAttr != "" {
		sortQuery = fmt.Sprintf("sort p.%s %s", sortAttr, sortType)
	}
	query := fmt.Sprintf(`for p in 1 inbound '%s' write_by
	%s %s return p`, a.ID, sortQuery, limitQuery)
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

// GenKey 返回 Key，需要 dblp 中 author 的 pid 属性，_key = pid.replaceAll("/", "-")，比如：journals-tocs-BalmauDZGCD20
func GenKey(dblpPid string) string {
	return strings.ReplaceAll(dblpPid, "/", "-")
}

// IsExist 判断 key 是否已存在，需要先设置 Key 属性
func IsExist(key string) (bool, error) {
	return pkg.DB.Cols[pkg.AuthorName].DocumentExists(nil, key)
}

// FindByKey 通过 key 查找
func FindByKey(key string) (*Author, error) {
	a := new(Author)
	docMeta, err := pkg.DB.Cols[pkg.AuthorName].ReadDocument(nil, key, a)
	if err != nil {
		return nil, err
	}
	a.ID = docMeta.ID.String()
	return a, nil
}

// Count 返回该集合中文档（记录）数量
func Count() (int64, error) {
	ctx := context.Background()
	return pkg.DB.Cols[pkg.AuthorName].Count(ctx)
}

// List
func List(offset uint64, count uint64) ([]*Author, error) {
	query := fmt.Sprintf("FOR d IN %s LIMIT %d, %d RETURN d",
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.AuthorName], offset, count)
	data, err := pkg.ComList(query, count)
	var authors []*Author
	b, _ := json.Marshal(&data)
	err = json.Unmarshal(b, &authors)
	return authors, err
}

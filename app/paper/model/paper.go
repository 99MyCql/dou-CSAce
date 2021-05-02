package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	authorModel "douCSAce/app/author/model"
	"douCSAce/pkg"
)

// Paper 论文实体模型
type Paper struct {
	ID             string `json:"_id,omitempty"` // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key            string `json:"_key"`          // 唯一标识，等同于 dblp 中文章的 key ，_key = dblpKey.replaceAll("/", "-")，比如：journals-tocs-BalmauDZGCD20
	Title          string `json:"title"`
	Type           uint   `json:"type"` // 0:Other, 1:Conference, 2:Journal
	Pages          string `json:"pages"`
	Year           string `json:"year"`
	BookTitle      string `json:"bookTitle"` // type=1
	Volume         string `json:"volume"`    // type=2
	Number         string `json:"number"`    // type=2
	DoiUrl         string `json:"doiUrl"`
	DblpUrl        string `json:"dblpUrl"`
	ReferenceCount uint64 `json:"referenceCount"`
	CitationCount  uint64 `json:"citationCount"`
}

// Create
func (p *Paper) Create() error {
	var err error
	if p.ID, err = pkg.ComCreate(p, pkg.PaperName); err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (p *Paper) Delete() error {
	_, err := pkg.DB.Cols[pkg.PaperName].RemoveDocument(nil, p.Key)
	return err
}

// Update 更新数据
func (p *Paper) Update(updateData map[string]interface{}) error {
	if updateData == nil {
		return pkg.ComUpdate(pkg.PaperName, p.Key, p)
	}
	return pkg.ComUpdate(pkg.PaperName, p.Key, updateData)
}

// DeletePublishOnJou 删除所有与 paper 关联的 PublishOnJou 边
func (p *Paper) DeletePublishOnJou() error {
	ctx := context.Background()
	query := fmt.Sprintf("for p, poJ in outbound '%s' %s remove poJ in %s",
		p.ID,
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.PublishOnJouName],
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.PublishOnJouName])
	_, err := pkg.DB.Database.Query(ctx, query, nil)
	return err
}

// DeletePublishOnConfIns 删除所有与 paper 关联的 PublishOnConfIns 边
func (p *Paper) DeletePublishOnConfIns() error {
	ctx := context.Background()
	query := fmt.Sprintf("for p, poCI in outbound '%s' %s remove poCI in %s", p.ID,
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.PublishOnConfInsName], pkg.Conf.ArangoDB.ModelColNameMap[pkg.PublishOnConfInsName])
	_, err := pkg.DB.Database.Query(ctx, query, nil)
	return err
}

// DeleteWriteBy 删除所有与 paper 关联的 WriteBy 边
func (p *Paper) DeleteWriteBy() error {
	ctx := context.Background()
	query := fmt.Sprintf("for p, wb in outbound '%s' %s remove wb in %s", p.ID,
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.WriteByName], pkg.Conf.ArangoDB.ModelColNameMap[pkg.WriteByName])
	_, err := pkg.DB.Database.Query(ctx, query, nil)
	return err
}

// DeleteCitBy 删除所有与 paper 关联的 CitBy 边
func (p *Paper) DeleteCitBy() error {
	ctx := context.Background()
	query := fmt.Sprintf("for p, cb in outbound '%s' %s remove cb in %s", p.ID,
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.CitByName], pkg.Conf.ArangoDB.ModelColNameMap[pkg.CitByName])
	_, err := pkg.DB.Database.Query(ctx, query, nil)
	return err
}

// GetDblpKey 返回 dblp 中的 key，比如：journals-tocs-BalmauDZGCD20 --> journals/tocs/BalmauDZGCD20
func (p *Paper) GetDblpKey() string {
	return strings.ReplaceAll(p.Key, "-", "/")
}

// ListAuthor
func (p *Paper) ListAuthor(offset uint64, count uint64) (
	[]*authorModel.Author, error) {
	query := fmt.Sprintf(`for author, wb in outbound '%s' write_by
		return author`, p.ID)
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

// GenKey 返回 Key，需传入 dblp 中 article 的 key 属性，比如：journals/tocs/BalmauDZGCD20 --> journals-tocs-BalmauDZGCD20
func GenKey(dblpKey string) string {
	return strings.ReplaceAll(dblpKey, "/", "-")
}

// IsExist 判断 key 是否已存在
func IsExist(key string) (bool, error) {
	return pkg.DB.Cols[pkg.PaperName].DocumentExists(nil, key)
}

// FindByKey 通过 key 查找
func FindByKey(key string) (*Paper, error) {
	p := new(Paper)
	docMeta, err := pkg.DB.Cols[pkg.PaperName].ReadDocument(nil, key, p)
	if err != nil {
		return nil, err
	}
	p.ID = docMeta.ID.String()
	return p, nil
}

// Count 返回该集合中文档（记录）数量
func Count() (int64, error) {
	ctx := context.Background()
	return pkg.DB.Cols[pkg.PaperName].Count(ctx)
}

// List 返回多个 Paper 文档（记录），Limit start, count
func List(offset uint64, count uint64) ([]*Paper, error) {
	query := fmt.Sprintf("FOR d IN %s LIMIT %d, %d RETURN d",
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.PaperName], offset, count)
	data, err := pkg.ComList(query, count)
	var papers []*Paper
	b, _ := json.Marshal(&data)
	err = json.Unmarshal(b, &papers)
	return papers, err
}

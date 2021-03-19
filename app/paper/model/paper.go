package model

import (
	"context"
	"fmt"
	"strings"

	"douCSAce/pkg"
)

// Paper 论文实体模型
type Paper struct {
	ID             string `json:"-"`    // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key            string `json:"_key"` // 唯一标识，等同于 dblp 中文章的 key ，_key = dblpKey.replaceAll("/", "-")，比如：journals-tocs-BalmauDZGCD20
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
	if p.ID, err = pkg.ComDocCreate(p, pkg.PaperName); err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (p *Paper) Delete() error {
	_, err := pkg.DB.Cols[pkg.PaperName].RemoveDocument(nil, p.Key)
	return err
}

// DeletePublishOnJou 删除所有与 paper 关联的 PublishOnJou 边
func (p *Paper) DeletePublishOnJou() error {
	ctx := context.Background()
	query := fmt.Sprintf("for p, poJ in outbound '%s' %s remove poJ in %s", p.ID,
		pkg.Conf.ArangoDB.ColNames[pkg.PublishOnJouName], pkg.Conf.ArangoDB.ColNames[pkg.PublishOnJouName])
	_, err := pkg.DB.Database.Query(ctx, query, nil)
	return err
}

// DeletePublishOnConfIns 删除所有与 paper 关联的 PublishOnConfIns 边
func (p *Paper) DeletePublishOnConfIns() error {
	ctx := context.Background()
	query := fmt.Sprintf("for p, poCI in outbound '%s' %s remove poCI in %s", p.ID,
		pkg.Conf.ArangoDB.ColNames[pkg.PublishOnConfInsName], pkg.Conf.ArangoDB.ColNames[pkg.PublishOnConfInsName])
	_, err := pkg.DB.Database.Query(ctx, query, nil)
	return err
}

// DeleteWriteBy 删除所有与 paper 关联的 WriteBy 边
func (p *Paper) DeleteWriteBy() error {
	ctx := context.Background()
	query := fmt.Sprintf("for p, wb in outbound '%s' %s remove wb in %s", p.ID,
		pkg.Conf.ArangoDB.ColNames[pkg.WriteByName], pkg.Conf.ArangoDB.ColNames[pkg.WriteByName])
	_, err := pkg.DB.Database.Query(ctx, query, nil)
	return err
}

// GenKey 返回 Key，需传入 dblp 中 article 的 key 属性，_key = dblpKey.replaceAll("/", "-")，比如：journals-tocs-BalmauDZGCD20
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

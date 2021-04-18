package model

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/arangodb/go-driver"

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

// Update 更新数据
func (p *Paper) Update(update map[string]interface{}) error {
	_, err := pkg.DB.Cols[pkg.PaperName].UpdateDocument(nil, p.Key, update)
	return err
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

// Count 返回该集合中文档（记录）数量
func Count() (int64, error) {
	ctx := context.Background()
	return pkg.DB.Cols[pkg.PaperName].Count(ctx)
}

// List 返回多个 Paper 文档（记录），Limit start, count
func List(start int64, count int) ([]*Paper, error) {
	if count > 1000 {
		pkg.Log.Error("一次不能获取超过1000条的数据")
		return nil, errors.New("一次不能获取超过1000条的数据")
	}

	query := fmt.Sprintf("FOR d IN %s LIMIT %d, %d RETURN d",
		pkg.Conf.ArangoDB.ModelColNameMap[pkg.PaperName], start, count)
	ctx := context.Background()
	cursor, err := pkg.DB.Database.Query(ctx, query, nil)
	if err != nil {
		pkg.Log.Error(err)
		return nil, err
	}
	defer cursor.Close()

	var papers []*Paper
	for {
		var tmp *Paper
		meta, err := cursor.ReadDocument(ctx, &tmp)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			pkg.Log.Error(err)
			return nil, err
		}
		tmp.ID = meta.ID.String()
		papers = append(papers, tmp)
	}
	return papers, nil
}

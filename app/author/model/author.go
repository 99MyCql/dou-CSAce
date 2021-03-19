package model

import (
	"strings"

	"douCSAce/pkg"
)

// Author 作者
type Author struct {
	ID            string `json:"-"`    // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key           string `json:"_key"` // 唯一标识，等同于 dblp 中作者的 pid ，_key = pid.replaceAll("/", "-")，比如：g-RajivGupta
	Name          string `json:"name"`
	ZhName        string `json:"zhName"`
	Urls          string `json:"urls"` // 作者主页，如果有多个则用空格分隔
	PaperCount    uint64 `json:"paperCount"`
	CitationCount uint64 `json:"citationCount"`
}

// Create 在数据库中创建数据
func (a *Author) Create() error {
	var err error
	if a.ID, err = pkg.ComDocCreate(a, pkg.AuthorName); err != nil {
		return err
	}
	return nil
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

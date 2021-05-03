package model

import "douCSAce/pkg"

// Affiliation 机构
type Affiliation struct {
	ID            string `json:"_id,omitempty"` // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key           string `json:"_key"`          // 唯一标识，自动生成
	Name          string `json:"name"`
	AuthorCount   uint64 `json:"authorCount"`
	PaperCount    uint64 `json:"paperCount"`
	CitationCount uint64 `json:"citationCount"`
}

// Create 在数据库中创建数据
func (a *Affiliation) Create() error {
	var err error
	if a.ID, err = pkg.ComCreate(a, pkg.AffiliationName); err != nil {
		return err
	}
	return nil
}

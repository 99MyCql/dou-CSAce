package journal

import (
	"douCSAce/pkg"
)

// Journal 期刊
type Journal struct {
	ID            string `json:"-"`    // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key           string `json:"_key"` // 唯一标识，等同于 dblp 中期刊的 key ，_key = shortName ，比如：tocs
	Name          string `json:"name"`
	ShortName     string `json:"shortName"`
	Publisher     string `json:"publisher"`
	DblpUrl       string `json:"dblpUrl"`
	PaperCount    uint   `json:"paperCount"`
	CitationCount uint   `json:"citationCount"`
}

// GetKey
func (j *Journal) GetKey() string {
	return j.ShortName
}

// Create 在数据库中创建数据
func (j *Journal) Create() error {
	var err error
	j.Key = j.GetKey()
	if j.ID, err = pkg.ComDocCreate(j, pkg.JournalName); err != nil {
		return err
	}
	return nil
}

// JouBelongToField Journal->Field 关系模型
type JouBelongToField struct {
	Key  string `json:"-"`     // 唯一标识，自动生成
	From string `json:"_from"` // From 表中对应文档的 _key
	To   string `json:"_to"`   // To 表中对应文档的 _key
	Note string `json:"note"`  // 备注，若指向的研究方向类型是 CCF ，则备注分类：A、B、C
}

// Create
func (j2f *JouBelongToField) Create() error {
	if _, err := pkg.ComDocCreate(j2f, pkg.JouBelToFieldName); err != nil {
		return err
	}
	return nil
}

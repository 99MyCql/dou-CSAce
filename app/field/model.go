package field

import (
	"strconv"
	"strings"

	"douCSAce/pkg"
)

// Field 研究方向
type Field struct {
	ID            string `json:"-"`    // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key           string `json:"_key"` // Key = <Type>-<Name.replaceAll(" ", "_")>
	Name          string `json:"name"`
	ZhName        string `json:"zhName"`
	Type          uint   `json:"type"` // 0:Other, 1:CCF
	PaperCount    uint   `json:"paperCount"`
	CitationCount uint   `json:"citationCount"`
}

// GetKey 返回 Key
func (f *Field) GetKey() string {
	return strconv.Itoa(int(f.Type)) + "-" + strings.ReplaceAll(f.Name, " ", "_")
}

// Create 在数据库中创建数据
func (f *Field) Create() error {
	var err error
	f.Key = f.GetKey()
	if f.ID, err = pkg.ComDocCreate(f, pkg.FieldName); err != nil {
		return err
	}
	return nil
}

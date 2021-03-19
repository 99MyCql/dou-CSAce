package model

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
	PaperCount    uint64 `json:"paperCount"`
	CitationCount uint64 `json:"citationCount"`
}

// Create 在数据库中创建数据
func (f *Field) Create() error {
	var err error
	if f.ID, err = pkg.ComDocCreate(f, pkg.FieldName); err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (f *Field) Delete() error {
	_, err := pkg.DB.Cols[pkg.FieldName].RemoveDocument(nil, f.Key)
	return err
}

// GenKey 返回 Key ，需先传入 Field 的 Type 和 Name 属性
func GenKey(typ uint, name string) string {
	return strconv.Itoa(int(typ)) + "-" + strings.ReplaceAll(name, " ", "_")
}

// IsExist 判断 key 是否已存在
func IsExist(key string) (bool, error) {
	return pkg.DB.Cols[pkg.FieldName].DocumentExists(nil, key)
}

// FindByKey 通过 key 查找
func FindByKey(key string) (*Field, error) {
	f := new(Field)
	docMeta, err := pkg.DB.Cols[pkg.FieldName].ReadDocument(nil, key, f)
	if err != nil {
		return nil, err
	}
	f.ID = docMeta.ID.String()
	return f, nil
}

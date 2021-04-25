package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	authorModel "douCSAce/app/author/model"
	paperModel "douCSAce/app/paper/model"
	"douCSAce/pkg"
)

// Field 研究方向
type Field struct {
	ID              string            `json:"_id,omitempty"` // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key             string            `json:"_key"`          // Key = <Type>-<Name.replaceAll(" ", "_")>
	Name            string            `json:"name"`
	ZhName          string            `json:"zhName"`
	Type            uint              `json:"type"` // 0:Other, 1:CCF
	PaperCount      uint64            `json:"paperCount"`
	CitationCount   uint64            `json:"citationCount"`
	CitCountPerYear map[string]uint64 `json:"citCountPerYear"`
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

// ListVenue 获取所属 Venue 列表
func (f *Field) ListVenue(offset uint64, count uint) ([]map[string]interface{}, error) {
	query := fmt.Sprintf(`for v in 1..1 inbound '%s' jou_belong_to_field,confSer_belong_to_field
	limit %d, %d
	return v`, f.ID, offset, count)
	return pkg.ComDocList(query, count)
}

// ListPaper
func (f *Field) ListPaper(offset uint64, count uint) ([]*paperModel.Paper, error) {
	query := fmt.Sprintf(`for v in 2..3
	inbound '%s'
	jou_belong_to_field,confSer_belong_to_field,publish_on_jou,publish_on_confIns,confIns_belong_to_confSer
	limit %d, %d
	return distinct v`, f.ID, offset, count)
	data, err := pkg.ComDocList(query, count)
	var papers []*paperModel.Paper
	b, _ := json.Marshal(&data)
	err = json.Unmarshal(b, &papers)
	return papers, err
}

// ListAuthor
func (f *Field) ListAuthor(offset uint64, count uint) ([]*authorModel.Author, error) {
	query := fmt.Sprintf(`for paper in 2..3
	inbound '%s'
	jou_belong_to_field,confSer_belong_to_field,publish_on_jou,publish_on_confIns,confIns_belong_to_confSer
	for author, wb in outbound paper._id write_by
		limit %d, %d
		return distinct author`, f.ID, offset, count)
	data, err := pkg.ComDocList(query, count)
	var authors []*authorModel.Author
	b, _ := json.Marshal(&data)
	err = json.Unmarshal(b, &authors)
	return authors, err
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

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
	PaperCountPYear map[string]uint64 `json:"paperCountPYear"`
	CitCountPYear   map[string]uint64 `json:"citCountPYear"`
}

// Create 在数据库中创建数据
func (f *Field) Create() error {
	var err error
	if f.ID, err = pkg.ComCreate(f, pkg.FieldName); err != nil {
		return err
	}
	return nil
}

// Delete 删除数据
func (f *Field) Delete() error {
	_, err := pkg.DB.Cols[pkg.FieldName].RemoveDocument(nil, f.Key)
	return err
}

// Update 更新数据
func (f *Field) Update(updateData map[string]interface{}) error {
	if updateData == nil {
		return pkg.ComUpdate(pkg.FieldName, f.Key, f)
	}
	return pkg.ComUpdate(pkg.FieldName, f.Key, updateData)
}

// UpdCountPYear 更新每年的引用数和论文数
func (f *Field) UpdCountPYear() error {
	query := fmt.Sprintf(`for v in 1 inbound '%s' jou_belong_to_field, confSer_belong_to_field
	return {'citCountPYear':v.citCountPYear, 'paperCountPYear':v.paperCountPYear}`, f.ID)
	data, err := pkg.ComList(query, 0)
	if err != nil {
		return err
	}
	f.CitCountPYear = make(map[string]uint64)
	f.PaperCountPYear = make(map[string]uint64)
	for _, tmp := range data {
		for year, citCount := range tmp["citCountPYear"].(map[string]interface{}) {
			if _, ok := f.CitCountPYear[year]; !ok {
				f.CitCountPYear[year] = 0
			}
			f.CitCountPYear[year] += uint64(citCount.(float64))
		}
		for year, paperCount := range tmp["paperCountPYear"].(map[string]interface{}) {
			if _, ok := f.PaperCountPYear[year]; !ok {
				f.PaperCountPYear[year] = 0
			}
			f.PaperCountPYear[year] += uint64(paperCount.(float64))
		}
	}
	if err := f.Update(map[string]interface{}{
		"paperCountPYear": f.PaperCountPYear,
		"citCountPYear":   f.CitCountPYear}); err != nil {
		return err
	}
	return nil
}

// ListVenue 获取所属 Venue 列表
func (f *Field) ListVenue(offset uint64, count uint64, sortAttr string, sortType string) (
	[]map[string]interface{}, error) {
	limitQuery := ""
	if count != 0 {
		limitQuery = fmt.Sprintf("limit %d, %d", offset, count)
	}
	sortQuery := ""
	if sortAttr != "" {
		sortQuery = fmt.Sprintf("sort v.%s %s", sortAttr, sortType)
	}
	query := fmt.Sprintf(`for v in 1..1 inbound '%s' jou_belong_to_field,confSer_belong_to_field
	%s %s return v`, f.ID, sortQuery, limitQuery)
	return pkg.ComList(query, count)
}

// ListPaper
func (f *Field) ListPaper(offset uint64, count uint64, sortAttr string, sortType string) (
	[]*paperModel.Paper, error) {
	limitQuery := ""
	if count != 0 {
		limitQuery = fmt.Sprintf("limit %d, %d", offset, count)
	}
	sortQuery := ""
	if sortAttr != "" {
		sortQuery = fmt.Sprintf("sort p.%s %s", sortAttr, sortType)
	}
	query := fmt.Sprintf(`for p in 2..3 inbound '%s'
	jou_belong_to_field,confSer_belong_to_field,publish_on_jou,publish_on_confIns,confIns_belong_to_confSer
	%s %s return p`, f.ID, sortQuery, limitQuery)
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

// ListAuthor
func (f *Field) ListAuthor(offset uint64, count uint64, sortAttr string, sortType string) (
	[]*authorModel.Author, error) {
	limitQuery := ""
	if count != 0 {
		limitQuery = fmt.Sprintf("limit %d, %d", offset, count)
	}
	sortQuery := ""
	if sortAttr != "" {
		sortQuery = fmt.Sprintf("sort author.%s %s", sortAttr, sortType)
	}
	query := fmt.Sprintf(`for p in 2..3 inbound '%s'
	jou_belong_to_field,confSer_belong_to_field,publish_on_jou,publish_on_confIns,confIns_belong_to_confSer
	for author, wb in outbound p._id write_by
		%s %s return author`, f.ID, sortQuery, limitQuery)
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

// List
func List() ([]*Field, error) {
	query := fmt.Sprintf("FOR d IN %s RETURN d", pkg.Conf.ArangoDB.ModelColNameMap[pkg.FieldName])
	data, err := pkg.ComList(query, 0)
	var fields []*Field
	b, _ := json.Marshal(&data)
	err = json.Unmarshal(b, &fields)
	return fields, err
}

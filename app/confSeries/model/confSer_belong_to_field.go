package model

import "douCSAce/pkg"

// ConfSerBelongToField confSeries->field 从属关系模型
type ConfSerBelongToField struct {
	Key  string `json:"-"`     // 唯一标识，自动生成
	From string `json:"_from"` // From 表中对应文档的 _id
	To   string `json:"_to"`   // To 表中对应文档的 _id
	Note string `json:"note"`  // 备注，若指向的研究方向类型是 CCF ，则备注分类：A、B、C
}

// Create 在数据库中创建数据
func (cs2f *ConfSerBelongToField) Create() error {
	if _, err := pkg.ComDocCreate(cs2f, pkg.ConfSerBelongToFieldName); err != nil {
		return err
	}
	return nil
}

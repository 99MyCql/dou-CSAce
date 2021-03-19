package model

import "douCSAce/pkg"

// JouBelongToField Journal->Field 关系模型
type JouBelongToField struct {
	Key  string `json:"-"`     // 唯一标识，自动生成
	From string `json:"_from"` // From 表中对应文档的 _id
	To   string `json:"_to"`   // To 表中对应文档的 _id
	Note string `json:"note"`  // 备注，若指向的研究方向类型是 CCF ，则备注分类：A、B、C
}

// Create
func (j2f *JouBelongToField) Create() error {
	if _, err := pkg.ComDocCreate(j2f, pkg.JouBelongToFieldName); err != nil {
		return err
	}
	return nil
}

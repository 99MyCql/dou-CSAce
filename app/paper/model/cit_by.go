package model

import "douCSAce/pkg"

// CitBy paper->paper 引用关系模型
type CitBy struct {
	Key  string `json:"_id,omitempty"` // 唯一标识，自动生成
	From string `json:"_from"`         // From 表中对应文档的 _id
	To   string `json:"_to"`           // To 表中对应文档的 _id
}

// Create
func (c *CitBy) Create() error {
	if _, err := pkg.ComDocCreate(c, pkg.CitByName); err != nil {
		return err
	}
	return nil
}

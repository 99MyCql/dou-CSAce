package model

import "douCSAce/pkg"

// CoWith author->author 合作关系模型
type CoWith struct {
	Key  string `json:"_id,omitempty"` // 唯一标识，自动生成
	From string `json:"_from"`         // From 表中对应文档的 _id
	To   string `json:"_to"`           // To 表中对应文档的 _id
}

// Create 在数据库中创建数据
func (c *CoWith) Create() error {
	if _, err := pkg.ComCreate(c, pkg.CoWithName); err != nil {
		return err
	}
	return nil
}

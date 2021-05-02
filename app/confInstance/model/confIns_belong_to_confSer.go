package model

import "douCSAce/pkg"

// ConfInsBelongToConfSer confInstance->confSeries 从属关系模型
type ConfInsBelongToConfSer struct {
	Key  string `json:"_id,omitempty"` // 唯一标识，自动生成
	From string `json:"_from"`         // From 表中对应文档的 _id
	To   string `json:"_to"`           // To 表中对应文档的 _id
}

// Create 在数据库中创建数据
func (ci2cs *ConfInsBelongToConfSer) Create() error {
	if _, err := pkg.ComCreate(ci2cs, pkg.ConfInsBelongToConfSerName); err != nil {
		return err
	}
	return nil
}

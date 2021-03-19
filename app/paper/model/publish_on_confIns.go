package model

import "douCSAce/pkg"

// PublishOnConfIns paper->conferenceIns 发表关系模型
type PublishOnConfIns struct {
	Key  string `json:"-"`     // 唯一标识，自动生成
	From string `json:"_from"` // From 表中对应文档的 _id
	To   string `json:"_to"`   // To 表中对应文档的 _id
}

// Create
func (poCI *PublishOnConfIns) Create() error {
	if _, err := pkg.ComDocCreate(poCI, pkg.PublishOnConfInsName); err != nil {
		return err
	}
	return nil
}

package model

import "douCSAce/pkg"

// PublishOnJou paper->journal 发表关系模型
type PublishOnJou struct {
	Key  string `json:"_id,omitempty"` // 唯一标识，自动生成
	From string `json:"_from"`         // From 表中对应文档的 _id
	To   string `json:"_to"`           // To 表中对应文档的 _id
}

// Create
func (poJ *PublishOnJou) Create() error {
	if _, err := pkg.ComDocCreate(poJ, pkg.PublishOnJouName); err != nil {
		return err
	}
	return nil
}

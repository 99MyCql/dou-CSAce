package model

import "douCSAce/pkg"

// WriteBy paper->author 著作关系模型
type WriteBy struct {
	Key  string `json:"_id,omitempty"` // 唯一标识，自动生成
	From string `json:"_from"`         // From 表中对应文档的 _id
	To   string `json:"_to"`           // To 表中对应文档的 _id
}

// Create
func (w *WriteBy) Create() error {
	if _, err := pkg.ComCreate(w, pkg.WriteByName); err != nil {
		return err
	}
	return nil
}

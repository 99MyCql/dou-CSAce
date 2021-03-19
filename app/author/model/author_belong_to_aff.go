package model

import "douCSAce/pkg"

// AuthorBelongToAff author->affiliation 从属关系模型
type AuthorBelongToAff struct {
	Key       string `json:"-"`     // 唯一标识，自动生成
	From      string `json:"_from"` // From 表中对应文档的 _id
	To        string `json:"_to"`   // To 表中对应文档的 _id
	StartYear string `json:"startYear"`
	EndYear   string `json:"endYear"`
	Note      string `json:"note"`
}

// Create 在数据库中创建数据
func (a2a *AuthorBelongToAff) Create() error {
	if _, err := pkg.ComDocCreate(a2a, pkg.AuthorBelongToAffName); err != nil {
		return err
	}
	return nil
}

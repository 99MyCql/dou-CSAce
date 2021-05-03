package model

// Author 作者
type Author struct {
	ID              string            `json:"_id,omitempty"` // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key             string            `json:"_key"`          // 唯一标识，等同于 dblp 中作者的 pid ，_key = pid.replaceAll("/", "-")，比如：g-RajivGupta
	Name            string            `json:"name"`
	ZhName          string            `json:"zhName"`
	Urls            string            `json:"urls"` // 作者主页，如果有多个则用空格分隔
	PaperCount      uint64            `json:"paperCount"`
	CitationCount   uint64            `json:"citationCount"`
	PaperCountPYear map[string]uint64 `json:"paperCountPYear,omitempty"`
	CitCountPYear   map[string]uint64 `json:"citCountPYear,omitempty"`
}

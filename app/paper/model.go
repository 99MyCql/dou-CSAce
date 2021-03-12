package paper

import (
	"strconv"

	"douCSAce/pkg"
)

type Paper struct {
	ID             string `json:"-"`    // ArangoDB 中文档的默认属性，_id = <collection name>/<_key>
	Key            string `json:"_key"` // 唯一标识，等同于 dblp 中文章的 key ，_key = <type>-<journal _key/conf _key>-<dblpkey>，比如：1-ppopp-GavrielatosKNGJ20
	Title          string `json:"title"`
	Abstract       string `json:"abstract"`
	Type           uint   `json:"type"` // 0:Other, 1:Conference, 2:Journal
	Pages          string `json:"pages"`
	Year           string `json:"year"`
	BookTitle      string `json:"bookTitle"` // type=1
	Volume         string `json:"volume"`    // type=2
	Number         string `json:"number"`    // type=2
	DoiUrl         string `json:"doiUrl"`
	DblpUrl        string `json:"dblpUrl"`
	ReferenceCount int64  `json:"referenceCount"`
	CitationCount  int64  `json:"citationCount"`
}

// GetKey 返回 Key，_key = <type>-<journal _key/conf _key>-<dblpkey>，比如：1-ppopp-GavrielatosKNGJ20
func (p *Paper) GetKey(venueKey string, dblpKey string) string {
	return strconv.Itoa(int(p.Type)) + "-" + venueKey + "-" + dblpKey
}

// Create
func (p Paper) Create(venueKey string, dblpKey string) error {
	var err error
	p.Key = p.GetKey(venueKey, dblpKey)
	if p.ID, err = pkg.ComDocCreate(p, pkg.PaperName); err != nil {
		return err
	}
	return nil
}

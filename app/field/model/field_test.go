package model

import (
	"os"
	"testing"

	"douCSAce/pkg"
)

func TestField_Create(t *testing.T) {
	f := &Field{
		Key:           GenKey(1, "Computer Architecture"),
		Name:          "Computer Architecture",
		ZhName:        "计算机体系结构/并行与分布计算/存储系统",
		Type:          1,
		PaperCount:    0,
		CitationCount: 0,
	}
	err := f.Create()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", f)
}

func TestFindByKey(t *testing.T) {
	f, err := FindByKey("1-Computer_Architecture")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", f)
}

func TestField_ListVenue(t *testing.T) {
	f, err := FindByKey("1-Computer_Architecture")
	if err != nil {
		t.Fatal(err)
	}
	venues, err := f.ListVenue(0, 100)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(venues))
	t.Logf("%+v", venues[0])
}

func TestField_ListPaper(t *testing.T) {
	f, err := FindByKey("1-Computer_Architecture")
	if err != nil {
		t.Fatal(err)
	}
	papers, err := f.ListPaper(0, 1000)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(papers))
	t.Logf("%+v", papers[0])
}

func TestField_ListAuthor(t *testing.T) {
	f, err := FindByKey("1-Computer_Architecture")
	if err != nil {
		t.Fatal(err)
	}
	authors, err := f.ListAuthor(0, 100)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(authors))
	t.Logf("%+v", authors[0])
}

func TestMain(m *testing.M) {
	pkg.TestSetup("../../../conf.yaml")
	code := m.Run()
	os.Exit(code)
}

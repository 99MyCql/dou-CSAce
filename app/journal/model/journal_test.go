package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	fieldModel "douCSAce/app/field/model"
	"douCSAce/pkg"
)

// TestMain 包内测试入口函数
func TestMain(m *testing.M) {
	pkg.TestSetup("../../../conf.yaml")
	code := m.Run()
	os.Exit(code)
}

func TestJournal_Create(t *testing.T) {
	j := &Journal{
		Key:           GenKey("test"),
		Name:          "test 111",
		ShortName:     "test",
		Publisher:     "test",
		Url:           "test",
		PaperCount:    0,
		CitationCount: 0,
	}
	err := j.Create()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", j))
}

func TestJournal_Delete(t *testing.T) {
	j := &Journal{
		Key: "testJou",
	}
	j.Create()
	if err := j.Delete(); err != nil {
		t.Fatal(err)
	}
}

func TestJournal_DeleteJouBelongToField(t *testing.T) {
	j := &Journal{
		Key: "testJou",
	}
	j.Create()
	f := &fieldModel.Field{
		Key: "testField",
	}
	f.Create()
	j2f := &JouBelongToField{
		From: j.ID,
		To:   f.ID,
	}
	j2f.Create()

	if err := j.DeleteJouBelongToField(); err != nil {
		t.Fatal(err)
	}

	f.Delete()
	j.Delete()
}

func TestJournal_ListPaper(t *testing.T) {
	j, err := FindByKey("TOCS")
	assert.Nil(t, err)
	papers, err := j.ListPaper(0, 10, "citationCount", pkg.SortDesc)
	assert.Nil(t, err)
	assert.NotEqual(t, len(papers), 0)
	t.Log(papers[0])
}

func TestJournal_ListAuthor(t *testing.T) {
	j, err := FindByKey("TOCS")
	assert.Nil(t, err)
	authors, err := j.ListAuthor(0, 10, "citationCount", pkg.SortDesc)
	assert.Nil(t, err)
	assert.NotEqual(t, len(authors), 0)
	t.Logf("%+v", authors[0])
}

func TestJournal_UpdCountPYear(t *testing.T) {
	j, err := FindByKey("TOCS")
	assert.Nil(t, err)
	err = j.UpdCountPYear()
	assert.Nil(t, err)
	j, err = FindByKey("TOCS")
	assert.Nil(t, err)
	t.Log(j)
}

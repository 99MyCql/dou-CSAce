package model

import (
	"fmt"
	"os"
	"testing"

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
		t.Error(err)
	}
	t.Log(fmt.Sprintf("%+v", j))
}

func TestJournal_Delete(t *testing.T) {
	j := &Journal{
		Key: "testJou",
	}
	j.Create()
	if err := j.Delete(); err != nil {
		t.Error(err)
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
		t.Error(err)
	}

	f.Delete()
	j.Delete()
}

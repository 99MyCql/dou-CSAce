package model

import (
	"os"
	"testing"

	jouModel "douCSAce/app/journal/model"
	"douCSAce/pkg"
)

// TestMain 包内测试入口函数
func TestMain(m *testing.M) {
	pkg.TestSetup("../../../conf.yaml")
	code := m.Run()
	os.Exit(code)
}

func TestJournal_DeletePublishOnJou(t *testing.T) {
	p := Paper{Key: "testPaper"}
	p.Create()
	j := &jouModel.Journal{Key: "testJou"}
	j.Create()
	poJ := &PublishOnJou{
		From: p.ID,
		To:   j.ID,
	}
	poJ.Create()

	if err := p.DeletePublishOnJou(); err != nil {
		t.Fatal(err)
	}

	j.Delete()
	p.Delete()
}

func TestList(t *testing.T) {
	papers, err := List(0, 10000)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(papers))
}

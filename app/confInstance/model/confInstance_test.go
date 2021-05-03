package model

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	confSerModel "douCSAce/app/confSeries/model"
	"douCSAce/pkg"
)

// TestMain 包内测试入口函数
func TestMain(m *testing.M) {
	pkg.TestSetup("../../../conf.yaml")
	code := m.Run()
	os.Exit(code)
}

func TestConfSeries_DeleteConfInsBelongToConfSer(t *testing.T) {
	ci := &ConfInstance{Key: "testConfIns"}
	ci.Create()
	cs := &confSerModel.ConfSeries{Key: "testConfSer"}
	cs.Create()
	ci2cs := &ConfInsBelongToConfSer{
		From: ci.ID,
		To:   cs.ID,
	}
	ci2cs.Create()

	if err := ci.DeleteConfInsBelongToConfSer(); err != nil {
		t.Fatal(err)
	}

	cs.Delete()
	ci.Delete()
}

func TestConfInstance_ListAuthor(t *testing.T) {
	ci, err := FindByKey("conf-ppopp-2020")
	assert.Nil(t, err)
	authors, err := ci.ListAuthor(0, 10, "citationCount", pkg.SortDesc)
	assert.Nil(t, err)
	t.Logf("%+v", authors[0])
}

func TestConfInstance_ListPaper(t *testing.T) {
	ci, err := FindByKey("conf-ppopp-2020")
	assert.Nil(t, err)
	papers, err := ci.ListPaper(0, 10, "citationCount", pkg.SortDesc)
	assert.Nil(t, err)
	t.Logf("%+v", papers[0])
}

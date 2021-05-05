package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"douCSAce/pkg"
)

// TestMain 包内测试入口函数
func TestMain(m *testing.M) {
	pkg.TestSetup("../../../conf.yaml")
	code := m.Run()
	os.Exit(code)
}

func TestAuthor_UpdCount(t *testing.T) {
	a, err := FindByKey("182-6438")
	assert.Nil(t, err)
	err = a.UpdCount()
	assert.Nil(t, err)
	a, err = FindByKey("182-6438")
	assert.Nil(t, err)
	t.Log(a)
}

func TestAuthor_UpdCountPYear(t *testing.T) {
	a, err := FindByKey("182-6438")
	assert.Nil(t, err)
	err = a.UpdCountPYear()
	assert.Nil(t, err)
	a, err = FindByKey("182-6438")
	assert.Nil(t, err)
	t.Log(a)
}

func TestAuthor_ListPaper(t *testing.T) {
	a, err := FindByKey("182-6438")
	assert.Nil(t, err)
	papers, err := a.ListPaper(0, 10, "citationCount", pkg.SortDesc)
	assert.Nil(t, err)
	t.Logf("%+v", papers[0])
}

package model

import (
	"fmt"
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
		t.Error(err)
	}
	t.Log(fmt.Sprintf("%+v", f))
}

func TestMain(m *testing.M) {
	pkg.TestSetup("../../../conf.yaml")
	code := m.Run()
	os.Exit(code)
}

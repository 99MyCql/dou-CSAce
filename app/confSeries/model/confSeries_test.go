package model

import (
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

func TestConfSeries_DeleteConfSerBelongToField(t *testing.T) {
	cs := &ConfSeries{Key: "testConfSer"}
	cs.Create()
	f := &fieldModel.Field{Key: "testField"}
	f.Create()
	cs2f := &ConfSerBelongToField{
		From: cs.ID,
		To:   f.ID,
		Note: "",
	}
	cs2f.Create()

	if err := cs.DeleteConfSerBelongToField(); err != nil {
		t.Fatal(err)
	}

	f.Delete()
	cs.Delete()
}

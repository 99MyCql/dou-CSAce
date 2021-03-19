package model

import (
	"os"
	"testing"

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
		t.Error(err)
	}

	cs.Delete()
	ci.Delete()
}

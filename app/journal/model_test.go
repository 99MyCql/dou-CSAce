package journal

import (
	"fmt"
	"os"
	"testing"

	"douCSAce/pkg"
)

func TestJournal_Create(t *testing.T) {
	j := &Journal{
		Name:          "test 111",
		ShortName:     "test",
		Publisher:     "test",
		DblpUrl:       "test",
		PaperCount:    0,
		CitationCount: 0,
	}
	err := j.Create()
	if err != nil {
		t.Error(err)
	}
	t.Log(fmt.Sprintf("%+v", j))
}

func TestJouBelongToField_Create(t *testing.T) {
	j2f := &JouBelongToField{
		From: "journals/test",
		To:   "fields/1-Computer_Architecture",
		Note: "A",
	}
	err := j2f.Create()
	if err != nil {
		t.Error(err)
	}
	t.Log(j2f)
}

func TestMain(m *testing.M) {
	pkg.TestSetup("../../conf.yaml")
	code := m.Run()
	os.Exit(code)
}

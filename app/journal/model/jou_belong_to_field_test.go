package model

import (
	"testing"
)

func TestJouBelongToField_Create(t *testing.T) {
	j2f := &JouBelongToField{
		From: "journals/test",
		To:   "fields/1-Computer_Architecture",
		Note: "A",
	}
	err := j2f.Create()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(j2f)
}

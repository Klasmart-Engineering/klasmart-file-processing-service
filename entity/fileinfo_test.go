package entity

import "testing"

func TestFileInfo(t *testing.T) {
	f := ParseFileInfo("kfps:assessments", "/partition/goodbye.jpg")
	t.Logf("%#v", f)
}

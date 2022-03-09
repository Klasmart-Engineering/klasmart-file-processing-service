package entity

import "testing"

func TestFileInfo(t *testing.T) {
	f := ParseFileInfo("/partition/goodbye.jpg")
	t.Logf("%#v", f)
}

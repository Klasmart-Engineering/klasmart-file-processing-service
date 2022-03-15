package entity

import "testing"

var testFilePath = "../piiFile.jpg"

func TestFileInfo(t *testing.T) {
	f := ParseFileInfo(testFilePath)
	if "jpg" != f.Extension || "piiFile.jpg" != f.Name || testFilePath != f.Path {
		t.Fail()
	}
	t.Logf("%#v", f)
}

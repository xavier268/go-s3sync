package gosync

import "testing"

func TestUpload(t *testing.T) {
	c := NewDefaultConfig()
	c.UploadFile(SrcFile{absPath: "/home/xavier/Desktop/test/ttt"})
}

package gosync

import "testing"

func TestUpload(t *testing.T) {
	c := NewDefaultConfig()
	c.uploadFile(SrcFile{absPath: "/home/xavier/Desktop/test/ttt"})
}

package synchro

import (
	"fmt"
	"time"
)

// SrcFile capture useful info about the file.
type SrcFile struct {
	abspath string
	updated time.Time
}

// String provides a human-readable string.
func (sf *SrcFile) String() string {
	s := fmt.Sprintf("[%v]\t%s", sf.updated, sf.abspath)
	return s
}

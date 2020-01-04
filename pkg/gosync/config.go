package gosync

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Config defines the configuration and context for a sync operation.
type Config struct {
	// bucket name
	bucket string
	// prefix to define the root of the file system to consider
	// also used as an implicit prefix for s3 keys
	prefix string
	// aws region
	region string
	// max key length - 1024 as per aws documentation
	maxKeyLength int

	// S3 client
	s3 *s3.S3

	// Channel for processing source files
	files chan SrcFile
	// Channel for processing S3 objects
	objects chan DstObject
	// WaitGroup used to detect when all workers are done.
	wait sync.WaitGroup
}

// SrcFile describe the source file locally.
// Dirs are ignored.
type SrcFile struct {
	absPath string
	updated time.Time
	size    int64
}

func (s *SrcFile) String() string {
	res := fmt.Sprintf("[%v]\t%d bytes\t%s", s.updated, s.size, s.absPath)
	return res
}

// GetKey remove the efix from the absPath of a SrcFile.
// Return empty key if prefix does not match.
func (c *Config) GetKey(s SrcFile) string {
	if !strings.HasPrefix(s.absPath, c.prefix) {
		return ""
	}
	return s.absPath[len(c.prefix):]
}

// DstObject describes the S3 object in the target bucket.
type DstObject struct {
	key     string
	updated time.Time
	size    int64
}

func (o *DstObject) String() string {
	res := fmt.Sprintf("[%v]\t%d bytes\t%s", o.updated, o.size, o.key)
	return res
}

// FromS3Object return a DstObject from an s3.Object
func FromS3Object(o *s3.Object) DstObject {
	if o == nil {
		panic(errors.New("cannot process a nil s3.object"))
	}
	d := DstObject{}
	d.key = *o.Key
	d.updated = o.LastModified.UTC()
	d.size = *o.Size
	return d
}

// NewTestConfig provides a test configuration
func NewTestConfig() *Config {
	c := new(Config)
	c.bucket = "bup3.photos.gandillot.com"
	c.prefix = "/home/xavier/Desktop/"
	c.region = "eu-west-1"
	c.maxKeyLength = 1000 // Name limit - real is 1024

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(c.region),
		})
	if err != nil {
		panic(err)
	}

	c.s3 = s3.New(sess)

	c.files = make(chan SrcFile, 2000)
	c.objects = make(chan DstObject, 2000)

	return c
}

package gosync

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	// also added to s3 keys to retrieve the absolute path.
	prefix string
	// aws region
	region string
	// max key length - 1024 as per aws documentation
	maxKeyLength int

	// What mode are syncing ?
	// Do we actually modify things or is it a "mock" operation ?
	mode Mode

	// permission mode used for creating intermediate directories.
	dirPerm os.FileMode

	// S3 session
	sess *session.Session
	// S3 client
	s3 *s3.S3

	// Channel for processing source files
	files chan SrcFile
	// Channel for processing S3 objects
	objects chan DstObject
}

func (c *Config) String() string {
	s := fmt.Sprintf("Configuration :\n\tMode:\t%s\n\tBucket:\t%s\n\tPrefix:\t%s\n\tRegion:\t%s\n",
		c.mode.String(), c.bucket, c.prefix, c.region)
	return s
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

// getKey remove the efix from the absPath of a SrcFile.
// Return empty key if prefix does not match.
func (c *Config) getKey(s SrcFile) string {
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

// getAbsPath constructs the absolute path equivalent.
func (o *DstObject) getAbsPath(c *Config) string {
	res := path.Join(c.prefix, o.key)
	res, err := filepath.Abs(res)
	if err != nil {
		panic(err)
	}
	return res
}

// dstObjectFromS3Object return a DstObject from an s3.Object
func (c *Config) dstObjectFromS3Object(o *s3.Object) DstObject {
	if o == nil {
		panic(errors.New("cannot process a nil s3.object"))
	}
	d := DstObject{}
	d.key = *o.Key
	d.updated = o.LastModified.UTC()
	d.size = *o.Size
	return d
}

// NewConfig creates a new configuration,
// starting with default values,
// then potentially overriden from CLI flags.
// Since this is parsing the CLI, it can only be called only once -
// and it's a feature, not a bug,
// intended to discourage simultaneous processing with different configs
func NewConfig() *Config {

	c := NewDefaultConfig()

	// Define and parse flags, overriding test values
	flag.StringVar(&c.bucket, "bucket", c.bucket, "the s3 bucket used to save the selected files")
	flag.StringVar(&c.bucket, "b", c.bucket, "the s3 bucket used to save the selected files")

	flag.StringVar(&c.prefix, "prefix", c.prefix, "the file directory to synchronize")
	flag.StringVar(&c.prefix, "p", c.prefix, "the file directory to synchronize")

	flag.StringVar(&c.region, "region", c.region, "the AWS region to use")

	flag.Parse()

	ap, err := filepath.Abs(c.prefix)
	if err != nil {
		fmt.Println("The provided prefix is invalid and could not be translated into an absolute path : ", c.prefix)
		panic(err)
	} else {
		c.prefix = ap
	}
	return c

}

// NewDefaultConfig provides a default configuration
func NewDefaultConfig() *Config {
	var err error

	c := new(Config)
	c.bucket = "test.gandillot.com"
	c.prefix = "/home/xavier/Desktop/test"
	c.region = "eu-west-1"
	c.maxKeyLength = 1000 // real limit is 1024 per AWS documentation

	c.mode = ModeBackupMock
	c.dirPerm = 0o_0777 // all permissions to anyone ...

	c.sess, err = session.NewSession(
		&aws.Config{
			Region: aws.String(c.region),
		})
	if err != nil {
		panic(err)
	}

	c.s3 = s3.New(c.sess)

	c.files = make(chan SrcFile, 2000)
	c.objects = make(chan DstObject, 2000)

	return c
}

// SetMode sets the mode for the sync operation (backup or restore)
func (c *Config) SetMode(m Mode) *Config {
	c.mode = m
	return c
}

// SetPerm sets the permission (FileMode) to use when creating missing directories.
func (c *Config) SetPerm(dirPermission os.FileMode) *Config {
	c.dirPerm = dirPermission
	return c
}

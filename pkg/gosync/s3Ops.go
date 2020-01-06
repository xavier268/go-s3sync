package gosync

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// uploadFile upload a potentially large file to S3
func (c *Config) uploadFile(sf SrcFile) {

	file, err := os.Open(sf.absPath)
	if err != nil {
		fmt.Println(sf.String())
		panic(err)
	}
	defer file.Close()

	up := s3manager.NewUploader(c.sess)
	_, err = up.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(c.getKey(sf)),
		Body:   file,
	})
	if err != nil {
		panic(err)
	}
}

// deleteFile does just that ...
func (c *Config) deleteFile(sf SrcFile) {
	err := os.Remove(sf.absPath)
	if err != nil {
		panic(err)
	}
}

// downloadFile downloads a potentially large object from S3 to file,
// overwriting existing file.
func (c *Config) downloadFile(sf SrcFile) {

	file, err := os.Create(sf.absPath)
	if err != nil {
		fmt.Println(sf.String())
		panic(err)
	}
	defer file.Close()

	down := s3manager.NewDownloader(c.sess)
	_, err = down.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(c.getKey(sf)),
		})

	if err != nil {
		fmt.Println(sf)
		panic(err)
	}
}

// deleteObject delete the provided object from s3
func (c *Config) deleteObject(ob DstObject) {

	_, err := c.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(ob.key),
	})

	if err != nil {
		panic(err)
	}

}

// uploadObject refresh the s3 object from corresponding file
func (c *Config) uploadObject(ob DstObject) {
	c.uploadFile(
		SrcFile{
			absPath: ob.getAbsPath(c),
		})
}

// downloadObject downloads an S3 object to the local file system.
func (c *Config) downloadObject(ob DstObject) {
	c.downloadFile(
		SrcFile{
			absPath: ob.getAbsPath(c),
		})

}

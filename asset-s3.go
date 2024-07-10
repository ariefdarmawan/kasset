package kasset

import (
	"bytes"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Asset struct {
	//key, secret, token, bucket string
	bucket string
	svc    *s3.S3
}

func NewS3(key, secret, token, bucket string) (*S3Asset, error) {
	cfg := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(key, secret, "")).
		WithRegion("ap-southeast-1")

	return NewS3WithConfig(bucket, cfg)
}

func NewS3WithConfig(bucket string, config *aws.Config) (*S3Asset, error) {
	fs := new(S3Asset)
	fs.bucket = bucket

	sess, e := session.NewSession(config)
	if e != nil {
		return nil, e
	}

	s3svc := s3.New(sess)
	listInput := &s3.ListObjectsInput{Bucket: aws.String(bucket), MaxKeys: aws.Int64(16)}
	if _, e := s3svc.ListObjects(listInput); e != nil {
		createInput := &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		}
		if _, e = s3svc.CreateBucket(createInput); e != nil {
			// fmt.Println(e.Error())
			if !strings.Contains(e.Error(), s3.ErrCodeBucketAlreadyOwnedByYou) {
				return nil, e
			}

		}
	}

	fs.svc = s3svc
	return fs, nil
}

func (sfs *S3Asset) getService() *s3.S3 {
	if sfs.svc == nil {
		sess, _ := session.NewSession()
		sfs.svc = s3.New(sess)
	}
	return sfs.svc
}

func (sfs *S3Asset) Save(name string, bs []byte) error {
	byteReader := bytes.NewReader(bs)

	svc := sfs.getService()
	putInput := s3.PutObjectInput{
		Bucket: aws.String(sfs.bucket),
		Body:   aws.ReadSeekCloser(byteReader),
		Key:    aws.String(name),
	}
	_, err := svc.PutObject(&putInput)
	if err != nil {
		return err
	}
	return nil
}

func (sfs *S3Asset) Read(name string) ([]byte, error) {
	bs := []byte{}
	buffer := bytes.NewBuffer(bs)

	svc := sfs.getService()
	readParam := s3.GetObjectInput{
		Bucket: aws.String(sfs.bucket),
		Key:    aws.String(name),
	}
	readResult, err := svc.GetObject(&readParam)
	if err != nil {
		return []byte{}, err
	}
	defer readResult.Body.Close()
	if _, err := io.Copy(buffer, readResult.Body); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

func (sfs *S3Asset) Delete(name string) error {
	svc := sfs.getService()
	delInput := s3.DeleteObjectInput{
		Bucket: aws.String(sfs.bucket),
		Key:    aws.String(name),
	}
	if _, err := svc.DeleteObject(&delInput); err != nil {
		return err
	}
	return nil
}

package store

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

type Store interface {
	Upload(ctx context.Context, filepath string) (downloadurl string, err error)
}

type S3Store struct {
	Region          string
	AccesskeyID     string
	SecretAccessKey string
	Bucket          string
	StoreDir        string
	Endpoint        string
	s3              *s3.S3
	uploader        *s3manager.Uploader
	PresignExpire   time.Duration
}

func (ss *S3Store) Upload(ctx context.Context, path string) (downloadurl string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	_, filename := filepath.Split(path)
	key := filepath.Join(ss.StoreDir, filename)
	_, err = ss.uploader.Upload(&s3manager.UploadInput{
		Bucket: &ss.Bucket,
		Key:    &key,
		Body:   f,
	}, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024 // Multipart upload
		u.LeavePartsOnError = true
		u.Concurrency = 3
	})
	gc := os.Getenv("GC")
	gc_type := os.Getenv("GC_TYPE")
	if err != nil {
		return
	} else if gc == "true" && gc_type == "rm" {
		_ = os.Remove(path)
	} else if gc == "true" && gc_type == "truncate" {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			fmt.Println("Error opening file:", err)
		} else {
			reader := bufio.NewReader(file)

			for {
				_, err := reader.ReadString('\n')
				if err != nil {
					break
				}
			}
		}
	}
	req, _ := ss.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(ss.Bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(ss.PresignExpire)
	if err != nil {
		err = errors.Wrap(err, "unable to sign request")
		return
	} else {
		downloadurl = urlStr
	}
	return
}

func NewS3Store(region, akid, aksecret, bucket, endpoint, storedir string, presignExpire int) (Store, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:           &region,
		Credentials:      credentials.NewStaticCredentials(akid, aksecret, ""),
		Endpoint:         &endpoint,
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false),
	}))

	svc := s3.New(sess)
	uploader := s3manager.NewUploader(sess)
	store := &S3Store{
		s3:              svc,
		Region:          region,
		AccesskeyID:     akid,
		SecretAccessKey: aksecret,
		Bucket:          bucket,
		StoreDir:        storedir,
		PresignExpire:   time.Duration(presignExpire * int(time.Second)),
	}
	store.uploader = uploader
	return store, nil
}

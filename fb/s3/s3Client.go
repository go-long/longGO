package s3

import (
	"github.com/go-long/goamz/aws"
	"github.com/go-long/goamz/s3"
	"fmt"

)


type (
	S3Client struct {
		*s3.S3
		config Config
	}

	Config struct {
	AccessKey string
	SecretKey string
	BucketName string
	S3Endpoint string
	Acl       s3.ACL
	Region string  //RegionName:"us-east-1"
       }
)

func NewS3Client(cnf Config)*S3Client{
	auth := aws.Auth{
		AccessKey: cnf.AccessKey,
		SecretKey: cnf.SecretKey,
	}
	if len(cnf.Region)==0 {cnf.Region= "us-east-1"}

	return &S3Client{
		S3: s3.New(auth, aws.Region{Name: cnf.Region, S3Endpoint: cnf.S3Endpoint }),
		config: cnf,
	}
}


func (s *S3Client) Bucket() *s3.Bucket {
	bucket := s.S3.Bucket(s.config.BucketName)
	//如果bucket不存在,创建bucket
	if err := bucket.PutBucket(s3.PublicReadWrite); err != nil {
		//log.Fatal(err)
		fmt.Println(err)

	}
	return bucket
}

//func (s *S3Client) PutFile(filename string) error{
//	f, er := os.Stat(filename)
//	if er != nil {
//		return   er
//	}
//	if f.Size()>
//	return f.Size(), nil
//
//
//	afile, err := os.Open(filename)
//	if err!=nil{
//	  return err
//	}
//	defer afile.Close()
//	afile.
//}
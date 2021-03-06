package s3

import (
	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-long/goamz/aws"
	"github.com/go-long/goamz/s3"
	"github.com/go-long/longGO/fb/reader"
	"os"
	"time"
	"fmt"
)


type (
	S3Client struct {
		*s3.S3
		config Config
		sliceSize int64
	}

	Config struct {
	  AccessKey string `form_widget:"text"  sql:"type:varchar(50)" valid:"Required"`
	  SecretKey string  `form_widget:"text"   sql:"type:varchar(50)" valid:"Required"`
	  BucketName string `form_widget:"text"  sql:"type:varchar(50)" valid:"Required"`
	  S3Endpoint string `form_widget:"text"  sql:"type:varchar(50)" valid:"Required"`
	  Acl       string ` gorm:"default:'private'" form_widget:"select" form_choices:"|private|私有||public-read|所有人可读||public-read-write|所有人可读写||authenticated-read|authenticated-read||bucket-owner-read|bucket-owner-read||bucket-owner-full-control|bucket-owner-full-control "  sql:"type:varchar(20)" valid:"Required"`
	  Region string  //RegionName:"us-east-1"
       }
)

//uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-west-2")}))
func NewS3Client(cnf Config,sliceSizeMB... int64)*S3Client{
	auth := aws.Auth{
		AccessKey: cnf.AccessKey,
		SecretKey: cnf.SecretKey,
	}
	if len(cnf.Acl)==0 { cnf.Acl=string(s3.PublicReadWrite) }
	if len(cnf.Region)==0 {cnf.Region= "us-east-1"}

	ssize:=int64(15*1024*1024) //默认15M分片大小
	if len(sliceSizeMB)>0{
		ssize=sliceSizeMB[0]*int64(1024*1024)
	}
	return &S3Client{
		S3: s3.New(auth, aws.Region{Name: cnf.Region, S3Endpoint: cnf.S3Endpoint }),
		config: cnf,
		sliceSize:ssize,
	}
}

func (s *S3Client) ListBuckets() ([]s3.Bucket,error) {
	buckets,err:=s.S3.ListBuckets()
	if err!=nil{
		return nil, err
	}

	return  buckets.Buckets,nil
}

func (s *S3Client) Bucket() (*s3.Bucket,error) {
	bucket := s.S3.Bucket(s.config.BucketName)


	//如果bucket不存在,创建bucket
	if err := bucket.PutBucket(s3.ACL(s.config.Acl)); err != nil {
		 return nil,err

	}
	return bucket,nil
}



func (s *S3Client) PutFile(filename string,key string,progressFunc... reader.ProgressReaderCallbackFunc) error{
	f, er := os.Stat(filename)
	if er != nil {
		return  er
	}

	b,err:=s.Bucket()
	if err != nil {
		return err
	}



	afile, err := os.Open(filename)
	if err!=nil{
	  return err
	}
	defer afile.Close()


	//shortName:=filepath.Base(filename)
	//fmt.Println("short:",shortName)
	var multi  *s3.Multi
	progressR := &reader.ReaderSeek{
		Reader: afile,
		Size:   f.Size(),
		DrawFunc: func(m_progress int64, m_size int64,m_speed string) error{
			if len(progressFunc)>0{
				er:=progressFunc[0](m_progress,m_size,m_speed)
				if er!=nil && multi!=nil {
					multi.PutCancel()
					return  er
				}
			}
			return nil
		},
	}






	if f.Size()>64*1024*1024{// *1024{
		fmt.Println("s3 Multi Upload")
		//文件大于64M,必须Multi方式
		multi, err = b.InitMulti(key, "application/octet-stream", s3.ACL(s.config.Acl))
		if err != nil {
			return err
		}
		defer multi.Abort()

		parts, err1 := multi.PutAll(progressR, s.sliceSize)

		if err1 != nil {
			return err1
		}
		fmt.Println("")
		for i, p := range parts {
			fmt.Printf("Processing %d part of %d and uploaded %d bytes. TAG:%s\n ", int(i), int(len(parts)), int(p.Size), p.ETag)
		}
		return   multi.Complete(parts)
	} else {
		err = b.PutReader(key, progressR, f.Size(), "application/octet-stream", s3.ACL(s.config.Acl))
		//err = b.Put("zoujtw2015-12-16.mkv", file, "content-type", s3.PublicReadWrite)
		if err != nil {
			return err
		}
	}
	return nil
}


func (s *S3Client)GetSignedURL(path string,duration time.Duration)string{
	  b,err:=s.Bucket()
	if err!=nil{
		fmt.Println("Bucket err:",err)
	}
	url:=b.SignedURL(path, time.Now().Add(duration))
	return url
}




